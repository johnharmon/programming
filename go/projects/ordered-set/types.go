package orderedset

import (
	"fmt"
	"io"
	"reflect"
	"sync"
	"sync/atomic"
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

type setOpEncoder[T comparable] struct {
	writer io.Writer
}

func NewSetOpEncoder[T comparable](writer io.Writer) setOpEncoder[T] {
	return setOpEncoder[T]{writer: writer}
}

func (s *setOpEncoder[T]) Encode(op setOp[T]) {
	fmt.Fprintf(s.writer,
		`Type: %v, 
		Value: %s, 
		Idx: %d, 
		opType: %d, 
		seqNo: %d\n`,
		reflect.TypeOf(op.opVal),
		op.opVal,
		op.opIdx,
		op.opType,
		op.seqNo)
}
