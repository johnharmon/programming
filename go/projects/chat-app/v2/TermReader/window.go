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
	EventChan     chan *SequenceNode
	RawEventChan  chan []byte
}

type LogConfig struct {
	File *os.File
	Link *os.File
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

type KeyAction struct {
	None int
}

type InputSequence struct {
	Bytes  []byte
	Action *KeyAction
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

func (w *Window) Listen() {
	input := []byte{}
	for {
		input = <-w.RawEventChan
		w.Buf.Write(input)
	}
}

func (w *Window) Scroll(scrollVector int) {
}

/*
func (w *Window) ScrollWindowLine(scrollVector int) {
	currentLength := w.ActiveLineLength
	w.Log("=========INPUT BREAK========")
	w.Log("Called ScrollLIne with a scrollVector of %d", scrollVector)
	w.Log("Current line index: %d", w.ActiveLineIdx)
	w.Log("Current line length: %d", currentLength)
	nextLineIndex := w.GetIncrActiveLine(scrollVector)
	w.Log("Next line index: %d", nextLineIndex)
	nextLineLength := w.GetLineLen(nextLineIndex)
	w.Log("Next line length: %d", nextLineLength)
	if currentLength <= 0 && nextLineLength <= 0 {
		w.Log("Line state not compabible with scrolling, ignoring...")
		return
	} else {
		w.Log("Scrolling %d", scrollVector)
		w.IncrActiveLine(scrollVector)
	}
}
*/

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
	var (
		res []byte
		ka  *SequenceNode
	)
	fd := int(os.Stdin.Fd())
	oldState, err := term.MakeRaw(fd)
	if err != nil {
		panic(err)
	}
	closer := make(chan interface{})
	go Cleanup(closer, fd, oldState, mc.LogConfig)
	buf := make([]byte, 1)
	byteHandler := MakeByteHandler(closer, mc.In)
	for {
		nb, err := mc.In.Read(buf)
		if err != nil {
			panic(err)
		}
		if nb > 0 {
			b := buf[0]
			// res = HandleByte(b, closer, mc.In) // this should return the final coerced byte or []byte that the window will be responsible for processing
			ka = byteHandler(b)
			mc.State.ActiveWindow.EventChan <- ka
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
			return NewSequenceNode(0x00, true, "Print", true)
		},
	}
}

func CoerceInputToAction(b []byte) *SequenceNode {
	if len(b) == 1 {
		return KeyActionTree[b[0]]
	}
	return ValidateSequence(b)
}

func MakeByteHandler(ch chan interface{}, in io.Reader) func(byte) *SequenceNode { // returns a byte handling function that will reuse an input buffer so re-allocation does not happen on every byte handled by the main loop
	res := make([]byte, 1, 8)
	var seqN *SequenceNode
	sp := MakeSequencePool() // Create a pool of *SequenceNode references for dispatching normal ascii printable characters on the event channel for the window (trying to avoid as much re-allocation as possible
	return func(b byte) *SequenceNode {
		if b == 3 {
			ch <- struct{}{}
		} else if b == 13 {
			os.Exit(0)
		} else if b >= 0x20 && b <= 0x7E {
			seqN = sp.Get().(*SequenceNode)
			seqN.Value = b
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
