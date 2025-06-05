package main

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"golang.org/x/term"
)

var (
	TERM_CLEAR_LINE   = []byte{0x1b, '[', '2', 'K'}
	TERM_CLEAR_SCREEN = []byte{0x1b, '[', '2', 'J'}
)

func (w Window) Size() int {
	return w.Height
}

func (w Window) GetCursorScreenLine() int {
	return w.CursorLine - w.BufTopLine + w.TermTopLine
}

func NewEmptyDisplayBuffer() *DisplayBuffer {
	db := &DisplayBuffer{
		Lines: make([][]byte, 10),
	}
	db.AllocateLines(256)
	return db
}

func NewWindow(line int, column int, height int, width int) (w *Window) {
	w = &Window{}
	w.TermTopLine = line
	w.StartCol = column
	w.Height = height
	w.Width = width
	w.BufTopLine = 0
	w.TermTopLine = 1
	w.Out = os.Stdout
	w.EventChan = make(chan *KeyAction, 10)
	w.KeyActionReturner = make(chan *KeyAction, 10)
	w.Buf = NewEmptyDisplayBuffer()
	return w
}

func (w Window) MoveCursorToPosition(line int, col int) {
	w.Logger.Logln("Moving Cursor to: %d,%d", line, col)
	fmt.Fprintf(w.Out, "\x1b[%d;%dH", line, col)
}

func (w *Window) RunKeyActionReturner(sp *sync.Pool) {
	for {
		sp.Put(<-w.KeyActionReturner)
	}
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
	w.Logger.Logln("Writing %b to %d", b, w.Buf.ActiveLine)
	w.Buf.Lines[w.Buf.ActiveLine] = InsertAt(w.Buf.Lines[w.Buf.ActiveLine], b, w.CursorCol)
}

func (w *Window) IncrCursorCol(incr int) {
	lLen := len(w.Buf.Lines[w.Buf.ActiveLine])
	newPos := w.CursorCol + incr
	if newPos < 0 {
		newPos = 0
		w.CursorCol = newPos
		w.DesiredCursorCol = newPos
	} else if newPos < lLen {
		w.CursorCol = newPos
		w.DesiredCursorCol = newPos
	} else if newPos >= -lLen && incr < 0 {
		w.CursorCol = lLen - 1
		w.DesiredCursorCol = lLen - 1
	} else if newPos >= lLen && incr > 0 {
		w.CursorCol = lLen
		w.DesiredCursorCol = lLen
	} else {
		w.CursorCol = lLen
	}
}

func (w *Window) GetDisplayCursorPosition() int {
	if len(w.Buf.Lines[w.Buf.ActiveLine]) < w.DesiredCursorCol {
		w.Logger.Logln("Setting cursor display position to: %d", len(w.Buf.Lines[w.Buf.ActiveLine]))
		return len(w.Buf.Lines[w.Buf.ActiveLine])
	} else {
		w.Logger.Logln("Desired Cursor position is compatible with the active line")
		return w.DesiredCursorCol
	}
}

// func (w *Window) GetPosition

func (w *Window) HandleArrowDown() (col int) {
	w.IncrCursorLine(1)
	col = w.GetDisplayCursorPosition()
	return col
}

func (w *Window) IncrCursorLine(vec int) {
	nextLine := w.CursorLine + vec
	if nextLine >= 0 && nextLine < len(w.Buf.Lines) {
		w.Buf.ActiveLine = nextLine
		w.CursorLine = nextLine
	}
}

/*
func (w *Window) SetCursorCol() {
	if w.DesiredCursorCol < len(w.Buf.Lines(w.
}
*/

func (w *Window) MakeNewLines(count int) [][]byte {
	newLines := make([][]byte, count, count)
	for i := range newLines {
		newLines[i] = make([]byte, 0, 4096)
	}
	return newLines
}

