package main

import (
	"bufio"
	"io"
	"os"
)

func (w *Window) LoadBuffer(in io.Reader) {
	if f, ok := in.(*os.File); ok {
		dispBuf := NewEmptyDisplayBuffer()
		dispBuf.Lines = make([][]byte, 0)
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			newLine := make([]byte, len(scanner.Bytes()))
			copy(newLine, scanner.Bytes())
			dispBuf.Lines = append(dispBuf.Lines, newLine)
		}
		if len(dispBuf.Lines) == 0 {
			dispBuf.Lines = MakeNewLines(1, 80)
		}
		w.Buf = dispBuf
		w.CursorLine = 0
	}
}

func (w *Window) LoadNewEmptyBuffer() {
	w.Buf.Lines = MakeNewLines(10, 256)
	w.CursorCol, w.CursorLine = 1, 0
	w.TermTopLine = 1
}

func (w *Window) WriteBuffer(out io.Writer) {
	for _, line := range w.Buf.Lines {
		out.Write(append(line, '\n'))
	}
}
