package main

import (
	"iter"
	"math/bits"
	"sync"
	"sync/atomic"
)

const (
	OP_APPEND = iota
	OP_DELETE
	OP_DELETE_IDX
)

type setOp[T comparable] struct {
	callback chan bool
	opVal    *T
	opIdx    int
	opType   int
	seqNo    uint64
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
	value  *T
	opIdx  int
}

type OrderedSet[T comparable] struct {
	dualWrite       bool
	bitmap          []uint64
	compacting      bool
	cBitmap         []uint64
	cItems          []uint64
	cDataToSequence map[T]uint64
	dataToSequence  map[T]uint64
	sequenceToData  map[uint64]setMember[T]
	cSequenceToData map[uint64]setMember[T]
	snapshotSeqNo   uint64
	opPool          *sync.Pool
	items           []uint64
	indexMap        map[T]int
	liveCount       atomic.Uint64
	cLiveCount      int
	seqNo           atomic.Uint64 // integer representing the most recent operation number so ops can be tagged with it
	lastAppliedSeq  uint64
	cLastAppliedSeq uint64
	rwLock          *sync.RWMutex
	pending         map[uint64]setOp[T]
	opsCh           chan *setOp[T]
	replayCh        chan replayOp[T]
	tombstones      int
	maxItems        int
	cAppendBitmap   func(int)
	appendBitmap    func(int)
	// cAppend func(
}

