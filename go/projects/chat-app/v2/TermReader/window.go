package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"

	"golang.org/x/term"
)

func CreateCleanupKeyGenerator(cleanupTaskMap map[string]*CleanupTask) func(string) string {
	return func(name string) string {
		iterator := 1
		iterName := name
		for {
			_, ok := cleanupTaskMap[iterName]
			if !ok {
				return iterName
			} else {
				iterName = fmt.Sprintf("%s-%d", name, iterator)
				iterator++
			}
		}
	}
}

func CreateCleanupKeyInserter(cleanupTaskMap map[string]*CleanupTask) func(string, *CleanupTask) {
	return func(key string, ct *CleanupTask) {
		cleanupTaskMap[key] = ct
	}
}

func CreateCleanupTaskStarter(cleanupTaskMap map[string]*CleanupTask) func() {
	return func() {
		for name, task := range cleanupTaskMap {
			if task.Start {
				GlobalLogger.Logln("Starting cleanup task %s", name)
				task.Func()
			}
		}
	}
}

func CreateCleanupTaskRegistrar(cleanupTaskMap map[string]*CleanupTask) func(chan *sync.WaitGroup, func(), string, bool) {
	return func(closer chan *sync.WaitGroup, Func func(), name string, start bool) {
		ct := &CleanupTask{Closer: closer, Name: name, Func: Func, Start: start}
		cleanupKey := GenCleanupKey(ct.Name)
		InsertCleanupKey(cleanupKey, ct)
		if ct.Start {
			go ct.Func()
		}
	}
}

// func RegisterCleanupTask(name string, closer chan *sync.WaitGroup, Func func(), start bool) {
//	ct := &CleanupTask{Closer: closer, Name: name, Func: Func, Start: start}
//	cleanupKey := GenCleanupKey(ct.Name)
//	InsertCleanupKey(cleanupKey, ct)

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
	w.CmdBuf = make([]byte, 1, 256)
	w.Width = width
	w.Mode = MODE_INSERT
	w.BufTopLine = 0
	w.CmdCursorCol = 2
	w.CursorCol = 1
	w.DesiredCursorCol = 1
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
		ka := <-w.KeyActionReturner
		if ka.FromPool {
			sp.Put(ka)
		}
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
	w.Buf.Lines[w.Buf.ActiveLine] = InsertAt(w.Buf.Lines[w.Buf.ActiveLine], b, w.CursorCol-1)
}

func WriteToLine(line []byte, b []byte, start int) (newLine []byte) {
	// w.Logger.Logln("Writing %b to %d", b, w.Buf.ActiveLine)
	return InsertAt(line, b, start)
}

func (w *Window) WriteToCmd(b []byte) {
	w.CmdBuf = InsertAt(w.CmdBuf, b, w.CursorCol-1)
}

func IncrCursorCol(line []byte, col int, incr int) (newCol int) {
	lLen := len(line)
	newPos := col + incr
	if newPos < 1 {
		newPos = 1
		newCol = newPos
	} else if newPos <= lLen+1 {
		newCol = newPos
	} else if newPos > lLen+1 {
		newPos = lLen + 1
	} else {
		newCol = lLen
	}
	return newCol
}

func IncrCursorColPtr(line []byte, col *int, incr int) {
	lLen := len(line)
	newPos := *col + incr
	GlobalLogger.Logln("New Cursor Col Target: %d", newPos)
	switch {
	case newPos < 1:
		newPos = 1
		*col = newPos
	case newPos <= lLen+1:
		*col = newPos
	case newPos > lLen+1:
		*col = lLen + 1
	default:
		*col = lLen
	}
	GlobalLogger.Logln("New Cursor Col: %d", *col)
}

