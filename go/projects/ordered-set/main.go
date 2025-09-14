package main

import (
	"iter"
	"math/bits"
	"sync"
)

const (
	OP_APPEND = iota
	OP_DELETE
	OP_DELETE_IDX
)

type setOp[T comparable] struct {
	replayOp[T]
	callback chan bool
	opVal    *T
}

type setMember[T comparable] struct {
	bitmapIdx int
	Value     *T
}

type setEntity[T comparable] struct {
	Value T
	ID    uint64
}

type replayOp[T comparable] struct {
	opType int
	seqNo  uint64
}

type OrderedSet[T comparable] struct {
	dualWrite       bool
	bitmap          []uint64
	compacting      bool
	cBitMap         []uint64
	cItems          []uint64
	cDataToSequence map[T]uint64
	cSequenceToData map[uint64]setMember[T]
	cLiveCount      int
	sequenceToData  map[uint64]setMember[T]
	dataToSequence  map[T]uint64
	snapshotSeqNo   uint64
	opPool          *sync.Pool
	items           []uint64
	indexMap        map[T]int
	liveCount       int
	seqNo           uint64 // integer representing the most recent operation number so ops can be tagged with it
	lastAppliedSeq  uint64
	rwLock          *sync.RWMutex
	pending         map[uint64]*T
	opsCh           chan *setOp[T]
	replayCh        chan replayOp[T]
	tombstones      int
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
			applied = s.append(*op.opVal)
		case OP_DELETE:
			applied = s.delete(*op.opVal)
		}
		op.seqNo = s.seqNo
		if s.compacting {
			replay = op.replayOp
			needReplay = true
			s.pending[s.seqNo] = op.opVal
		}
		s.rwLock.Unlock()
		if needReplay {
			s.replayCh <- replay
		}
		op.callback <- applied
		op.opVal, op.callback = nil, nil
		s.opPool.Put(op)
	}
}

func (s *OrderedSet[T]) Append(elem T) bool {
	op := s.opPool.Get().(*setOp[T])
	op.opType = OP_APPEND
	op.opVal = &elem
	op.callback = make(chan bool, 1)
	s.opsCh <- op
	return <-op.callback
}

func (s *OrderedSet[T]) append(elem T) bool {
	if _, ok := s.dataToSequence[elem]; !ok {
		bitmapIdx := len(s.items) - 1
		s.dataToSequence[elem] = s.seqNo
		s.sequenceToData[s.seqNo] = setMember[T]{bitmapIdx: bitmapIdx - 1, Value: &elem}
		s.items = append(s.items, s.seqNo)
		s.liveCount++
		s.appendBitMap(s.indexMap[elem])
		return true
	}
	return false
}

func (s *OrderedSet[T]) Get(value T) (setEntity[T], bool) {
	if seqNo, ok := s.dataToSequence[value]; ok {
		if s.isSeqAlive(seqNo) {
			return setEntity[T]{
				Value: *s.sequenceToData[seqNo].Value,
				ID:    seqNo,
			}, true
		}
	}
	return setEntity[T]{}, false
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

func (s *OrderedSet[T]) snapshotUnderLock() { // snapshot necessary values for compacting; THIS ASSUMES THE CALLER HAS ACQUIRED THE GLOBAL RWMUTEX, INTEGRITY CANNOT BE VERIFIED OTHERWISE
	s.cBitMap = make([]uint64, (s.liveCount+63)/64)
	s.cItems = make([]uint64, s.liveCount)
}

func (s *OrderedSet[T]) copyItemsUnderLock(newItems *[]uint64, sequenceToData map[uint64]setMember[T], dataToSequence map[T]uint64) int {
	liveCount := 0
	for idx := range s.items {
		if isAlive(s.cBitMap, idx) {
			itemSeqNo := s.items[idx]
			if liveCount > len(*newItems)-1 {
				*newItems = append(*newItems, uint64(0))
			}
			(*newItems)[liveCount] = itemSeqNo
			sequenceToData[itemSeqNo] = setMember[T]{Value: s.sequenceToData[s.items[idx]].Value, bitmapIdx: liveCount}
			dataToSequence[*(sequenceToData[itemSeqNo].Value)] = itemSeqNo
			liveCount++
		}
	}
	return liveCount
}

func (s *OrderedSet[T]) compact() {
	s.rwLock.RLock()
	s.compacting = true
	s.cBitMap = make([]uint64, (s.liveCount+63)/64)
	s.cItems = make([]uint64, s.liveCount)
	s.snapshotSeqNo = s.seqNo
	copy(s.cBitMap, s.bitmap[:len(s.cBitMap)-1])
	liveCount := s.copyItemsUnderLock(&s.cItems, s.cSequenceToData)

	//	for idx := range s.items {
	//		if isAlive(s.cBitMap, idx) {
	//			itemSeqNo = s.items[idx]
	//			s.cItems[liveCount] = itemSeqNo
	//			s.cSequenceToData[itemSeqNo] = setMember{Value: s.sequenceToData[s.items[idx]], bitmapIdx: liveCount}
	//			liveCount++
	//		}
	//	}
	s.replayCh = make(chan replayOp[T], 10000)
	s.rwLock.RUnlock()
	bitmapRemainder := liveCount % 64
	for idx := range s.cBitMap {
		s.cBitMap[idx] = ^uint64(0)
	}
	if bitmapRemainder != 0 {
		s.cBitMap[len(s.cBitMap)-1] = (uint64(1) << bitmapRemainder) - 1
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
	s.bitmap, s.indexMap, s.items = s.cBitMap, s.cIndexMap, s.cItems
	s.cBitMap, s.cIndexMap, s.cItems = nil, nil, nil
	s.compacting = false
	s.rwLock.Unlock()
}

func (s *OrderedSet[T]) compactingAppendBitMap(idx int) {
	uintIndex := idx / 64
	bitOffset := idx % 64
	if uintIndex >= len(s.cBitMap) {
		s.cBitMap = append(s.bitmap, uint64(0))
	}
	bitMask := uint64(1) << bitOffset
	s.cBitMap[uintIndex] |= bitMask
}

// func appendClosure(s *OrderedSet[T], setItems []T,
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
		s.cBitMap[uintIndex] &^= bitMask
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
	if idx, ok := s.cIndexMap[elem]; ok {
		s.compactingDeleteBitMap(idx)
		delete(s.cIndexMap, elem)
		s.cLiveCount--
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

func (s *OrderedSet[T]) deleteLiveIdx(idx int) bool {
	liveIdx := -1
	for i := range s.items {
		if s.isIdxAlive(i) {
			liveIdx++
		}
		if liveIdx == idx {
			s.deleteBitMap(i)
		}
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

func (s *OrderedSet[T]) isSeqAlive(seqNo uint64) bool {
	if member, ok := s.sequenceToData[seqNo]; ok {
		if s.isIdxAlive(member.bitmapIdx) {
			return true
		}
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
