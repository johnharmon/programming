package main

import (
	"bufio"
	"io"
	"os"
)

func (w *Window) LoadBuffer(in io.Reader) {
	if f, ok := in.(*os.File); ok {
		dispBuf := NewEmptyDisplayBuffer()
		dispBuf.Lines = make([][]byte, 1)
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			newLine := make([]byte, len(scanner.Bytes()))
			copy(newLine, scanner.Bytes())
			dispBuf.Lines = append(dispBuf.Lines, newLine)
		}
	}
}