func IncrTwoCursorColPtr(line []byte, col1 *int, col2 *int, incr int) {
	lLen := len(line)
	newPos := *col1 + incr
	GlobalLogger.Logln("New Cursor col1 Target: %d", newPos)
	switch {
	case newPos < 1:
		newPos = 1
		*col1 = newPos
		*col2 = newPos
	case newPos <= lLen+1:
		*col1 = newPos
		*col2 = newPos
	case newPos > lLen+1:
		*col1 = lLen + 1
		*col2 = lLen + 1
	default:
		*col1 = lLen
	}
	GlobalLogger.Logln("New Cursor Col: %d,%d", *col1, *col2)
}

func (w *Window) IncrCmdCursorCol(incr int) {
	lLen := len(w.CmdBuf)
	newPos := w.CursorCol + incr
	if newPos < 1 {
		newPos = 1
		w.CursorCol = newPos
		w.DesiredCursorCol = newPos
	} else if newPos <= lLen+1 {
		w.CursorCol = newPos
		w.DesiredCursorCol = newPos
	} else if newPos > lLen+1 {
		newPos = lLen + 1
	} else {
		w.CursorCol = lLen
	}
}

func (w *Window) IncrCmdCursorCol2(incr int) {
	lLen := len(w.CmdBuf)
	newPos := w.CmdCursorCol + incr
	if newPos < 1 {
		newPos = 1
		w.CmdCursorCol = newPos
		// w.DesiredCursorCol = newPos
	} else if newPos <= lLen+1 {
		w.CmdCursorCol = newPos
		// w.DesiredCursorCol = newPos
	} else if newPos > lLen+1 {
		newPos = lLen + 1
	} else {
		w.CmdCursorCol = lLen
	}
}

func (w *Window) IncrCursorCol(incr int) {
	lLen := len(w.Buf.Lines[w.Buf.ActiveLine])
	newPos := w.CursorCol + incr
	if newPos < 1 {
		newPos = 1
		w.CursorCol = newPos
		w.DesiredCursorCol = newPos
	} else if newPos <= lLen+1 {
		w.CursorCol = newPos
		w.DesiredCursorCol = newPos
	} else if newPos > lLen+1 {
		newPos = lLen + 1
	} else {
		w.CursorCol = lLen
	}
}

func (w *Window) GetDisplayCursorCol() int {
	if len(w.Buf.Lines[w.Buf.ActiveLine])+1 < w.DesiredCursorCol {
		w.Logger.Logln("Setting cursor display position to: %d", len(w.Buf.Lines[w.Buf.ActiveLine]))
		return len(w.Buf.Lines[w.Buf.ActiveLine])
	} else {
		w.Logger.Logln("Desired Cursor position is compatible with the active line")
		return w.DesiredCursorCol
	}
}

func (w *Window) HandleArrowDown() (col int) {
	w.IncrCursorLine(1)
	col = w.GetDisplayCursorCol()
	return col
}

func (w *Window) IncrCursorLine(vec int) {
	nextLine := w.CursorLine + vec
	if nextLine >= 0 && nextLine < len(w.Buf.Lines) {
		w.Buf.ActiveLine = nextLine
		w.CursorLine = nextLine
	}
}

func MakeNewLines(lines int, lineSize int) [][]byte {
	newLines := make([][]byte, lines)
	for i := range newLines {
		newLines[i] = make([]byte, 0, lineSize)
	}
	return newLines
}

func (w *Window) RedrawLine(ln int) {
	screenLine := ln - w.BufTopLine
	if screenLine >= 0 && screenLine < w.BufTopLine+w.Height {
		w.MoveCursorToPosition(screenLine+w.TermTopLine, 1)
		w.Logger.Logln("Writing content to line #%d:", screenLine)
		w.Logger.Logln("%s", w.Buf.Lines[ln])
		w.Out.Write(TERM_CLEAR_LINE)
		w.Out.Write(w.Buf.Lines[ln])
	}
}

