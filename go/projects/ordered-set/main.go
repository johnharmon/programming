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
	opIdx    int
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
	cBitmap         []uint64
	cItems          []uint64
	cDataToSequence map[T]uint64
	dataToSequence  map[T]uint64
	cLiveCount      int
	sequenceToData  map[uint64]setMember[T]
	cSequenceToData map[uint64]setMember[T]
	snapshotSeqNo   uint64
	opPool          *sync.Pool
	items           []uint64
	indexMap        map[T]int
	liveCount       int
	seqNo           uint64 // integer representing the most recent operation number so ops can be tagged with it
	lastAppliedSeq  uint64
	rwLock          *sync.RWMutex
	pending         map[uint64]*T
	opsCh           chan setOp[T]
	replayCh        chan replayOp[T]
	tombstones      int
	cAppendBitmap   func(int)
	maxItems        int
	appendBitmap    func(int)
}

func NewOrderedSet[T comparable]() *OrderedSet[T] {
	s := OrderedSet[T]{}
	s.dataToSequence = make(map[T]uint64)
	s.cDataToSequence = make(map[T]uint64)
	s.sequenceToData = make(map[uint64]setMember[T])
	s.cSequenceToData = make(map[uint64]setMember[T])
	s.indexMap = make(map[T]int)
	s.pending = make(map[uint64]*T)
	s.opsCh = make(chan setOp[T], 10000)
	s.replayCh = make(chan replayOp[T], 10000)
	s.bitmap = make([]uint64, 10000)
	s.cBitmap = make([]uint64, 10000)
	s.items = make([]uint64, 10000)
	s.cItems = make([]uint64, 10000)
	s.opPool = &sync.Pool{}
	s.opPool.New = func() any {
		return setOp[T]{}
	}
	s.appendBitmap = appendBitmapClosure(&s.bitmap)
	s.cAppendBitmap = appendBitmapClosure(&s.cBitmap)
	return &s
}

func getAliveIndicesUnderLock(bitmap []uint64) (liveIndices []int) {
	bitmapLength := len(bitmap)
	if bitmapLength < 1 {
		return liveIndices
	}
	for idx, word := range bitmap {
		if word != 0 {
			for word != 0 {
				liveIndices = append(liveIndices, (64*idx)+bits.TrailingZeros64(word))
				word &= word - 1
			}
		}
	}
	return liveIndices
}