func (w *Window) RedrawLine(ln int) {
	screenLine := ln - w.BufTopLine
	if screenLine >= 0 && screenLine < w.BufTopLine+w.Height {
		w.MoveCursorToPosition(screenLine+w.TermTopLine, 0)
		w.Logger.Logln("Writing content to line #%d:", screenLine)
		w.Logger.Logln("%s", w.Buf.Lines[ln])
		w.Out.Write(TERM_CLEAR_LINE)
		w.Out.Write(w.Buf.Lines[ln])
	}
}

func (w *Window) NewBuffer() {
	w.Buf = &DisplayBuffer{}
	w.Buf.Lines = make([][]byte, 1, 100)
	for i := 0; i < len(w.Buf.Lines); i++ {
		w.Buf.Lines[i] = make([]byte, 0, 256)
	}
	w.Buf.Size = len(w.Buf.Lines)
	w.Buf.ActiveLine = 0
	w.CursorLine = 0
	w.CursorCol = 0
	w.Redraw(w.MakeRedrawHandler())
}

func (w *Window) Listen() {
	// redrawHandler := w.MakeRedrawHandler()
	gl := GlobalLogger
	w.Logger = gl
	var ka *KeyAction
	for {
		ka = <-w.EventChan
		gl.Logln("Window received *KeyAction: %s", ka.String())
		if ka.PrintRaw && len(ka.Value) == 1 {
			gl.Logln("Raw write triggered for %s", ka.String())
			w.WriteRaw(ka.Value)
			w.IncrCursorCol(1)
			w.RedrawLine(w.Buf.ActiveLine)
			w.MoveCursorToDisplayPosition()
			w.KeyActionReturner <- ka
		} else {
			switch ka.Action {
			case "Backspace":
				w.Buf.Lines[w.Buf.ActiveLine] = DeleteByteAt(w.Buf.Lines[w.Buf.ActiveLine], w.CursorCol-1)
				w.RedrawLine(w.CursorLine)
				w.IncrCursorCol(-1)
			case "Delete":
				w.Buf.Lines[w.Buf.ActiveLine] = DeleteByteAt(w.Buf.Lines[w.Buf.ActiveLine], w.CursorCol)
				w.RedrawLine(w.CursorLine)
				w.IncrCursorCol(-1)
				w.Buf.Write(ka.Value)
			case "ArrowRight":
				w.IncrCursorCol(1)
				// w.WriteRaw(ka.Value)
				w.Out.Write(ka.Value)
			case "ArrowLeft":
				w.IncrCursorCol(-1)
				// w.WriteRaw(ka.Value)
				w.Out.Write(ka.Value)
			case "ArrowUp":
				w.IncrCursorLine(-1)
				// w.SetCursorCol()
				// w.WriteRaw(ka.Value)
				w.Out.Write(ka.Value)
			case "ArrowDown":
				col := w.HandleArrowDown()
				// w.WriteRaw(ka.Value)
				w.MoveCursorToPosition(w.CursorLine+1, col)
			case "Enter":
				newLine := w.MakeNewLines(1)
				w.Logger.Logln("Enter detected")
				// w.WriteRaw([]byte("\r\n"))
				w.Logger.Logln("Inserting new line at index %d", w.CursorLine+1)
				w.Buf.Lines = InsertLineAt(w.Buf.Lines, newLine, w.CursorLine+1)
				w.Logger.Logln("New Byte buffer: %b", w.Buf.Lines)
				w.IncrCursorLine(1)
				w.CursorCol = 0
				// w.Redraw(redrawHandler)
				w.RedrawAllLines()
				w.MoveCursorToDisplayPosition()
			}
		}
		w.DisplayStatusLine()
		w.MoveCursorToDisplayPosition()
	}
}

func (w *Window) MoveCursorToDisplayPosition() {
	cursorDisplayLine := (w.CursorLine - w.BufTopLine) + 1
	cursorCol := w.GetDisplayCursorPosition()
	w.MoveCursorToPosition(cursorDisplayLine, cursorCol)
}

