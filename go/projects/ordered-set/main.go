package main

import (
	"iter"
	"sync"
)

const (
	OP_APPEND = iota
	OP_DELETE
)

type setOp[T comparable] struct {
	replayOp[T]
	callback chan bool
}

type replayOp[T comparable] struct {
	opType int
	opVal  T
	seqNo  uint64
}

type OrderedSet[T comparable] struct {
	dualWrite           bool
	bitmap              []uint64
	compacting          bool
	compactingBitMap    []uint64
	compactingItems     []T
	compactingIndexMap  map[T]int
	compactingLiveCount int
	// callbackPool        *sync.Pool
	opPool         *sync.Pool
	items          []T
	indexMap       map[T]int
	liveCount      int
	seqNo          uint64 // integer representing the most recent operation number so ops can be tagged with it
	lastAppliedSeq uint64
	rwLock         *sync.RWMutex
	opsCh          chan *setOp[T]
	replayCh       chan replayOp[T]
	tombstones     int
}

func (s *OrderedSet[T]) sequencer() {
	var applied bool
	var replay replayOp[T]
	var needReplay bool
	for op := range s.opsCh {
		needReplay = false
		applied = false
		s.rwLock.Lock()
		s.lastAppliedSeq = s.seqNo
		s.seqNo++
		switch op.opType {
		case OP_APPEND:
			applied = s.append(op.opVal)
		case OP_DELETE:
			applied = s.delete(op.opVal)
		}
		op.seqNo = s.seqNo
		if s.compacting {
			replay = op.replayOp
			needReplay = true
		}
		s.rwLock.Unlock()
		if needReplay {
			s.replayCh <- replay
		}
		op.callback <- applied
		s.opPool.Put(op)
	}
}

func (s *OrderedSet[T]) Append(elem T) bool {
	op := s.opPool.Get().(*setOp[T])
	op.opType = OP_APPEND
	op.opVal = elem
	op.callback = make(chan bool, 1)
	s.opsCh <- op
	return <-op.callback
}

func (s *OrderedSet[T]) append(elem T) bool {
	if _, ok := s.indexMap[elem]; !ok {
		s.items = append(s.items, elem)
		s.indexMap[elem] = len(s.items) - 1
		s.liveCount++
		s.appendBitMap(s.indexMap[elem])
		return true
	}
	return false
}

func (s *OrderedSet[T]) appendBitMap(idx int) {
	uintIndex := idx / 64
	bitOffset := idx % 64
	if uintIndex >= len(s.bitmap) {
		s.bitmap = append(s.bitmap, uint64(0))
	}
	bitMask := uint64(1) << bitOffset
	s.bitmap[uintIndex] |= bitMask
}

func (s *OrderedSet[T]) compact() {
	newIndexMap := make(map[T]int)
	s.rwLock.RLock()
	s.compacting = true
	s.compactingBitMap = make([]uint64, (s.liveCount+63)/64)
	s.compactingItems = make([]T, s.liveCount)
	copy(s.compactingBitMap, s.bitmap[:len(s.compactingBitMap)-1])
	liveCount := 0
	for idx := range s.items {
		if isAlive(s.compactingBitMap, idx) {
			s.compactingItems[liveCount] = s.items[idx]
			newIndexMap[s.items[idx]] = liveCount
			liveCount++
		}
	}
	s.replayCh = make(chan replayOp[T], 10000)
	s.rwLock.RUnlock()
	bitmapRemainder := liveCount % 64
	for idx := range s.compactingBitMap {
		s.compactingBitMap[idx] = ^uint64(0)
	}
	if bitmapRemainder != 0 {
		s.compactingBitMap[len(s.compactingBitMap)-1] = (uint64(1) << bitmapRemainder) - 1
	}

	s.rwLock.Lock()
	close(s.replayCh)
	s.replayCh = nil
	s.dualWrite = true
	s.rwLock.Unlock()
	wg := &sync.WaitGroup{}
	go s.replay(wg)
	wg.Wait()
	s.rwLock.Lock()
	s.bitmap, s.indexMap, s.items = s.compactingBitMap, s.compactingIndexMap, s.compactingItems
	s.compactingBitMap, s.compactingIndexMap, s.compactingItems = nil, nil, nil
	s.compacting = false
	s.rwLock.Unlock()
}

