package set

import (
	"fmt"
)

type setOp[T comparable] struct {
	callback chan bool
	opVal    *T
	opIdx    int
	opType   int
	seqNo    uint64
}

func (s setOp[T]) String() string {
	var opType string
	switch s.opType {
	case OP_APPEND:
		opType = "append"
	case OP_DELETE:
		opType = "delete"
	case OP_DELETE_IDX:
		opType = "deleteIdx"
	}
	return fmt.Sprintf("opIdx: %d, opType: %s, seqNo: %d, opVal: %d", s.opIdx, opType, s.seqNo, *s.opVal)
}