func (w *Window) Redraw(handler func() []int) {
	w.Logger.Logln("Redrawing window")
	linesToRedraw := handler()
	w.Logger.Logln("Lines to redraw: %s", linesToRedraw)
	lastIndex := 0
	w.MoveCursorToPosition(w.TermTopLine, 0)
	for _, lineNum := range linesToRedraw {
		w.MoveCursorToPosition(w.TermTopLine+lineNum, 0)
		w.RedrawLine(lineNum + w.BufTopLine)
		lastIndex++
	}
	// w.DisplayStatusLine()
}

func (w *Window) RedrawAllLines() {
	w.Logger.Logln("Forcing full redraw of all lines")
	w.MoveCursorToPosition(w.TermTopLine, 0)
	var lineLimit int
	linesLeftInBuffer := len(w.Buf.Lines) - w.BufTopLine - 1
	if linesLeftInBuffer < w.Height {
		lineLimit = linesLeftInBuffer
	} else {
		lineLimit = w.Height
	}
	for i := w.BufTopLine; i < w.BufTopLine+lineLimit; i++ {
		w.RedrawLine(i)
	}
}

func (w *Window) DisplayStatusLine() {
	termStatusLineNum := w.TermTopLine + w.Height + 1
	w.MoveCursorToPosition(termStatusLineNum, 0)
	w.Out.Write(TERM_CLEAR_LINE)
	fmt.Fprintf(
		w.Out,
		"CursorLine: %d | CursorColumn: %d | DesiredCursorColumn: %d | TermLine: %d | LineLength: %d | BufferLength: %d | WindowHeight: %d",
		w.CursorLine,
		w.CursorCol,
		w.DesiredCursorCol,
		w.TermTopLine+(w.CursorLine-w.BufTopLine),
		len(w.Buf.Lines[w.Buf.ActiveLine]),
		len(w.Buf.Lines),
		w.Height)
}