func (s *OrderedSet[T]) compactingAppendBitMap(idx int) {
	uintIndex := idx / 64
	bitOffset := idx % 64
	if uintIndex >= len(s.compactingBitMap) {
		s.compactingBitMap = append(s.bitmap, uint64(0))
	}
	bitMask := uint64(1) << bitOffset
	s.compactingBitMap[uintIndex] |= bitMask
}

func (s *OrderedSet[T]) compactingAppend(elem T) bool {
	s.rwLock.Lock()
	if _, ok := s.compactingIndexMap[elem]; !ok {
		s.items = append(s.items, elem)
		s.compactingIndexMap[elem] = len(s.items) - 1
		s.compactingLiveCount += 1
		s.compactingAppendBitMap(s.compactingIndexMap[elem])
		s.rwLock.Unlock()
		return true
	}
	s.rwLock.Unlock()
	return false
}

func (s *OrderedSet[T]) replay(wg *sync.WaitGroup) {
	wg.Add(1)
	for op := range s.replayCh {
		if op.opType == OP_APPEND {
			s.compactingAppend(op.opVal)
		}
		if op.opType == OP_DELETE {
			s.compactingDelete(op.opVal)
		}
	}
	wg.Done()
}

func (s *OrderedSet[T]) deleteBitMap(idx int) {
	uintIndex := idx / 64
	bitOffset := idx % 64
	if uintIndex < len(s.bitmap) {
		bitMask := uint64(1) << bitOffset
		s.bitmap[uintIndex] &^= bitMask
	}
}

func (s *OrderedSet[T]) compactingDeleteBitMap(idx int) {
	uintIndex := idx / 64
	bitOffset := idx % 64
	if uintIndex < len(s.bitmap) {
		bitMask := uint64(1) << bitOffset
		s.compactingBitMap[uintIndex] &^= bitMask
	}
}

func (s *OrderedSet[T]) Delete(elem T) bool {
}

func (s *OrderedSet[T]) delete(elem T) bool {
	if idx, ok := s.indexMap[elem]; ok {
		s.deleteBitMap(idx)
		delete(s.indexMap, elem)
		s.liveCount--
		return true
	}
	return false
}

func (s *OrderedSet[T]) compactingDelete(elem T) bool {
	s.rwLock.Lock()
	if idx, ok := s.compactingIndexMap[elem]; ok {
		s.compactingDeleteBitMap(idx)
		delete(s.compactingIndexMap, elem)
		s.compactingLiveCount--
		s.rwLock.Unlock()
		return true
	}
	s.rwLock.Unlock()
	return false
}

func (s *OrderedSet[T]) DeleteIdx(idx int) bool {
	if s.isIdxAlive(idx) {
		s.deleteBitMap(idx)
		return true
	}
	return false
}

func (s *OrderedSet[T]) isValAlive(elem T) bool {
	if _, ok := s.indexMap[elem]; ok {
		return true
	}
	return false
}

// func getBitMapIndex(idx int) uint64
func (s *OrderedSet[T]) isIdxAlive(idx int) bool {
	if idx < 0 || idx > len(s.bitmap) {
		return false
	}
	uintIndex := idx / 64
	bitOffset := idx % 64
	bitMask := uint64(1) << bitOffset
	return (s.bitmap[uintIndex] & bitMask) != 0
}

func (s *OrderedSet[T]) In(elem T) bool {
	if _, ok := s.indexMap[elem]; ok {
		return true
	} else {
		return false
	}
}

func (s *OrderedSet[T]) Iter() iter.Seq[T] {
	return func(yield func(T) bool) {
		for _, item := range s.items {
			if s.isValAlive(item) {
				if !(yield(item)) {
					return
				}
			}
		}
	}
}

func (s *OrderedSet[T]) IterIndex() iter.Seq2[int, T] {
	return func(yield func(int, T) bool) {
		liveIdx := 0
		for _, item := range s.items {
			if s.isValAlive(item) {
				if !(yield(liveIdx, item)) {
					return
				}
				liveIdx++
			}
		}
	}
}

func isAlive(bitmap []uint64, idx int) bool {
	if idx < 0 || idx > len(bitmap) {
		return false
	}
	uintIndex := idx / 64
	bitOffset := idx % 64
	bitMask := uint64(1) << bitOffset
	return (bitmap[uintIndex] & bitMask) != 0
}
