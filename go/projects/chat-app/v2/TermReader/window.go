package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"
	"time"

	"golang.org/x/term"
)

type Window struct { // Represents a sliding into its backing buffer of Window.Buf as well as the space it takes up in the terminal window
	TermTopLine   int
	BufTopLine    int
	StartIndex    int
	StartLine     int
	Buf           *DisplayBuffer
	Height        int
	Width         int
	StartCol      int
	CursorLine    int
	CursorCol     int
	EndIndex      int
	RawStartIndex int
	RawEndIndex   int
	Out           io.Writer
}

func (w Window) Size() int {
	return w.Height
}

func NewWindow(line int, column int, height int, width int) (w *Window) {
	w = &Window{}
	w.StartLine = line
	w.StartCol = column
	w.Height = height
	w.Width = width
	w.Buf = &DisplayBuffer{}
	return w
}

func (w Window) MoveCursorToPosition(line int, col int) {
	fmt.Fprintf(w.Out, "\x1b[%d;%dH", line, col)
}

func (w Window) Render(out io.Writer) {
	fmt.Fprintf(out, "\x1b[%d;0H", w.TermTopLine)
	for i := 0; i <= w.Height; i++ {
		targetLine := i + w.BufTopLine
		if targetLine < len(w.Buf.Lines) {
			out.Write(w.Buf.Lines[i+w.BufTopLine])
		} else {
			break
		}
	}
}

func (w *Window) Scroll(scrollVector int) {
}

func (cell *Cell) ScrollWindowLine(scrollVector int) {
	currentLength := cell.ActiveLineLength
	cell.Log("=========INPUT BREAK========")
	cell.Log("Called ScrollLIne with a scrollVector of %d", scrollVector)
	cell.Log("Current line index: %d", cell.ActiveLineIdx)
	cell.Log("Current line length: %d", currentLength)
	nextLineIndex := cell.GetIncrActiveLine(scrollVector)
	cell.Log("Next line index: %d", nextLineIndex)
	nextLineLength := cell.GetLineLen(nextLineIndex)
	cell.Log("Next line length: %d", nextLineLength)
	if currentLength <= 0 && nextLineLength <= 0 {
		cell.Log("Line state not compabible with scrolling, ignoring...")
		return
	} else {
		cell.Log("Scrolling %d", scrollVector)
		cell.IncrActiveLine(scrollVector)
	}
}