func (w *Window) MakeRedrawHandler() func() []int { // makes a handler that will track indicies to redraw (reuses slice)
	redrawIndicies := make([]int, 0, w.Height)
	return func() []int {
		GlobalLogger.Logln("Redraw Handler invoked")
		clear(redrawIndicies)
		GlobalLogger.Logln("Redraw Indicies slice cleared")
		redrawIndicies = redrawIndicies[:]
		GlobalLogger.Logln("Redraw Indicies slice zeroed")
		var loopUpperBound int
		remainingBufLen := len(w.Buf.Lines) - w.BufTopLine
		if remainingBufLen > w.Height {
			loopUpperBound = w.Height
		} else {
			loopUpperBound = remainingBufLen
		}
		for i := w.BufTopLine; i < loopUpperBound; i++ {
			GlobalLogger.Logln("Redraw Handler loop: Iteration: %d", i)
			// if !BufCmp(w.Buf.Lines[i], w.Buf.DisplayedLines[i-w.Buf.TopLine]) {
			redrawIndicies = append(redrawIndicies, i-w.BufTopLine)
			//}
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

func ReturnKeyActionsToPool(p *sync.Pool, returner chan *KeyAction) {
	var ka *KeyAction
	for {
		ka = <-returner
		if ka.FromPool {
			p.Put(ka)
		}
	}
}

func MainEventHandler(mc *MainConfig) {
	gl := GlobalLogger
	var ka *KeyAction
	fd := int(os.Stdin.Fd())
	gl.Logln("Setting terminal to raw mode")
	oldState, err := term.MakeRaw(fd)
	if err != nil {
		panic(err)
	}
	gl.Logln("Making closer channel")
	closer := make(chan interface{})
	gl.Logln("Making *KeyAction pool")
	sp := MakeKeyActionPool() // Create a pool of *SequenceNode references for dispatching normal ascii printable characters on the event channel for the window (trying to avoid as much re-allocation as possible
	go mc.State.ActiveWindow.RunKeyActionReturner(sp)
	gl.Logln("Making *KeyAction return channel for the pool")
	keyActionReturner := make(chan *KeyAction, 1000)
	gl.Logln("Spinning off goroutine for returning *KeyActions to the pool")
	go ReturnKeyActionsToPool(sp, keyActionReturner)
	gl.Logln("Spinning off cleanup goroutine")
	go Cleanup(closer, fd, oldState, mc.LogConfig)
	buf := make([]byte, 1)
	gl.Logln("Creating byte handler from closure")
	byteHandler := MakeByteHandler(closer, mc.In, sp)
	gl.Logln("Entering main event loop")
	for {
		nb, err := mc.In.Read(buf)
		if err != nil {
			gl.Logln("error encountered while reading from stdin")
			panic(err)
		}
		if nb > 0 {
			b := buf[0]
			gl.Logln("Read byte %x", b)
			// res = HandleByte(b, closer, mc.In) // this should return the final coerced byte or []byte that the window will be responsible for processing
			ka = byteHandler(b)
			if ka == nil {
				ka = sp.Get().(*KeyAction)
			}
			gl.Logln("Generated *KeyAction: %s", ka.String())
			mc.State.ActiveWindow.EventChan <- ka
			gl.Logln("*KeyAction Passed to active window")
			// mc.State.ActiveWindow.RawEventChan <- res
		}
	}
}

func (mc *MainConfig) CoerceInput(b byte) (inputSeq []byte) { // Will coerce input to an actionable sequence, will possibly read more bytes from the main input source, bypassing the main event loop
	// HandleModSequence(cell, modSeq)
	return []byte{}
}

func MakeKeyActionPool() *sync.Pool {
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
	sp = MakeKeyActionPool() // Create a pool of *SequenceNode references for dispatching normal ascii printable characters on the event channel for the window (trying to avoid as much re-allocation as possible
	return func(b byte) *KeyAction {
		res = res[:1]
		GlobalLogger.Logln("Byte Handler result buffer: %b", res)
		if b == 3 {
			ch <- struct{}{}
			//		} else if b == 13 {
			// os.Exit(0)
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

func (cl *ConcreteLogger) Log(message string, vars ...any) {
	cl.LogCh <- fmt.Sprintf(message, vars...)
}

func (cl *ConcreteLogger) Logln(message string, vars ...any) {
	cl.LogCh <- fmt.Sprintf(message+"\n", vars...)
}

func (cl *ConcreteLogger) Start() {
	for {
		select {
		case msg := <-cl.LogCh:
			fmt.Fprintf(cl.Out, msg)
		case <-cl.RunCh:
			return
		}
	}
}

func (cl *ConcreteLogger) Stop() {
	cl.RunCh <- struct{}{}
}

func (cl *ConcreteLogger) InitWithBuffer() {
	cl.LogCh = make(chan string, 1000)
	cl.RunCh = make(chan interface{})
	f, err := os.CreateTemp("./", ".term-reader-logger.txt.")
	if err != nil {
		fmt.Printf("Error opening tmp file: %s\n", err)
		cl.Out = io.Discard
	} else {
		os.Remove("term-reader-logger.txt")
		err := os.Symlink(f.Name(), "term-reader-logger.txt")
		if err != nil {
			fmt.Printf("Error creating logger symlink: %s\n", err)
		}
		cl.Out = f
		go cl.Start()
		cl.Log("Opened New logger at %s", f.Name())
	}
}

func (cl *ConcreteLogger) Init() {
	cl.LogCh = make(chan string)
	cl.RunCh = make(chan interface{})
	f, err := os.CreateTemp("./", ".term-reader-logger.txt.")
	if err != nil {
		fmt.Printf("Error opening tmp file: %s\n", err)
		cl.Out = io.Discard
	} else {
		os.Remove("term-reader-logger.txt")
		err := os.Symlink(f.Name(), "term-reader-logger.txt")
		if err != nil {
			fmt.Printf("Error creating logger symlink: %s\n", err)
		}
		cl.Out = f
		go cl.Start()
		cl.Log("Opened New logger at %s", f.Name())
	}
}

func NewConcreteLogger() (cl *ConcreteLogger) {
	cl = &ConcreteLogger{}
	cl.Init()
	return cl
}
