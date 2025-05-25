package main

import (
	"fmt"
	"log"
)

type SequenceNode struct {
	Children   map[byte]*SequenceNode
	Value      byte
	IsTerminal bool
	Action     string
	PrintRaw   bool
}

var KeyActionTree map[byte]*SequenceNode

func NewSequenceNode(value byte, terminal bool, action string, raw bool) (sq *SequenceNode) {
	sq = &SequenceNode{Children: make(map[byte]*SequenceNode), Value: value, IsTerminal: terminal, Action: action, PrintRaw: raw}
	return sq
}

func InitializeArrowKeys() error {
	if _, ok := KeyActionTree[0x1b]; !ok {
		return fmt.Errorf("error initializing the arrow keys; start of escape sequence not set")
	}
	escapeSequences := KeyActionTree[0x1b]
	escapeSequences.Children[0x5b] = NewSequenceNode(0x5b, false, "ArrowKeyPrefix", false)
	arrowKeyParent := escapeSequences.Children[0x5b]
	arrowKeyParent.Children[0x41] = NewSequenceNode(0x41, true, "ArrowUp", true)
	arrowKeyParent.Children[0x42] = NewSequenceNode(0x42, true, "ArroDown", true)
	arrowKeyParent.Children[0x43] = NewSequenceNode(0x43, true, "ArrowRight", true)
	arrowKeyParent.Children[0x44] = NewSequenceNode(0x44, true, "ArrowLeft", true)
	return nil
}

func InitializeControlCodes() error {
	if _, ok := KeyActionTree[0x1b]; !ok {
		return fmt.Errorf("error initializing the arrow keys; start of escape sequence not set")
	}
	KeyActionTree[0x7F] = NewSequenceNode(0x7F, true, "Delete", false)
	KeyActionTree[0x08] = NewSequenceNode(0x08, true, "Backspace", false)
	KeyActionTree[0x0A] = NewSequenceNode(0x0A, true, "Enter", false)
	return nil
}

func init() {
	KeyActionTree = make(map[byte]*SequenceNode)
	KeyActionTree[0x1b] = NewSequenceNode(0x1b, false, "TerminalEscape", false)
	err := InitializeArrowKeys()
	if err != nil {
		log.Fatal(err.Error())
	}
	err = InitializeControlCodes()
	if err != nil {
		log.Fatal(err.Error())
	}
}

func ValidateSequence(seq []byte) (sq *SequenceNode) {
	seqLen := len(seq)
	if seqLen == 0 {
		return nil
	}
	tmpNode, ok := KeyActionTree[seq[0]]
	if ok {
		for i := 1; i < seqLen; i++ {
			tmpNode = tmpNode.Children[seq[i]]
			if tmpNode == nil {
				return nil
			}
		}
		return tmpNode
	}
	return nil
}