func NewOrderedSet[T comparable]() *OrderedSet[T] {
	s := OrderedSet[T]{}
	s.dataToSequence = make(map[T]uint64)
	s.cDataToSequence = make(map[T]uint64)
	s.sequenceToData = make(map[uint64]setMember[T])
	s.cSequenceToData = make(map[uint64]setMember[T])
	s.indexMap = make(map[T]int)
	s.pending = make(map[uint64]setOp[T])
	s.opsCh = make(chan *setOp[T], 10000)
	s.replayCh = make(chan replayOp[T], 10000)
	s.bitmap = make([]uint64, 10000)
	s.cBitmap = make([]uint64, 10000)
	s.items = make([]uint64, 10000)
	s.cItems = make([]uint64, 10000)
	s.opPool = &sync.Pool{}
	s.opPool.New = func() any {
		return &setOp[T]{}
	}
	s.appendBitmap = appendBitmapClosure(&s.bitmap)
	s.cAppendBitmap = appendBitmapClosure(&s.cBitmap)
	s.seqNo = atomic.Uint64{}
	s.seqNo.Store(0)
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
	for op := range s.opsCh {
		applied = false
		s.rwLock.Lock()
		seqNo := s.seqNo.Load() + 1
		s.lastAppliedSeq = seqNo - 1
		op.seqNo = seqNo
		switch op.opType {
		case OP_APPEND:
			applied = s.append(*op.opVal, seqNo)
		case OP_DELETE:
			applied = s.delete(*op.opVal)
		case OP_DELETE_IDX:
			applied = s.deleteLiveIdx(op.opIdx)
		}
		op.seqNo = seqNo
		if s.compacting {
			s.pending[seqNo] = *op
		}
		s.seqNo.Store(seqNo)
		s.rwLock.Unlock()
		op.callback <- applied
		op.opVal, op.callback, op.opIdx = nil, nil, 0
		s.opPool.Put(&op)
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

func (s *OrderedSet[T]) append(elem T, seqNo uint64) bool {
	if _, ok := s.dataToSequence[elem]; !ok {
		bitmapIdx := len(s.items) - 1
		s.dataToSequence[elem] = seqNo
		s.sequenceToData[seqNo] = setMember[T]{bitmapIdx: bitmapIdx, Value: &elem}
		s.items = append(s.items, seqNo)
		s.liveCount.Add(1)
		s.appendBitMap(bitmapIdx)
		return true
	}
	return false
}

func (s *OrderedSet[T]) cAppend(elem T, seqNo uint64) bool {
	if _, ok := s.cDataToSequence[elem]; !ok {
		bitmapIdx := len(s.cItems) - 1
		s.cDataToSequence[elem] = seqNo
		s.cSequenceToData[seqNo] = setMember[T]{bitmapIdx: bitmapIdx, Value: &elem}
		s.cItems = append(s.cItems, seqNo)
		s.cLiveCount++
		s.cAppendBitMap(bitmapIdx)
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

//func (s *OrderedSet[T]) snapshotUnderLock() { // snapshot necessary values for compacting; THIS ASSUMES THE CALLER HAS ACQUIRED THE GLOBAL RWMUTEX, INTEGRITY CANNOT BE VERIFIED OTHERWISE
//	s.cBitmap = make([]uint64, (s.liveCount+63)/64)
//	s.cItems = make([]uint64, s.liveCount)
//}

func (s *OrderedSet[T]) copyItemsUnderLock(newItems *[]uint64, newSequenceToData map[uint64]setMember[T], newDataToSequence map[T]uint64, newBitmap *[]uint64) int {
	liveCount := 0
	appendBitmap := appendBitmapClosure(newBitmap)
	for _, idx := range getAliveIndicesUnderLock(s.bitmap) {
		itemSeqNo := s.items[idx]
		if liveCount > len(*newItems)-1 {
			*newItems = append(*newItems, uint64(0))
		}
		liveCount++
		(*newItems)[liveCount-1] = itemSeqNo
		newSequenceToData[itemSeqNo] = setMember[T]{Value: s.sequenceToData[s.items[idx]].Value, bitmapIdx: liveCount - 1}
		newDataToSequence[*(newSequenceToData[itemSeqNo].Value)] = itemSeqNo
		appendBitmap(liveCount - 1)
	}
	return liveCount
}

func (s *OrderedSet[T]) compact() {
	s.cBitmap = s.cBitmap[:0]
	s.cItems = s.cItems[:0]
	liveCount := int(s.liveCount.Load())
	if s.maxItems > 2*(int(s.liveCount.Load())) {
		s.cSequenceToData = make(map[uint64]setMember[T], liveCount+liveCount/3)
		s.cDataToSequence = make(map[T]uint64, liveCount+liveCount/3)
	} else {
		clear(s.cSequenceToData)
		clear(s.cDataToSequence)
	}

	s.rwLock.Lock()
	s.compacting = true
	if s.pending == nil {
		s.pending = make(map[uint64]setOp[T])
	} else {
		clear(s.pending)
	}
	s.rwLock.Unlock()

	s.rwLock.RLock()
	s.snapshotSeqNo = s.seqNo.Load()
	liveCount = s.copyItemsUnderLock(&s.cItems, s.cSequenceToData, s.cDataToSequence, &s.cBitmap)
	s.cItems = make([]uint64, liveCount, int(float64(liveCount)*1.33))
	s.cLastAppliedSeq = s.snapshotSeqNo
	s.rwLock.RUnlock()

	for currentSeq := s.snapshotSeqNo + 1; (s.seqNo.Load() - currentSeq) > 0; currentSeq++ {
		if op, ok := s.pending[currentSeq]; ok {
			switch op.opType {
			case OP_APPEND:
				s.cAppend(*op.opVal, currentSeq)
			case OP_DELETE:
				s.cDelete(*op.opVal)
			case OP_DELETE_IDX:
				s.cDeleteIdx(op.opIdx)
			}
			s.cLastAppliedSeq = currentSeq
			delete(s.pending, currentSeq)
		}
	}
	s.rwLock.Lock()
	for currentSeq := s.cLastAppliedSeq + 1; (s.seqNo.Load() - currentSeq) > 0; currentSeq++ {
		if op, ok := s.pending[currentSeq]; ok {
			switch op.opType {
			case OP_APPEND:
				s.cAppend(*op.opVal, currentSeq)
			case OP_DELETE:
				s.cDelete(*op.opVal)
			case OP_DELETE_IDX:
				s.cDeleteIdx(op.opIdx)
			}
			s.cLastAppliedSeq = currentSeq
			delete(s.pending, currentSeq)
		}
	}
	s.bitmap, s.dataToSequence, s.sequenceToData, s.items = s.cBitmap, s.cDataToSequence, s.cSequenceToData, s.cItems
	s.cBitmap, s.cDataToSequence, s.cSequenceToData, s.cItems = nil, nil, nil, nil
	s.compacting = false
	s.rwLock.Unlock()
}

func (s *OrderedSet[T]) cAppendBitMap(idx int) {
	uintIndex := idx / 64
	bitOffset := idx % 64
	if uintIndex >= len(s.cBitmap) {
		s.cBitmap = append(s.bitmap, uint64(0))
	}
	bitMask := uint64(1) << bitOffset
	s.cBitmap[uintIndex] |= bitMask
}

//func (s *OrderedSet[T]) appendClosure(items *[]uint64, sequenceToData map[uint64]setMember[T], dataToSequence map[T]uint64, bitmap *[]uint64) func(T) {
//}

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
	op := s.opPool.Get().(*setOp[T])
	op.opVal = nil
	op.opIdx = -1
	op.callback = make(chan bool, 1)
	s.opsCh <- op
	return <-op.callback
}

func (s *OrderedSet[T]) delete(elem T) bool {
	if seqNo, ok := s.dataToSequence[elem]; ok {
		if member, ok2 := s.sequenceToData[seqNo]; ok2 {
			s.deleteBitMap(member.bitmapIdx)
			//			delete(s.dataToSequence, elem)
			//			delete(s.sequenceToData, seqNo)
			s.liveCount.Store(s.liveCount.Load() - 1)
			return true
		}
	}
	return false
}

func (s *OrderedSet[T]) cDelete(elem T) bool {
	if seqNo, ok := s.cDataToSequence[elem]; ok {
		if member, ok2 := s.cSequenceToData[seqNo]; ok2 {
			s.deleteBitMap(member.bitmapIdx)
			//			delete(s.cDataToSequence, elem)
			//			delete(s.cSequenceToData, seqNo)
			s.liveCount.Store(s.liveCount.Load() - 1)
			return true
		}
	}
	return false
}

func (s *OrderedSet[T]) DeleteIdx(idx int) bool {
	op := s.opPool.Get().(*setOp[T])
	op.opIdx = idx
	op.opVal = nil
	op.callback = make(chan bool, 1)
	s.opsCh <- op
	return <-op.callback
}

func (s *OrderedSet[T]) deleteLiveIdx(idx int) bool {
	itemIdx := getNthAliveIndexUnderLock(s.bitmap, idx)
	if itemIdx >= 0 {
		seqNo := s.items[itemIdx]
		member, ok := s.sequenceToData[seqNo]
		if ok {
			s.deleteBitMap(member.bitmapIdx)
			//			delete(s.cDataToSequence, *member.Value)
			//			delete(s.cSequenceToData, seqNo)
			s.liveCount.Store(s.liveCount.Load() - 1)
			return true
		}
	}
	return false
}

func (s *OrderedSet[T]) cDeleteIdx(idx int) bool {
	itemIdx := getNthAliveIndexUnderLock(s.bitmap, idx)
	if itemIdx >= 0 {
		seqNo := s.items[itemIdx]
		member, ok := s.sequenceToData[seqNo]
		if ok {
			s.deleteBitMap(member.bitmapIdx)
			delete(s.cDataToSequence, *member.Value)
			delete(s.cSequenceToData, seqNo)
			s.liveCount.Store(s.liveCount.Load() - 1)
			return true
		}
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
		for _, seqNo := range s.items {
			val := *s.sequenceToData[seqNo].Value
			if s.isValAlive(val) {
				if !(yield(val)) {
					return
				}
			}
		}
	}
}

func (s *OrderedSet[T]) IterIndex() iter.Seq2[int, T] {
	return func(yield func(int, T) bool) {
		liveIdx := 0
		for _, item := range getAliveIndicesUnderLock(s.bitmap) {
			val := *s.sequenceToData[s.items[item]].Value
			if s.isValAlive(val) {
				if !(yield(item, val)) {
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
