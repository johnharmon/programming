package main

import (
	"fmt"
)

func main() {
	for b := 0x20; b <= 0x7E; b++ {
		fmt.Printf("KeyActionTree[0x%X] = &KeyAction{Value: []byte{0x%X}, IsTerminal: true, Action: \"%s\", PrintRaw: false, FromPool: false}\n",
			b, b, string(b))
	}
}
