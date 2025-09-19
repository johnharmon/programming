package set

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"sync"
	"sync/atomic"
)

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
	logger          *ConcreteLogger
	rwLock          *sync.RWMutex
	pending         map[uint64]setOp[T]
	pendingLock     *sync.RWMutex
	opsCh           chan *setOp[T]
	replayCh        chan replayOp[T]
	tombstones      int
	maxItems        int
	encoder         setOpEncoder[T]
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
		Value: %v, 
		Idx: %d, 
		opType: %d, 
		seqNo: %d
		`,
		reflect.TypeOf(op.opVal),
		*op.opVal,
		op.opIdx,
		op.opType,
		op.seqNo)
}

type LogEntry struct {
	Message   string `json:"message"`
	Timestamp any    `json:"timestamp"`
}

type RawLogArgs struct {
	FormatMessage string
	FormatArgs    []any
	Timestamp     string
}

type FlushToken struct {
	Iteration  int
	HandledBy  string
	SentBy     string
	ReceivedBy string
	Values     map[string]any
}

type ConcreteLogger struct {
	ActiveBuffer      *bytes.Buffer
	FlushBuffer       *bytes.Buffer
	Out               io.Writer
	Mu                *sync.Mutex
	FlushMu           *sync.Mutex
	SwapMu            *sync.Mutex
	FlushSender       chan *FlushToken
	FlushReceiver     chan *FlushToken
	LogOutput         chan []byte
	LogEntryCh        chan *LogEntry
	RawLogCh          chan *RawLogArgs
	RunCh             chan *sync.WaitGroup
	Done              chan struct{}
	LogFileName       string
	LogEntryPool      *sync.Pool
	RawLogArgPool     *sync.Pool
	MessageBufferPool *sync.Pool
}