func getNthAliveIndexUnderLock(bitmap []uint64, target int) int { // only do this under lock, else it cannot guarantee index accuracy or item liveleness
	if len(bitmap) < 1 || target < 0 {
		return -1
	}
	actualIndex := 0 // 'index' value for the bit that represents 'Nth' 1 found
	aliveIndex := 0  // variable to keep track of how many ones we have found
	targetIdx := 0   // index of the word in the bitmap that holds the value we care about
	found := false
	for idx := range bitmap {
		onesCount := bits.OnesCount64(bitmap[idx])
		tmp := aliveIndex + onesCount
		if tmp <= target {
			aliveIndex += onesCount
			actualIndex += 64
		} else {
			found = true
			targetIdx = idx
			break
		}
	}
	if !found {
		return -1
	}
	word := bitmap[targetIdx]
	for i := 0; i < (target - aliveIndex); i++ {
		word &= word - 1
	}
	return actualIndex + bits.TrailingZeros64(word)
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
		case OP_DELETE_IDX:
			applied = s.deleteLiveIdx(op.opIdx)
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
		s.sequenceToData[s.seqNo] = setMember[T]{bitmapIdx: bitmapIdx, Value: &elem}
		s.items = append(s.items, s.seqNo)
		s.liveCount++
		s.appendBitMap(bitmapIdx)
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

func appendBitmapClosure(bitmap *[]uint64) func(int) {
	return func(idx int) {
		uintIndex := idx / 64
		bitOffset := idx % 64
		if uintIndex >= len(*bitmap) {
			*bitmap = append(*bitmap, uint64(0))
		}
		bitMask := uint64(1) << bitOffset
		(*bitmap)[uintIndex] |= bitMask
	}
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
	s.cBitmap = make([]uint64, (s.liveCount+63)/64)
	s.cItems = make([]uint64, s.liveCount)
}

func (s *OrderedSet[T]) copyItemsUnderLock(newItems *[]uint64, newSequenceToData map[uint64]setMember[T], newDataToSequence map[T]uint64, newBitmap *[]uint64) int {
	liveCount := 0
	for _, idx := range getAliveIndicesUnderLock(s.bitmap) {
		itemSeqNo := s.items[idx]
		if liveCount > len(*newItems)-1 {
			*newItems = append(*newItems, uint64(0))
		}
		(*newItems)[liveCount] = itemSeqNo
		newSequenceToData[itemSeqNo] = setMember[T]{Value: s.sequenceToData[s.items[idx]].Value, bitmapIdx: liveCount}
		newDataToSequence[*(newSequenceToData[itemSeqNo].Value)] = itemSeqNo
		s.cAppendBitmap(liveCount)
		liveCount++
	}
	return liveCount
}

func (s *OrderedSet[T]) compact() {
	s.cBitmap = s.cBitmap[:0]
	s.cItems = s.cItems[:0]
	if s.maxItems > 2*s.liveCount {
		s.cSequenceToData = make(map[uint64]setMember[T], s.liveCount+s.liveCount/3)
		s.cDataToSequence = make(map[T]uint64, s.liveCount+s.liveCount/3)
	} else {
		clear(s.cSequenceToData)
		clear(s.cDataToSequence)
	}

	s.rwLock.Lock()
	s.compacting = true
	if s.pending == nil {
		s.pending = make(map[uint64]*T)
	} else {
		clear(s.pending)
	}
	s.rwLock.Unlock()

	s.rwLock.RLock()
	s.snapshotSeqNo = s.seqNo
	liveCount := s.copyItemsUnderLock(&s.cItems, s.cSequenceToData, s.cDataToSequence, &s.cBitmap)
	s.cItems = make([]uint64, liveCount, int(float64(liveCount)*1.33))
	s.rwLock.RUnlock()

	s.replayCh = make(chan replayOp[T], 10000)

	s.rwLock.Lock()
	close(s.replayCh)
	s.replayCh = nil
	s.dualWrite = true
	s.rwLock.Unlock()

	wg := &sync.WaitGroup{}
	go s.replay(wg)
	wg.Wait()

	s.rwLock.Lock()
	s.bitmap, s.dataToSequence, s.sequenceToData, s.items = s.cBitmap, s.cDataToSequence, s.cSequenceToData, s.cItems
	s.cBitmap, s.cDataToSequence, s.cSequenceToData, s.cItems = nil, nil, nil, nil
	s.compacting = false
	s.rwLock.Unlock()
}

func (s *OrderedSet[T]) compactingAppendBitMap(idx int) {
	uintIndex := idx / 64
	bitOffset := idx % 64
	if uintIndex >= len(s.cBitmap) {
		s.cBitmap = append(s.bitmap, uint64(0))
	}
	bitMask := uint64(1) << bitOffset
	s.cBitmap[uintIndex] |= bitMask
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
		s.cBitmap[uintIndex] &^= bitMask
	}
}

func (s *OrderedSet[T]) Delete(elem T) bool {
	op := s.opPool.Get().(setOp[T])
	op.opVal = nil
	op.opIdx = -1
	op.callback = make(chan bool, 1)
	s.opsCh <- op
	return <-op.callback
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
	if seqNo, ok := s.cDataToSequence[elem]; ok {
		s.compactingDeleteBitMap(int(seqNo))
		delete(s.cDataToSequence, elem)
		delete(s.cSequenceToData(seqNo)
		s.cLiveCount--
		s.rwLock.Unlock()
		return true
	}
	s.rwLock.Unlock()
	return false
}

func (s *OrderedSet[T]) DeleteIdx(idx int) bool {
	op := s.opPool.Get().(setOp[T])
	op.opIdx = idx
	op.opVal = nil
	op.callback = make(chan bool, 1)
	s.opsCh <- op
	return <-op.callback
}

func (s *OrderedSet[T]) deleteLiveIdx(idx int) bool {
	liveIdx := -1
	for i := range s.items {
		if s.isIdxAlive(i) {
			liveIdx++
		}
		if liveIdx == seqNo {
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
	if seqNo < 0 || seqNo > len(s.bitmap) {
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
	if seqNo < 0 || seqNo > len(bitmap) {
		return false
	}
	uintIndex := idx / 64
	bitOffset := idx % 64
	bitMask := uint64(1) << bitOffset
	return (bitmap[uintIndex] & bitMask) != 0
}
