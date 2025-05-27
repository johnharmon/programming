package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"
	"sync"
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
	EventChan     chan *KeyAction
	RawEventChan  chan []byte
}

type LogConfig struct {
	File *os.File
	Link *os.File
}

type PooledKeyAction struct {
	KA       *KeyAction
	FromPool bool
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

type State struct {
	Windows      []*Window
	ActiveWindow *Window
}

type MainConfig struct {
	LogConfig *LogConfig
	In        io.Reader
	Out       io.Writer
	State     State
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

func (w *Window) WriteRaw(b []byte) {
	w.Buf.Lines[w.Buf.ActiveLine] = InsertAt(w.Buf.Lines[w.Buf.ActiveLine], b, w.CursorCol)
}

func (w *Window) IncrCursorCol(incr int) {
	lLen := len(w.Buf.Lines[w.Buf.ActiveLine])
	newPos := w.CursorCol + incr
	if newPos < 0 {
		newPos = 0
	} else if newPos > lLen {
		newPos = lLen
	}
	w.CursorCol = newPos
}

func (w *Window) IncrCursorLine(vec int) {
	nextLine := w.Buf.ActiveLine + vec
	if nextLine > 0 && nextLine < len(w.Buf.Lines) {
		w.Buf.ActiveLine = nextLine
	}
}

func (w *Window) MakeNewLines(count int) [][]byte {
	newLines := make([][]byte, count, count)
	for i := range newLines {
		newLines[i] = make([]byte, 0, 4096)
	}
	return newLines
}

func (w *Window) Listen() {
	redrawHandler := w.MakeRedrawHandler()
	var ka *KeyAction
	for {
		ka = <-w.EventChan
		if ka.PrintRaw && len(ka.Value) == 1 {
			w.WriteRaw(ka.Value)
		} else {
			switch ka.Action {
			case "Backspace":
				w.Buf.Lines[w.Buf.ActiveLine] = DeleteByteAt(w.Buf.Lines[w.Buf.ActiveLine], w.CursorCol-1)
				w.IncrCursorCol(-1)
			case "Delete":
				w.Buf.Lines[w.Buf.ActiveLine] = DeleteByteAt(w.Buf.Lines[w.Buf.ActiveLine], w.CursorCol)
				w.IncrCursorCol(-1)
				w.Buf.Write(ka.Value)
			case "RightArrow":
				w.IncrCursorCol(1)
				w.Buf.Write(ka.Value)
			case "UpArrow":
				w.IncrCursorLine(-1)
				w.Buf.Write(ka.Value)
			case "DownArrow":
				w.IncrCursorLine(1)
				w.Buf.Write(ka.Value)
			case "Enter":
				newLine := w.MakeNewLines(1)
				w.WriteRaw([]byte("\r\n"))
				w.Buf.Lines = InsertLineAt(w.Buf.Lines, newLine, w.CursorLine)
				w.Redraw(redrawHandler)
			}
		}

	}
}

func (w *Window) Redraw(handler func() []int) {
	linesToRedraw := handler()
	lastIndex := 0
	w.MoveCursorToPosition(w.TermTopLine, 0)
	for _, lineNum := range linesToRedraw {
		w.MoveCursorToPosition(w.TermTopLine+lineNum, 0)
		RedrawLine(w.Buf.Lines[lineNum+w.Buf.TopLine])
		lastIndex++
	}
}

func (w *Window) MakeRedrawHandler() func() []int {
	redrawIndicies := make([]int, 0, w.Height)
	return func() []int {
		clear(redrawIndicies)
		redrawIndicies = redrawIndicies[:0]
		for i := w.Buf.TopLine; i <= w.Buf.TopLine+w.Buf.Height; i++ {
			if !BufCmp(w.Buf.Lines[i], w.Buf.DisplayedLines[i-w.Buf.TopLine]) {
				redrawIndicies = append(redrawIndicies, i-w.Buf.TopLine)
			}
		}
		return redrawIndicies
	}
}

func (w *Window) MarkForReddraw() {
	for i := w.Buf.TopLine; i <= w.Buf.TopLine+w.Buf.Height; i++ {
	}
}

func (db *DisplayBuffer) UpdateDisplayedLines(start int, end int) {
	db.DisplayedLines = db.DisplayedLines[:0]
	for i := start; i <= end; i++ {
		db.DisplayedLines = append(db.DisplayedLines, db.Lines[i])
	}
}

func (w *Window) Scroll(scrollVector int) {
}

func Cleanup(closer chan interface{}, fd int, oldState *term.State, logConfig *LogConfig) {
	<-closer
	fmt.Println("\n\rRestoring old state")
	term.Restore(fd, oldState)
	logConfig.File.Close()
	os.Remove(logConfig.File.Name())
	os.Remove(logConfig.Link.Name())
	os.Exit(0)
}

func MainEventHandler(mc *MainConfig) {
	var ka *KeyAction
	fd := int(os.Stdin.Fd())
	oldState, err := term.MakeRaw(fd)
	if err != nil {
		panic(err)
	}
	closer := make(chan interface{})
	go Cleanup(closer, fd, oldState, mc.LogConfig)
	buf := make([]byte, 1)
	sp := MakeSequencePool() // Create a pool of *SequenceNode references for dispatching normal ascii printable characters on the event channel for the window (trying to avoid as much re-allocation as possible
	byteHandler := MakeByteHandler(closer, mc.In, sp)
	for {
		nb, err := mc.In.Read(buf)
		if err != nil {
			panic(err)
		}
		if nb > 0 {
			b := buf[0]
			// res = HandleByte(b, closer, mc.In) // this should return the final coerced byte or []byte that the window will be responsible for processing
			ka = byteHandler(b)
			if ka != nil {
				mc.State.ActiveWindow.EventChan <- ka
			}
			// mc.State.ActiveWindow.RawEventChan <- res
		}
	}
}

func (mc *MainConfig) CoerceInput(b byte) (inputSeq []byte) { // Will coerce input to an actionable sequence, will possibly read more bytes from the main input source, bypassing the main event loop
	// HandleModSequence(cell, modSeq)
	return []byte{}
}

func MakeSequencePool() *sync.Pool {
	return &sync.Pool{
		New: func() any {
			return NewKeyActionFromPool(true, "Print", true, 0x00)
		},
	}
}

func CoerceInputToAction(b []byte) *KeyAction {
	if len(b) == 1 {
		return KeyActionTree[b[0]]
	}
	return ValidateSequence(b)
}

func MakeByteHandler(ch chan interface{}, in io.Reader, sp *sync.Pool) func(byte) *KeyAction { // returns a byte handling function that will reuse an input buffer so re-allocation does not happen on every byte handled by the main loop
	res := make([]byte, 1, 8)
	var seqN *KeyAction
	sp := MakeSequencePool() // Create a pool of *SequenceNode references for dispatching normal ascii printable characters on the event channel for the window (trying to avoid as much re-allocation as possible
	return func(b byte) *KeyAction {
		if b == 3 {
			ch <- struct{}{}
		} else if b == 13 {
			os.Exit(0)
		} else if b >= 0x20 && b <= 0x7E {
			seqN = sp.Get().(*KeyAction)
			seqN.Value[0] = b
			seqN.Value = seqN.Value[0:1]
			return seqN
		}
		res[0] = b
		res = ParseByte(b, res, in)
		defer clear(res)
		return CoerceInputToAction(res)
	}
}

func ParseByte(b byte, result []byte, in io.Reader) []byte { // Should handle initial detection for multi-byte sequences, if a single byte sequence then just return the byte as a slide
	result[0] = b
	if b == 0x1b {
		n, _ := ReadMultiByteSequence(result, in, time.Millisecond*25)
		result = result[:1+n]
	}
	return result
}

func ReadMultiByteSequence(buf []byte, input io.Reader, timeout time.Duration) (n int, err error) { // will read a multi-byte sequence into buf, respecting any existing elements
	bufLen := len(buf)
	bufCap := cap(buf)
	deadline := time.Now().Add(timeout)
	if f, ok := input.(*os.File); ok {
		f.SetReadDeadline(deadline)
	}
	defer ClearReadDeadline(input)
	n, err = input.Read(buf[bufLen:bufCap])
	if err != nil {
		return 0, err
	}
	return n, nil
}

func HandleByte(b byte, ch chan interface{}, in io.Reader) (res []byte) {
	res = make([]byte, 1, 8)
	if b == 3 {
		ch <- struct{}{}
	} else if b == 13 {
		//		cell.ScrollLine(1)
		//		cell.SetCursorPositionFromActiveLine()
		//		cell.DisplayActiveLine()
	} else {
		res[0] = b
		ParseByte(b, res, in)

	}
	return res
}