func (w *Window) GetActiveLine() []byte {
	if w.CursorLine >= 0 && w.CursorLine < len(w.Buf.Lines) {
		return w.Buf.Lines[w.CursorLine]
	} else {
		return nil
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

func (w *Window) MoveCursorToCmdPosition() {
	cursorDisplayLine := w.TermTopLine + w.Height
	w.MoveCursorToPosition(cursorDisplayLine, w.CursorCol)
}

func (w *Window) ProcessCmd() (err error) {
	GlobalLogger.Logln("Processing cmd: %s", w.CmdBuf)
	cmd, cmdArgs := ProcessCmdArgs(w.CmdBuf)
	GlobalLogger.Logln("cmd: %s, cmdArgs: %s", cmd, strings.Join(cmdArgs, "|"))
	fmt.Fprintf(os.Stderr, "cmd: %s, cmdArgs: %s\n", cmd, strings.Join(cmdArgs, "|"))
	w.CmdBuf = w.CmdBuf[:1]
	w.CmdCursorCol = 2
	cmdDispatch, ok := COMMANDS[cmd]
	if !ok {
		err = errors.New(fmt.Sprintf("Error: cmd not found: \"%s\"", cmd))
		w.DisplayCmdMessage(err.Error())
		return err
	}
	err = cmdDispatch.ExecFunc(w, cmdArgs...)
	if err != nil {
		w.DisplayCmdMessage(fmt.Sprintf("Error: %s", err))
	}
	return err
}

func (w *Window) MoveCursorToDisplayPosition() {
	cursorDisplayLine := (w.CursorLine - w.BufTopLine) + 1
	if cursorDisplayLine >= w.TermTopLine+w.Height {
		cursorDisplayLine = w.TermTopLine + w.Height - 1
	}
	cursorCol := w.GetDisplayCursorCol()
	w.MoveCursorToPosition(cursorDisplayLine, cursorCol)
}

func (w *Window) Redraw(handler func() []int) {
	w.Logger.Logln("Redrawing window")
	linesToRedraw := handler()
	w.Logger.Logln("Lines to redraw: %s", linesToRedraw)
	lastIndex := 0
	w.MoveCursorToPosition(w.TermTopLine, 1)
	for _, lineNum := range linesToRedraw {
		w.MoveCursorToPosition(w.TermTopLine+lineNum, 1)
		w.RedrawLine(lineNum + w.BufTopLine)
		lastIndex++
	}
	// w.DisplayStatusLine()
}

func (w *Window) RedrawAllLines() {
	w.Logger.Logln("Forcing full redraw of all lines")
	w.MoveCursorToPosition(w.TermTopLine, 1)
	var lineLimit int
	linesLeftInBuffer := len(w.Buf.Lines) - w.BufTopLine - 1
	if linesLeftInBuffer < w.Height-1 {
		lineLimit = linesLeftInBuffer
	} else {
		lineLimit = w.Height - 1
	}
	w.Logger.Logln("Line limit calculated: %d", lineLimit)
	for i := w.BufTopLine; i <= w.BufTopLine+lineLimit; i++ {
		w.RedrawLine(i)
	}
}

func (w *Window) DisplayCmdLine() {
	termStatusLineNum := w.TermTopLine + w.Height
	w.MoveCursorToPosition(termStatusLineNum, 0)
	w.Out.Write(TERM_CLEAR_LINE)
	padding := strings.Repeat(" ", w.Width-len(w.CmdBuf))
	fmt.Fprintf(w.Out, "\x1b[48;5;202m\x1b[38;5;16m%s%s\x1b[00m", w.CmdBuf, padding)
}

func (w *Window) DisplayCmdMessage(msg string) {
	termStatusLineNum := w.TermTopLine + w.Height + 2
	w.MoveCursorToPosition(termStatusLineNum, 0)
	w.Out.Write(TERM_CLEAR_LINE)
	fmt.Fprintf(w.Out, "\x1b[48;5;202m\x1b[38;5;16m%s\x1b[00m", msg)
}

func (w *Window) DisplayStatusLine() {
	termStatusLineNum := w.TermTopLine + w.Height
	w.MoveCursorToPosition(termStatusLineNum, 0)
	w.Out.Write(TERM_CLEAR_LINE)
	statusLine := fmt.Sprintf(
		"%sCursorLine: %d | CursorColumn: %d | BufTopLine: %d | DesiredCursorColumn: %d | TermLine: %d | LineLength: %d | BufferLength: %d | WindowHeight: %d",
		"\x1b[41;30m",
		w.CursorLine,
		w.CursorCol,
		w.BufTopLine,
		w.DesiredCursorCol,
		w.TermTopLine+(w.CursorLine-w.BufTopLine),
		len(w.Buf.Lines[w.Buf.ActiveLine]),
		len(w.Buf.Lines),
		w.Height)
	padding := strings.Repeat(" ", w.Width-len(statusLine))
	fmt.Fprintf(w.Out, "%s%s", statusLine, padding)
	fmt.Fprintf(w.Out, "\r\x1b[00m")
	fmt.Fprintf(w.Out, "\nMode: %d", w.Mode)
	fmt.Fprintf(w.Out, "\r\x1b[00m")
	// fmt.Fprintf(w.Out, "FlushToken: %s", GlobalLogger.(*ConcreteLogger).FlushBuffer.Bytes())
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

func (cl *ConcreteLogger) Cleanup() {
	wg := <-cl.RunCh
	defer wg.Done()
	cl.Logln("Closing the logging channel")
	close(cl.Done)
	// close(cl.LogOutput)
	cl.Mu.Lock()
}

func Cleanup(closer chan struct{}, fd int, oldState *term.State, taskMap map[string]*CleanupTask) {
	wg1 := &sync.WaitGroup{}
	wg2 := &sync.WaitGroup{}
	<-closer
	GlobalLogger.Logln("Caught Main Cleanup termination signal")
	for key, task := range taskMap {
		if key != LOGGER_CLEANUP_UNIQUE_KEY {
			GlobalLogger.Logln("Sending termination signal to task: %s", key)
			wg1.Add(1)
			task.Closer <- wg1
		}
	}
	GlobalLogger.Logln("Waiting for cleanup tasks to finsh...")
	wg1.Wait()
	if gl, ok := GlobalLogger.(*ConcreteLogger); ok {
		CleanLogFiles(gl.LogFileName)
	}
	GlobalLogger.Logln("Cleanup tasks finished, shutting down logger")
	if task, ok := taskMap[LOGGER_CLEANUP_UNIQUE_KEY]; ok {
		wg2.Add(1)
		GlobalLogger.Logln("wg2 added")
		task.Closer <- wg2
		wg2.Wait()
	}
	fmt.Println("\n\rRestoring old state")
	term.Restore(fd, oldState)
	os.Exit(0)
}

func CleanLogFiles(currentFile string) {
	// assume files are in cwd for now
	entries, err := os.ReadDir("./")
	entryRegex, _ := regexp.Compile(`\.term-reader-logger\.json.*`)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Could not open ./")
	}
	// GlobalLogger.Logln("Current log file: %s", currentFile)
	fmt.Fprintf(os.Stderr, "Current log file: %s\n", currentFile)
	for _, entry := range entries {
		// GlobalLogger.Logln("Processing log file: %s", entry.Name())
		fmt.Fprintf(os.Stderr, "Pocessing log file; %s\n", "./"+entry.Name())
		if "./"+entry.Name() != currentFile && entryRegex.Match([]byte(entry.Name())) {
			err = os.Remove(entry.Name())
			if err != nil {
				fmt.Fprintln(os.Stderr, "Could not open ./")
			}
		}
	}
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
	gl := GlobalLogger.(*ConcreteLogger)
	RegisterCleanupTask(gl.RunCh, gl.Cleanup, LOGGER_CLEANUP_UNIQUE_KEY, true)
	var ka *KeyAction
	fd := int(os.Stdin.Fd())
	gl.Logln("Setting terminal to raw mode")
	oldState, err := term.MakeRaw(fd)
	if err != nil {
		panic(err)
	}
	gl.Logln("Making closer channel")
	closer := make(chan struct{})
	gl.Logln("Making *KeyAction pool")
	sp := MakeKeyActionPool() // Create a pool of *SequenceNode references for dispatching normal ascii printable characters on the event channel for the window (trying to avoid as much re-allocation as possible
	go mc.State.ActiveWindow.RunKeyActionReturner(sp)
	gl.Logln("Making *KeyAction return channel for the pool")
	keyActionReturner := make(chan *KeyAction, 1000)
	gl.Logln("Spinning off goroutine for returning *KeyActions to the pool")
	go ReturnKeyActionsToPool(sp, keyActionReturner)
	gl.Logln("Spinning off cleanup goroutine")
	go Cleanup(closer, fd, oldState, CleanupTaskMap)
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

func MakeByteHandler(ch chan struct{}, in io.Reader, kaPool *sync.Pool) func(byte) *KeyAction { // returns a byte handling function that will reuse an input buffer so re-allocation does not happen on every byte handled by the main loop
	res := make([]byte, 1, 8)
	var seqN *KeyAction
	kaPool = MakeKeyActionPool() // Create a pool of *SequenceNode references for dispatching normal ascii printable characters on the event channel for the window (trying to avoid as much re-allocation as possible
	return func(b byte) *KeyAction {
		GlobalLogger.Logln("Byte Handler: len(res) = %d, cap(res) = %d", len(res), cap(res))
		res = res[:1]
		res[0] = b
		// GlobalLogger.Logln("Byte Handler: len(res) = %d, cap(res) = %d", len(res), cap(res))
		// GlobalLogger.Logln("Byte passed to handler: %b", b)
		// GlobalLogger.Logln("Byte Handler result buffer: %b", res)
		if b == 3 {
			GlobalLogger.Logln("Interrupt caught, sending shutdown signal to global cleanup")
			ch <- struct{}{}
			//		} else if b == 13 {
			// os.Exit(0)
			select {}
			// return nil
		} else if b >= 0x20 && b <= 0x7E {
			GlobalLogger.Logln("Getting *KeyAction from pool")
			seqN = kaPool.Get().(*KeyAction)
			seqN.Value[0] = b
			// seqN.Action = "Raw"
			seqN.Value = seqN.Value[0:1]
			return seqN
		}
		// GlobalLogger.Logln("Result before byte Parsing: %b", res)
		res = ParseByte(b, res, in)
		// GlobalLogger.Logln("Result after byte Parsing: %b", res)
		// time.Sleep(time.Millisecond * 100)
		defer clear(res)
		return CoerceInputToAction(res)
	}
}

func ParseByte(b byte, result []byte, in io.Reader) []byte { // Should handle initial detection for multi-byte sequences, if a single byte sequence then just return the byte as a slice
	result[0] = b
	// GlobalLogger.Logln("Result before byte parsing: %b", result)
	if b == 0x1b {
		GlobalLogger.Logln("Escape sequence detectet, entering read timeout")
		n, _ := ReadMultiByteSequence(result, in, time.Millisecond*25)
		// GlobalLogger.Logln("Number of bytes read into result: %d", n)
		result = result[0 : 1+n]
	}
	// GlobalLogger.Logln("Result after byte parsing: %b", result)
	return result
}

func ReadMultiByteSequence(buf []byte, input io.Reader, timeout time.Duration) (n int, err error) { // will read a multi-byte sequence into buf, respecting any existing elements
	bufLen := len(buf)
	bufCap := cap(buf)
	numB := new(int)
	// readBuf := make([]byte, 0, bufCap-bufLen)
	// GlobalLogger.Logln("Result inside byte paring: %b", buf)
	go func() {
		n, err = input.Read(buf[bufLen:bufCap])
		numB = &n
	}()
	<-time.After(timeout)
	GlobalLogger.Logln(fmt.Sprintf("Result right after byte read: %b | bufLen: %d | bufCap: %d", buf[0:bufCap], len(buf), cap(buf)))
	return (*numB), nil
	// deadline := time.Now().Add(timeout)
	// if f, ok := input.(*os.File); ok {
	// f.SetReadDeadline(deadline)
	// defer ClearReadDeadline(f)
	// GlobalLogger.Logln("Result right before byte read: %b | bufLen: %d | bufCap: %d", buf[0:bufCap], bufLen, bufCap)
	// n, err = f.Read(buf[bufLen:bufCap])
	// GlobalLogger.Logln("Number of bytes read: %d", n)
	//		if err != nil {
	//			return 0, err
	//		}
	//		return n, nil
	//	}
	//	return 0, nil
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

//func (cl *ConcreteLogger) Log(message string, vars ...any) {
//	cl.LogCh <- []byte(fmt.Sprintf(message, vars...))
//}

//	func (cl *ConcreteLogger) Logln(message string, vars ...any) {
//		go func() {
//			timestamp := time.Now().Format(time.StampMicro)
//			entry := cl.LogEntryPool.Get().(*LogEntry)
//			entry.Message = fmt.Sprintf(message, vars...)
//			entry.Timestamp = timestamp
//			msg, err := json.Marshal(entry)
//			if err != nil {
//				os.Exit(1)
//			}
//			select {
//			case cl.LogCh <- string(msg) + "\n":
//			case <-cl.Done:
//				fmt.Fprint(cl.Out, message)
//			}
//		}()
//	}

//func (cl *ConcreteLogger) Logln(message string, vars ...any) {
//	// fmt.Fprintf(os.Stderr, "Logln called\n")
//	rawLogArgs := cl.RawLogArgPool.Get().(*RawLogArgs)
//	rawLogArgs.FormatMessage = message
//	rawLogArgs.FormatArgs = vars
//	rawLogArgs.Timestamp = time.Now().Format(time.StampMicro)
//	cl.RawLogCh <- rawLogArgs
//}
//
//func (cl *ConcreteLogger) RawLogHandler() {
//	for logArgs := range cl.RawLogCh {
//		// fmt.Fprintf(os.Stderr, "RawLogHandler received log: %s\n", logArgs.FormatMessage)
//		entry := cl.LogEntryPool.Get().(*LogEntry)
//		entry.Message = fmt.Sprintf(logArgs.FormatMessage, logArgs.FormatArgs...)
//		entry.Timestamp = logArgs.Timestamp
//		cl.LogEntryCh <- entry
//		cl.RawLogArgPool.Put(logArgs)
//	}
//	close(cl.LogEntryCh)
//}
//
//func (cl *ConcreteLogger) JsonMarshaler() {
//	encodeBuffer := bytes.NewBuffer(make([]byte, 0, 2048))
//	encoder := json.NewEncoder(encodeBuffer)
//	for rawLog := range cl.LogEntryCh {
//		encodeSlice := cl.MessageBufferPool.New().([]byte)
//		encodeSlice = encodeSlice[:0]
//		// fmt.Fprintf(os.Stderr, "JsonMarshaler received Log entry: +%v\n", rawLog)
//		encodeBuffer.Reset()
//		encoder.Encode(rawLog)
//		encodeSlice = append(encodeSlice, encodeBuffer.Bytes()...)
//		cl.LogOutput <- encodeSlice
//		os.Stderr.Write(encodeBuffer.Bytes())
//		cl.LogEntryPool.Put(rawLog)
//	}
//	close(cl.LogOutput)
//}
//
//func (cl *ConcreteLogger) JsonWriter() {
//	// cl.Mu.Lock()
//	var flushToken *FlushToken
//	bufFlushSize := 1024
//	ticker := time.NewTicker(time.Millisecond * 100)
//LogWriteLoop:
//	for {
//		select {
//		case msg, ok := <-cl.LogOutput:
//			// fmt.Fprintf(os.Stderr, "JsonWriter received Log: %s\n", msg)
//			if !ok {
//				if cl.ActiveBuffer.Len() > 0 {
//					flushToken = <-cl.FlushReceiver
//					flushToken.HandledBy = "JsonWriter(): case msg, ok := <- cl.Logch; if !ok {<this>}"
//					flushToken.SentBy = "JsonWriter(): case msg, ok := <- cl.Logch; if !ok {<this>}"
//					cl.Out.Write(cl.ActiveBuffer.Bytes())
//				}
//				cl.Mu.Unlock()
//				return
//			}
//			cl.SwapMu.Lock()
//			cl.ActiveBuffer.Write(msg)
//			os.Stderr.Write(msg)
//			cl.SwapMu.Unlock()
//			cl.MessageBufferPool.Put(msg)
//			if cl.ActiveBuffer.Len() >= bufFlushSize {
//				select {
//				case flushToken = <-cl.FlushReceiver:
//					flushToken.HandledBy = "JsonWriter(): case msg, ok := <- cl.Logch"
//					flushToken.SentBy = "JsonWriter(): case msg, ok := <- cl.Logch"
//					cl.FlushAndSwapActiveBuffer()
//					cl.FlushSender <- flushToken
//				default:
//					continue LogWriteLoop
//				}
//			}
//		case <-ticker.C:
//			if cl.ActiveBuffer.Len() > 0 {
//				select {
//				case flushToken = <-cl.FlushReceiver:
//					flushToken.HandledBy = "JsonWriter(): case <- ticker.C:"
//					flushToken.SentBy = "JsonWriter(): case <- ticker.C:"
//					cl.FlushAndSwapActiveBuffer()
//					cl.FlushSender <- flushToken
//				default:
//					continue LogWriteLoop
//				}
//			}
//		}
//	}
//}
//
//func (cl *ConcreteLogger) FlushAndSwapActiveBuffer() {
//	go func() {
//		cl.SwapMu.Lock()
//		cl.ActiveBuffer, cl.FlushBuffer = cl.FlushBuffer, cl.ActiveBuffer
//		cl.SwapMu.Unlock()
//		flushToken := <-cl.FlushSender
//		flushToken.Iteration++
//		flushToken.HandledBy = "FlushAndSwapActiveBuffer()"
//		flushToken.ReceivedBy = "FlushAndSwapActiveBuffer()"
//		cl.Out.Write(cl.FlushBuffer.Bytes())
//		cl.FlushBuffer.Reset()
//		cl.FlushReceiver <- flushToken
//	}()
//}
//
//func (cl *ConcreteLogger) StartAsync() {
//	cl.Mu.Lock()
//	flushToken := &FlushToken{Iteration: 0}
//	cl.FlushReceiver <- flushToken
//	cl.RawLogHandler()
//	cl.JsonMarshaler()
//	cl.JsonWriter()
//	for {
//		_, ok := <-cl.Done
//		if !ok {
//			// fmt.Fprintf(os.Stderr, "Logging close received, exiting")
//			close(cl.RawLogCh)
//		}
//	}
//}
//
//func (cl *ConcreteLogger) Start() {
//	cl.Mu.Lock()
//	flushToken := &FlushToken{Iteration: 0}
//	go cl.RawLogHandler()
//	go cl.JsonMarshaler()
//	go cl.JsonWriter()
//	cl.FlushReceiver <- flushToken
//	for {
//		// fmt.Fprintf(os.Stderr, "Waiting on cl.Done to be closed\n")
//		_, ok := <-cl.Done
//		if !ok {
//			// fmt.Fprintf(os.Stderr, "cl.Done closed, closing cl.RawLogCh\n")
//			cl.RawLogCh <- &RawLogArgs{FormatMessage: "FlushToken{Iterations: %d HandledBy: %s, SentBy: %s", FormatArgs: []any{flushToken.Iteration, flushToken.HandledBy, flushToken.SentBy}}
//			close(cl.RawLogCh)
//			break
//		}
//	}
//}

//func (cl *ConcreteLogger) Start() {
//	cl.Mu.Lock()
//	bufFlushSize := 1024
//	flushBuffer := &bytes.Buffer{}
//	activeBuffer := &bytes.Buffer{}
//	ticker := time.NewTicker(time.Millisecond * 100)
//	cl.FlushCh <- struct{}{}
//listenLoop:
//	for {
//		select {
//		case msg := <-cl.LogCh:
//			activeBuffer.Write([]byte(msg))
//			if activeBuffer.Len() >= bufFlushSize {
//				cl.Out.Write(activeBuffer.Bytes())
//				activeBuffer.Reset()
//				activeBuffer, flushBuffer = flushBuffer, activeBuffer
//			}
//		case <-ticker.C:
//			if activeBuffer.Len() > 0 {
//				cl.Out.Write(activeBuffer.Bytes())
//				activeBuffer.Reset()
//				activeBuffer, flushBuffer = flushBuffer, activeBuffer
//			}
//		case <-cl.Done:
//			for msg := range cl.LogCh {
//				activeBuffer.Write([]byte(msg))
//			}
//			cl.Out.Write(activeBuffer.Bytes())
//			break listenLoop
//		}
//	}
//	cl.Mu.Unlock()
//}

func (cl *ConcreteLogger) Stop() {
	cl.RunCh <- &sync.WaitGroup{}
}

//func (cl *ConcreteLogger) InitWithBuffer() {
//	cl.LogCh = make(chan []byte, 1000)
//	cl.RunCh = make(chan *sync.WaitGroup)
//	f, err := os.CreateTemp("./", ".term-reader-logger.json.")
//	if err != nil {
//		fmt.Printf("Error opening tmp file: %s\n", err)
//		cl.Out = io.Discard
//	} else {
//		os.Remove("term-reader-logger.json")
//		err := os.Symlink(f.Name(), "term-reader-logger.json")
//		if err != nil {
//			fmt.Printf("Error creating logger symlink: %s\n", err)
//		}
//		cl.Out = f
//		go cl.Start()
//		cl.Log("Opened New logger at %s", f.Name())
//	}
//}

func (cl *ConcreteLogger) Init() {
	// cl.LogCh = make(chan string)
	cl.Mu = &sync.Mutex{}
	cl.FlushMu = &sync.Mutex{}
	cl.SwapMu = &sync.Mutex{}
	cl.ActiveBuffer = &bytes.Buffer{}
	cl.FlushBuffer = &bytes.Buffer{}
	cl.FlushReceiver = make(chan *FlushToken)
	cl.FlushSender = make(chan *FlushToken)
	cl.LogOutput = make(chan []byte, 1000)
	cl.RawLogCh = make(chan *RawLogArgs, 1000)
	cl.LogEntryCh = make(chan *LogEntry, 1000)
	cl.RunCh = make(chan *sync.WaitGroup)
	cl.Done = make(chan struct{})
	cl.LogEntryPool = &sync.Pool{}
	cl.LogEntryPool.New = func() any {
		return &LogEntry{}
	}
	cl.RawLogArgPool = &sync.Pool{}
	cl.RawLogArgPool.New = func() any {
		return &RawLogArgs{}
	}
	cl.MessageBufferPool = &sync.Pool{}
	cl.MessageBufferPool.New = func() any {
		return make([]byte, 2048)
	}

	f, err := os.CreateTemp("./", ".term-reader-logger.json.")
	if err != nil {
		fmt.Printf("Error opening tmp file: %s\n", err)
		cl.Out = os.Stderr
	} else {
		os.Remove("term-reader-logger.json")
		cl.LogFileName = f.Name()
		err := os.Symlink(f.Name(), "term-reader-logger.json")
		if err != nil {
			fmt.Printf("Error creating logger symlink: %s\n", err)
		}
		cl.Out = f
		go cl.Start()
		cl.Logln("Opened New logger at %s", f.Name())
	}
}

func NewConcreteLogger() (cl *ConcreteLogger) {
	cl = &ConcreteLogger{}
	cl.Init()
	return cl
}
