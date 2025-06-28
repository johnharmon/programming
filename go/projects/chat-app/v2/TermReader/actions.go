package main

import (
	"fmt"
	"log"
)

var KeyActionTree map[byte]*KeyAction

func NewKeyAction(terminal bool, action string, raw bool, value ...byte) (sq *KeyAction) {
	sq = &KeyAction{Children: make(map[byte]*KeyAction), Value: value, IsTerminal: terminal, Action: action, PrintRaw: raw, FromPool: false}
	return sq
}

func NewKeyActionFromPool(terminal bool, action string, raw bool, value ...byte) (sq *KeyAction) {
	sq = &KeyAction{Children: make(map[byte]*KeyAction), Value: make([]byte, len(value), 8), IsTerminal: terminal, Action: action, PrintRaw: raw, FromPool: true}
	copy(sq.Value, value)
	return sq
}

func InitializeArrowKeys() error {
	if _, ok := KeyActionTree[0x1b]; !ok {
		return fmt.Errorf("error initializing the arrow keys; start of escape sequence not set")
	}
	escapeSequences := KeyActionTree[0x1b]
	escapeSequences.Children[0x5b] = NewKeyAction(false, "ArrowKeyPrefix", false, 0x5b)
	arrowKeyParent := escapeSequences.Children[0x5b]
	arrowKeyParent.Children[0x41] = NewKeyAction(true, "ArrowUp", false, []byte{0x1b, 0x5b, 0x41}...)
	arrowKeyParent.Children[0x42] = NewKeyAction(true, "ArrowDown", false, []byte{0x1b, 0x5b, 0x42}...)
	arrowKeyParent.Children[0x43] = NewKeyAction(true, "ArrowRight", false, []byte{0x1b, 0x5b, 0x43}...)
	arrowKeyParent.Children[0x44] = NewKeyAction(true, "ArrowLeft", false, []byte{0x1b, 0x5b, 0x44}...)
	return nil
}

func InitializeControlCodes() error {
	if _, ok := KeyActionTree[0x1b]; !ok {
		return fmt.Errorf("error initializing the arrow keys; start of escape sequence not set")
	}
	KeyActionTree[0x7F] = NewKeyAction(true, "Delete", false, 0x7F)
	KeyActionTree[0x08] = NewKeyAction(true, "Backspace", false, 0x08)
	KeyActionTree[0x0D] = NewKeyAction(true, "Enter", false, 0x0D)
	return nil
}

func init() {
	KeyActionTree = make(map[byte]*KeyAction)
	KeyActionTree[0x1b] = NewKeyAction(false, "Escape", false, 0x1b)
	err := InitializeArrowKeys()
	if err != nil {
		log.Fatal(err.Error())
	}
	err = InitializeControlCodes()
	if err != nil {
		log.Fatal(err.Error())
	}
}

func ValidateSequence(seq []byte) (sq *KeyAction) {
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

func (ka *KeyAction) String() string {
	return fmt.Sprintf("KeyAction: { \"Value\": %X, \"FromPool\": %t, \"Action\": \"%s\"}", ka.Value, ka.FromPool, ka.Action)
}
