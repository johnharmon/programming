package set

//
//import (
//	"iter"
//	"sync"
//	"sync/atomic"
//)
//
//func NewAnyAnyOrderedSet() *AnyOrderedSet {
//	s := AnyOrderedSet{}
//	s.dataToSequence = make(map[any]uint64)
//	s.cDataToSequence = make(map[any]uint64)
//	s.sequenceToData = make(map[uint64]anySetMember)
//	s.cSequenceToData = make(map[uint64]anySetMember)
//	s.indexMap = make(map[any]int)
//	s.pending = make(map[uint64]anySetOp)
//	s.opsCh = make(chan *anySetOp, 10000)
//	s.bitmap = make([]uint64, 10000)
//	s.cBitmap = make([]uint64, 10000)
//	s.items = make([]uint64, 10000)
//	s.cItems = make([]uint64, 10000)
//	s.opPool = &sync.Pool{}
//	s.opPool.New = func() any {
//		return &anySetOp{}
//	}
//	s.appendBitmap = appendBitmapClosure(&s.bitmap)
//	s.cAppendBitmap = appendBitmapClosure(&s.cBitmap)
//	s.seqNo = atomic.Uint64{}
//	s.seqNo.Store(0)
//	go s.sequencer()
//	return &s
//}
//
//func (s *AnyOrderedSet) sequencer() {
//	var applied bool
//	for op := range s.opsCh {
//		applied = false
//		s.rwLock.Lock()
//		seqNo := s.seqNo.Load() + 1
//		s.lastAppliedSeq = seqNo - 1
//		op.seqNo = seqNo
//		switch op.opType {
//		case OP_APPEND:
//			applied = s.append(op.opVal, seqNo)
//		case OP_DELETE:
//			applied = s.delete(op.opVal)
//		case OP_DELETE_IDX:
//			applied = s.deleteLiveIdx(op.opIdx)
//		}
//		op.seqNo = seqNo
//		if s.compacting {
//			s.pending[seqNo] = *op
//		}
//		s.seqNo.Store(seqNo)
//		s.rwLock.Unlock()
//		op.callback <- applied
//		op.opVal, op.callback, op.opIdx = nil, nil, 0
//		s.opPool.Put(&op)
//	}
//}
//
//func (s *AnyOrderedSet) Append(elem any) bool {
//	op := s.opPool.Get().(*anySetOp)
//	op.opType = OP_APPEND
//	op.opVal = &elem
//	op.callback = make(chan bool, 1)
//	s.opsCh <- op
//	return <-op.callback
//}
//
//func (s *AnyOrderedSet) append(elem any, seqNo uint64) bool {
//	if _, ok := s.dataToSequence[elem]; !ok {
//		bitmapIdx := len(s.items) - 1
//		s.dataToSequence[elem] = seqNo
//		s.sequenceToData[seqNo] = anySetMember{bitmapIdx: bitmapIdx, Value: &elem}
//		s.items = append(s.items, seqNo)
//		s.liveCount.Add(1)
//		s.appendBitMap(bitmapIdx)
//		return true
//	}
//	return false
//}
//
//func (s *AnyOrderedSet) cAppend(elem any, seqNo uint64) bool {
//	if _, ok := s.cDataToSequence[elem]; !ok {
//		bitmapIdx := len(s.cItems) - 1
//		s.cDataToSequence[elem] = seqNo
//		s.cSequenceToData[seqNo] = anySetMember{bitmapIdx: bitmapIdx, Value: &elem}
//		s.cItems = append(s.cItems, seqNo)
//		s.cLiveCount++
//		s.cAppendBitMap(bitmapIdx)
//		return true
//	}
//	return false
//}
//
//func (s *AnyOrderedSet) Get(value any) (anySetEntity, bool) {
//	if seqNo, ok := s.dataToSequence[value]; ok {
//		if s.isSeqAlive(seqNo) {
//			return anySetEntity{
//				Value: s.sequenceToData[seqNo].Value,
//				ID:    seqNo,
//			}, true
//		}
//	}
//	return anySetEntity{}, false
//}
//
//func (s *AnyOrderedSet) appendBitMap(idx int) {
//	uintIndex := idx / 64
//	bitOffset := idx % 64
//	if uintIndex >= len(s.bitmap) {
//		s.bitmap = append(s.bitmap, uint64(0))
//	}
//	bitMask := uint64(1) << bitOffset
//	s.bitmap[uintIndex] |= bitMask
//}
//
//func (s *AnyOrderedSet) copyItemsUnderLock(newItems *[]uint64, newSequenceToData map[uint64]anySetMember, newDataToSequence map[any]uint64, newBitmap *[]uint64) int {
//	liveCount := 0
//	appendBitmap := appendBitmapClosure(newBitmap)
//	for _, idx := range getAliveIndicesUnderLock(s.bitmap) {
//		itemSeqNo := s.items[idx]
//		if liveCount > len(*newItems)-1 {
//			*newItems = append(*newItems, uint64(0))
//		}
//		liveCount++
//		(*newItems)[liveCount-1] = itemSeqNo
//		newSequenceToData[itemSeqNo] = setMember{Value: s.sequenceToData[s.items[idx]].Value, bitmapIdx: liveCount - 1}
//		newDataToSequence[*(newSequenceToData[itemSeqNo].Value)] = itemSeqNo
//		appendBitmap(liveCount - 1)
//	}
//	return liveCount
//}
//
//func (s *AnyOrderedSet) compact() {
//	s.cBitmap = s.cBitmap[:0]
//	s.cItems = s.cItems[:0]
//	liveCount := int(s.liveCount.Load())
//	if s.maxItems > 2*(int(s.liveCount.Load())) {
//		s.cSequenceToData = make(map[uint64]setMember, liveCount+liveCount/3)
//		s.cDataToSequence = make(mapuint64, liveCount+liveCount/3)
//	} else {
//		clear(s.cSequenceToData)
//		clear(s.cDataToSequence)
//	}
//
//	s.rwLock.Lock()
//	s.compacting = true
//	if s.pending == nil {
//		s.pending = make(map[uint64]setOp)
//	} else {
//		clear(s.pending)
//	}
//	s.rwLock.Unlock()
//
//	s.rwLock.RLock()
//	s.snapshotSeqNo = s.seqNo.Load()
//	liveCount = s.copyItemsUnderLock(&s.cItems, s.cSequenceToData, s.cDataToSequence, &s.cBitmap)
//	s.cItems = make([]uint64, liveCount, int(float64(liveCount)*1.33))
//	s.cLastAppliedSeq = s.snapshotSeqNo
//	s.rwLock.RUnlock()
//
//	for currentSeq := s.snapshotSeqNo + 1; (s.seqNo.Load() - currentSeq) > 0; currentSeq++ {
//		if op, ok := s.pending[currentSeq]; ok {
//			switch op.opType {
//			case OP_APPEND:
//				s.cAppend(*op.opVal, currentSeq)
//			case OP_DELETE:
//				s.cDelete(*op.opVal)
//			case OP_DELETE_IDX:
//				s.cDeleteIdx(op.opIdx)
//			}
//			s.cLastAppliedSeq = currentSeq
//			delete(s.pending, currentSeq)
//		}
//	}
//	s.rwLock.Lock()
//	for currentSeq := s.cLastAppliedSeq + 1; (s.seqNo.Load() - currentSeq) > 0; currentSeq++ {
//		if op, ok := s.pending[currentSeq]; ok {
//			switch op.opType {
//			case OP_APPEND:
//				s.cAppend(*op.opVal, currentSeq)
//			case OP_DELETE:
//				s.cDelete(*op.opVal)
//			case OP_DELETE_IDX:
//				s.cDeleteIdx(op.opIdx)
//			}
//			s.cLastAppliedSeq = currentSeq
//			delete(s.pending, currentSeq)
//		}
//	}
//	s.bitmap, s.dataToSequence, s.sequenceToData, s.items = s.cBitmap, s.cDataToSequence, s.cSequenceToData, s.cItems
//	s.cBitmap, s.cDataToSequence, s.cSequenceToData, s.cItems = nil, nil, nil, nil
//	s.compacting = false
//	s.rwLock.Unlock()
//}
//
//func (s *AnyOrderedSet) cAppendBitMap(idx int) {
//	uintIndex := idx / 64
//	bitOffset := idx % 64
//	if uintIndex >= len(s.cBitmap) {
//		s.cBitmap = append(s.bitmap, uint64(0))
//	}
//	bitMask := uint64(1) << bitOffset
//	s.cBitmap[uintIndex] |= bitMask
//}
//
////func (s *AnyOrderedSet) appendClosure(items *[]uint64, sequenceToData map[uint64]setMember[T], dataToSequence map[T]uint64, bitmap *[]uint64) func(T) {
////}
//
//func (s *AnyOrderedSet) deleteBitMap(idx int) {
//	uintIndex := idx / 64
//	bitOffset := idx % 64
//	if uintIndex < len(s.bitmap) {
//		bitMask := uint64(1) << bitOffset
//		s.bitmap[uintIndex] &^= bitMask
//	}
//}
//
//func (s *AnyOrderedSet) compactingDeleteBitMap(idx int) {
//	uintIndex := idx / 64
//	bitOffset := idx % 64
//	if uintIndex < len(s.bitmap) {
//		bitMask := uint64(1) << bitOffset
//		s.cBitmap[uintIndex] &^= bitMask
//	}
//}
//
//func (s *AnyOrderedSet) Delete(elem T) bool {
//	op := s.opPool.Get().(*setOp)
//	op.opVal = nil
//	op.opIdx = -1
//	op.callback = make(chan bool, 1)
//	s.opsCh <- op
//	return <-op.callback
//}
//
//func (s *AnyOrderedSet) delete(elem T) bool {
//	if seqNo, ok := s.dataToSequence[elem]; ok {
//		if member, ok2 := s.sequenceToData[seqNo]; ok2 {
//			s.deleteBitMap(member.bitmapIdx)
//			//			delete(s.dataToSequence, elem)
//			//			delete(s.sequenceToData, seqNo)
//			s.liveCount.Store(s.liveCount.Load() - 1)
//			return true
//		}
//	}
//	return false
//}
//
//func (s *AnyOrderedSet) cDelete(elem T) bool {
//	if seqNo, ok := s.cDataToSequence[elem]; ok {
//		if member, ok2 := s.cSequenceToData[seqNo]; ok2 {
//			s.deleteBitMap(member.bitmapIdx)
//			//			delete(s.cDataToSequence, elem)
//			//			delete(s.cSequenceToData, seqNo)
//			s.liveCount.Store(s.liveCount.Load() - 1)
//			return true
//		}
//	}
//	return false
//}
//
//func (s *AnyOrderedSet) DeleteIdx(idx int) bool {
//	op := s.opPool.Get().(*setOp)
//	op.opIdx = idx
//	op.opVal = nil
//	op.callback = make(chan bool, 1)
//	s.opsCh <- op
//	return <-op.callback
//}
//
//func (s *AnyOrderedSet) deleteLiveIdx(idx int) bool {
//	itemIdx := getNthAliveIndexUnderLock(s.bitmap, idx)
//	if itemIdx >= 0 {
//		seqNo := s.items[itemIdx]
//		member, ok := s.sequenceToData[seqNo]
//		if ok {
//			s.deleteBitMap(member.bitmapIdx)
//			//			delete(s.cDataToSequence, *member.Value)
//			//			delete(s.cSequenceToData, seqNo)
//			s.liveCount.Store(s.liveCount.Load() - 1)
//			return true
//		}
//	}
//	return false
//}
//
//func (s *AnyOrderedSet) cDeleteIdx(idx int) bool {
//	itemIdx := getNthAliveIndexUnderLock(s.bitmap, idx)
//	if itemIdx >= 0 {
//		seqNo := s.items[itemIdx]
//		member, ok := s.sequenceToData[seqNo]
//		if ok {
//			s.deleteBitMap(member.bitmapIdx)
//			delete(s.cDataToSequence, *member.Value)
//			delete(s.cSequenceToData, seqNo)
//			s.liveCount.Store(s.liveCount.Load() - 1)
//			return true
//		}
//	}
//	return false
//}
//
//func (s *AnyOrderedSet) isValAlive(elem T) bool {
//	if _, ok := s.indexMap[elem]; ok {
//		return true
//	}
//	return false
//}
//
//func (s *AnyOrderedSet) isSeqAlive(seqNo uint64) bool {
//	if member, ok := s.sequenceToData[seqNo]; ok {
//		if s.isIdxAlive(member.bitmapIdx) {
//			return true
//		}
//	}
//	return false
//}
//
//// func getBitMapIndex(idx int) uint64
//func (s *AnyOrderedSet) isIdxAlive(idx int) bool {
//	if idx < 0 || idx > len(s.bitmap) {
//		return false
//	}
//	uintIndex := idx / 64
//	bitOffset := idx % 64
//	bitMask := uint64(1) << bitOffset
//	return (s.bitmap[uintIndex] & bitMask) != 0
//}
//
//func (s *AnyOrderedSet) In(elem T) bool {
//	if _, ok := s.indexMap[elem]; ok {
//		return true
//	} else {
//		return false
//	}
//}
//
//func (s *AnyOrderedSet) Iter() iter.Seq[T] {
//	return func(yield func(T) bool) {
//		for _, seqNo := range s.items {
//			val := *s.sequenceToData[seqNo].Value
//			if s.isValAlive(val) {
//				if !(yield(val)) {
//					return
//				}
//			}
//		}
//	}
//}
//
//func (s *AnyOrderedSet) IterIndex() iter.Seq2[int, any] {
//	return func(yield func(int, any) bool) {
//		liveIdx := 0
//		for _, item := range getAliveIndicesUnderLock(s.bitmap) {
//			val := *s.sequenceToData[s.items[item]].Value
//			if s.isValAlive(val) {
//				if !(yield(item, val)) {
//					return
//				}
//				liveIdx++
//			}
//		}
//	}
//}
//
//func isAlive(bitmap []uint64, idx int) bool {
//	if idx < 0 || idx > len(bitmap) {
//		return false
//	}
//	uintIndex := idx / 64
//	bitOffset := idx % 64
//	bitMask := uint64(1) << bitOffset
//	return (bitmap[uintIndex] & bitMask) != 0
//}
