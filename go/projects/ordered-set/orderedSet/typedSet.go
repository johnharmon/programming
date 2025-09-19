package set

import (
	"fmt"
	"io"
	"iter"
	"math/bits"
	"os"
	"sync"
	"sync/atomic"
)

func NewOrderedSet[T comparable]() *OrderedSet[T] {
	s := OrderedSet[T]{}
	s.dataToSequence = make(map[T]uint64)
	s.cDataToSequence = make(map[T]uint64)
	s.sequenceToData = make(map[uint64]setMember[T])
	s.cSequenceToData = make(map[uint64]setMember[T])
	s.indexMap = make(map[T]int)
	s.pending = make(map[uint64]setOp[T])
	s.opsCh = make(chan *setOp[T], 10000)
	s.bitmap = make([]uint64, 1, 10000)
	s.cBitmap = make([]uint64, 1, 10000)
	s.items = make([]uint64, 0, 10000)
	s.cItems = make([]uint64, 1, 10000)
	s.opPool = &sync.Pool{}
	s.rwLock = &sync.RWMutex{}
	s.logger = NewConcreteLogger()
	s.logger.Init()

	s.pendingLock = &sync.RWMutex{}
	s.opPool.New = func() any {
		return &setOp[T]{}
	}
	s.encoder = NewSetOpEncoder[T](os.Stdout)
	s.appendBitmap = appendBitmapClosure(&s.bitmap)
	s.cAppendBitmap = appendBitmapClosure(&s.cBitmap)
	s.seqNo = atomic.Uint64{}
	s.seqNo.Store(1)
	go s.sequencer()
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

func getNthAliveIndexUnderLock(bitmap []uint64, target int) (realIndex int) { // only do this under lock, else it cannot guarantee index accuracy or item liveleness
	// Scans a bitmap for Nth alive item and returns the 'bitmap index' of that 1
	if len(bitmap) < 1 || target < 0 {
		return -1
	}
	remaining := target
	for idx := range bitmap {
		popCount := bits.OnesCount64(bitmap[idx])
		if remaining <= popCount {
			word := bitmap[idx]
			for i := 1; i < remaining; i++ {
				word &= word - 1
			}
			return idx*64 + bits.TrailingZeros64(word)
		}
		remaining -= popCount
	}
	return -1

	// onesCount := bits.OnesCount64(bitmap[idx])
	//			tmp := aliveIndex + onesCount
	//			if tmp <= target {
	//				aliveIndex += onesCount
	//				actualIndex += 64
	//			} else {
	//				found = true
	//				targetIdx = idx
	//				break
	//			}
	//		}
}

func (s *OrderedSet[T]) sequencer() {
	var applied bool
	for op := range s.opsCh {
		applied = false
		s.rwLock.Lock()
		seqNo := s.seqNo.Load() + 1
		s.lastAppliedSeq = seqNo - 1
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
			s.pendingLock.Lock()
			s.pending[seqNo] = *op
			s.pendingLock.Unlock()
		}
		s.seqNo.Store(seqNo)
		s.logger.Logln(op.String())
		// fmt.Fprintf(os.Stdout, "%s\n", op.String())
		if s.liveCount.Load() > 100 && (int(s.liveCount.Load()) < len(s.items)*2) {
			go s.compact()
		}
		s.rwLock.Unlock()
		// s.encoder.Encode(*op)
		op.callback <- applied
		op.opVal, op.callback, op.opIdx = nil, nil, 0
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

func (s *OrderedSet[T]) append(elem T, seqNo uint64) bool {
	if _, ok := s.dataToSequence[elem]; !ok {
		bitmapIdx := len(s.items)
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
		bitmapIdx := len(s.cItems)
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

func (s *OrderedSet[T]) GetIdx(idx int) (*setEntity[T], bool) {
	s.rwLock.Lock()
	realIdx := getNthAliveIndexUnderLock(s.bitmap, idx+1)
	if realIdx != -1 {
		if realIdx < len(s.items) {
			seqNo := s.items[realIdx]
			if member, ok := s.sequenceToData[seqNo]; ok {
				s.rwLock.Unlock()
				return &setEntity[T]{Value: *member.Value, ID: seqNo}, true
			}
			// fmt.Printf("Error: non-existent sequence number: %d given by s.items[%d]\n", seqNo, realIdx)
		}
		// fmt.Printf("Error, out of bounds index given by 'getNthAliveIndexUnderLock(%d)'\n", idx)
	}
	s.rwLock.Unlock()
	return nil, false
}

func (s *OrderedSet[T]) Len() int {
	return int(s.liveCount.Load())
}

func (s *OrderedSet[T]) DumpItems(out io.Writer) {
	for idx, item := range s.items {
		fmt.Fprintf(out, "Index: %d, SequenceNumber: %d, Alive: %t\n", idx, item, s.isSeqAlive(item))
		fmt.Println("")
	}
}

func (s *OrderedSet[T]) DumpBitMap(out io.Writer) {
	for idx, word := range s.bitmap {
		fmt.Printf("%d: %b\n", idx, word)
	}
	fmt.Println("")
}

func (s *OrderedSet[T]) DumpSequenceMap(out io.Writer) {
	for k, v := range s.sequenceToData {
		fmt.Fprintf(out, "sequenceNumber: %d, bitmapIdx: %d, value: %d\n", k, v.bitmapIdx, *v.Value)
	}
	fmt.Println("")
}

func (s *OrderedSet[T]) DumpDataMap(out io.Writer) {
	for k, v := range s.dataToSequence {
		fmt.Fprintf(out, "value: %d, sequenceNumber: %d\n", k, v)
	}
	fmt.Println("")
}

func (s *OrderedSet[T]) DumpMetadata(out io.Writer) {
	//"\x1b[0;0H\x1b[2JTotalItemCount: %d, LiveItemCount: %d, sequenceNumber: %d\r",
	fmt.Fprintf(out,
		"TotalItemCount: %d, LiveItemCount: %d, sequenceNumber: %d\n",
		len(s.items),
		s.liveCount.Load(),
		s.seqNo.Load())
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
	// fmt.Printf("appendBitMap(%d):::: idx: %d, bitOffset: %d\n", idx, idx, bitOffset)
	bitMask := uint64(1) << bitOffset
	s.bitmap[uintIndex] |= bitMask
}

func (s *OrderedSet[T]) copyItemsUnderLock(newItems *[]uint64, newSequenceToData map[uint64]setMember[T], newDataToSequence map[T]uint64, newBitmap *[]uint64) int {
	liveCount := 0
	appendBitmap := appendBitmapClosure(newBitmap)
	for _, idx := range getAliveIndicesUnderLock(s.bitmap) {
		itemSeqNo := s.items[idx]
		// s.logger.Logln("copyItemsUnderLock: itemSeqNo: %d", itemSeqNo)
		fmt.Fprintf(os.Stderr, "copyItemsUnderLock: itemSeqNo: %d\n", itemSeqNo)
		if liveCount > len(*newItems)-1 {
			*newItems = append(*newItems, uint64(0))
		}
		liveCount++
		(*newItems)[liveCount-1] = itemSeqNo
		seqNo := s.items[idx]
		sm := s.sequenceToData[seqNo]
		newSequenceToData[itemSeqNo] = setMember[T]{Value: sm.Value, bitmapIdx: liveCount - 1}
		sMember := newSequenceToData[itemSeqNo]
		val := sMember.Value
		s.logger.Logln("%d", itemSeqNo)
		message := fmt.Sprintf("Attempting to access value for seqNo: %d\n", itemSeqNo)
		fmt.Fprint(os.Stderr, message)
		newDataToSequence[*val] = itemSeqNo
		appendBitmap(liveCount - 1)
	}
	return liveCount
}

func (s *OrderedSet[T]) compact() {
	s.cBitmap = s.cBitmap[:0]
	s.cItems = s.cItems[:0]
	liveCount := int(s.liveCount.Load())
	fmt.Printf("Making new maps\n")
	s.cSequenceToData = make(map[uint64]setMember[T], liveCount+liveCount/3)
	s.cDataToSequence = make(map[T]uint64, liveCount+liveCount/3)

	s.rwLock.Lock()
	s.compacting = true
	s.pendingLock.Lock()
	if s.pending == nil {
		s.pending = make(map[uint64]setOp[T])
	} else {
		clear(s.pending)
	}
	s.pendingLock.Unlock()
	s.rwLock.Unlock()

	s.rwLock.RLock()
	s.snapshotSeqNo = s.seqNo.Load()
	liveCount = s.copyItemsUnderLock(&s.cItems, s.cSequenceToData, s.cDataToSequence, &s.cBitmap)
	s.cItems = make([]uint64, liveCount, int(float64(liveCount)*1.33))
	s.cLastAppliedSeq = s.snapshotSeqNo
	s.rwLock.RUnlock()

	for currentSeq := s.snapshotSeqNo + 1; (s.seqNo.Load() - currentSeq) > 0; currentSeq++ {
		s.pendingLock.Lock()
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
		s.pendingLock.Unlock()
	}
	s.rwLock.Lock()
	for currentSeq := s.cLastAppliedSeq + 1; (s.seqNo.Load() - currentSeq) > 0; currentSeq++ {
		s.pendingLock.Lock()
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
		s.pendingLock.Unlock()
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
	op.opVal = &elem
	op.opIdx = -1
	op.opType = OP_DELETE
	op.callback = make(chan bool, 1)
	s.opsCh <- op
	return <-op.callback
}

func (s *OrderedSet[T]) delete(elem T) bool {
	if seqNo, ok := s.dataToSequence[elem]; ok {
		if member, ok2 := s.sequenceToData[seqNo]; ok2 {
			s.deleteBitMap(member.bitmapIdx)
			delete(s.dataToSequence, elem)
			delete(s.sequenceToData, seqNo)
			s.liveCount.Store(s.liveCount.Load() - 1)
			s.logger.Logln("Deleted sequence number: %d, with value: %d, flippeed bitmapIdx: %d", seqNo, elem, member.bitmapIdx)
			return true
		}
	}
	return false
}

func (s *OrderedSet[T]) cDelete(elem T) bool {
	if seqNo, ok := s.cDataToSequence[elem]; ok {
		if member, ok2 := s.cSequenceToData[seqNo]; ok2 {
			s.deleteBitMap(member.bitmapIdx)
			delete(s.cDataToSequence, elem)
			delete(s.cSequenceToData, seqNo)
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
	// fmt.Printf("isSeqAlive(%d)\n", seqNo)
	if member, ok := s.sequenceToData[seqNo]; ok {
		// fmt.Printf("isIdxAlive(%d)\n", member.bitmapIdx)
		if s.isIdxAlive(member.bitmapIdx) {
			return true
		}
	}
	return false
}

// func getBitMapIndex(idx int) uint64
func (s *OrderedSet[T]) isIdxAlive(idx int) bool {
	if idx < 0 || idx > len(s.items) {
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
