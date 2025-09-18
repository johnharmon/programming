package set

import (
	"fmt"
	"io"
	"reflect"
	"sync"
	"sync/atomic"
)

type anySetOp struct {
	callback chan bool
	opVal    any
	opIdx    int
	opType   int
	seqNo    uint64
}

type anySetMember struct {
	bitmapIdx int
	Value     any
}

type anySetEntity struct {
	Value any
	ID    uint64
}

type anySeplayOp struct {
	opType int
	seqNo  uint64
	value  any
	opIdx  int
}

type AnyOrderedSet struct {
	dualWrite       bool
	bitmap          []uint64
	compacting      bool
	cBitmap         []uint64
	cItems          []uint64
	cDataToSequence map[any]uint64
	dataToSequence  map[any]uint64
	sequenceToData  map[uint64]anySetMember
	cSequenceToData map[uint64]anySetMember
	snapshotSeqNo   uint64
	opPool          *sync.Pool
	items           []uint64
	indexMap        map[any]int
	liveCount       atomic.Uint64
	cLiveCount      int
	seqNo           atomic.Uint64 // integer representing the most recent operation number so ops can be tagged with it
	lastAppliedSeq  uint64
	cLastAppliedSeq uint64
	rwLock          *sync.RWMutex
	pending         map[uint64]anySetOp
	opsCh           chan *anySetOp
	tombstones      int
	maxItems        int
	cAppendBitmap   func(int)
	appendBitmap    func(int)
	// cAppend func(
}

type anySetOpEncoder struct {
	writer io.Writer
}

func NewAnySetOpEncoder(writer io.Writer) anySetOpEncoder {
	return anySetOpEncoder{writer: writer}
}

func (s *anySetOpEncoder) Encode(op anySetOp) {
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
