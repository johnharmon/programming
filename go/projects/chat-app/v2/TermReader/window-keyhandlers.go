package main

import (
	"bytes"
	"fmt"
	"os"
)

func NormalHandleForwardFind(w *Window, ka *KeyAction) bool {
	findBytes := ka.Value[0]
	var nextCursorCol int
	if w.CursorCol-1 < len(w.Buf.Lines[w.CursorLine]) {
		nextCursorCol = bytes.IndexByte(w.Buf.Lines[w.CursorLine][w.CursorCol-1:], findBytes)
	} else {
		nextCursorCol = -1
	}
	fmt.Fprintf(os.Stderr, "Next cursor Column := %d", nextCursorCol)

	if nextCursorCol != -1 {
		w.IncrCursorCol(nextCursorCol)
	}
	return false
}

func FindByteIndex(searchBuf []byte, b byte) (idx int) {
	return bytes.IndexByte(searchBuf, b)
}

func NormalHandleLeftMove(w *Window, count int) {
	w.IncrCursorCol(-count)
	w.MoveCursorToDisplayPosition()
}

func NormalHandleRightMove(w *Window, count int) {
	w.IncrCursorCol(count)
	w.MoveCursorToDisplayPosition()
}

func NormalHandleUpMove(w *Window, count int) {
	w.IncrCursorLine(-count)
	if w.CursorLine < w.BufTopLine && w.BufTopLine > 0 {
		w.BufTopLine -= count
		w.NeedRedraw = true
		// w.RedrawAllLines()
	}
	w.MoveCursorToDisplayPosition()
}

func NormalHandleDownMove(w *Window, count int) {
	oldLine := w.CursorLine
	w.IncrCursorLine(count)
	w.MoveCursorToDisplayPosition()
	if w.CursorLine > w.BufTopLine+w.Height && oldLine != w.CursorLine {
		w.BufTopLine += count
		w.NeedRedraw = true
		// w.RedrawAllLines()
	}
}

func InsertHandleDelete(w *Window) {
	w.Logger.Logln("Backspace Detected, content before deletion: %s", w.GetActiveLine())
	if len(w.GetActiveLine()) == 0 {
		w.Buf.Lines = DeleteLineAt(w.Buf.Lines, w.CursorLine, 1)
		w.IncrCursorLine(-1)
		numDisplayedLines := len(w.Buf.Lines) - w.BufTopLine
		if numDisplayedLines < w.Height && w.BufTopLine > 0 {
			w.BufTopLine--
		}
		w.NeedRedraw = true
		// w.RedrawAllLines()
		w.IncrCursorCol(len(w.GetActiveLine()))

	} else {
		w.Buf.Lines[w.Buf.ActiveLine] = DeleteByteAt(w.GetActiveLine(), w.CursorCol-1)
		w.Logger.Logln("Content After deletion: %s", w.GetActiveLine())
		w.RedrawLine(w.CursorLine)
		w.IncrCursorCol(-1)
	}
}

func InsertHandleEnter(w *Window) {
	newLine := MakeNewLines(1, 256)
	w.Logger.Logln("Enter detected")
	w.Logger.Logln("Inserting new line at index %d", w.CursorLine+1)
	w.Buf.Lines = InsertLineAt(w.Buf.Lines, newLine, w.CursorLine+1)
	w.IncrCursorLine(1)
	w.CursorCol = 1
	if w.CursorLine-w.BufTopLine >= w.Height {
		w.BufTopLine++
	}
	w.NeedRedraw = true
	// w.RedrawAllLines()
	w.MoveCursorToDisplayPosition()
}

func NormalHandleEnter(w *Window) {
	newLine := MakeNewLines(1, 256)
	w.Logger.Logln("Enter detected")
	w.Logger.Logln("Inserting new line at index %d", w.CursorLine+1)
	w.Buf.Lines = InsertLineAt(w.Buf.Lines, newLine, w.CursorLine+1)
	w.IncrCursorLine(1)
	w.CursorCol = 1
	if w.CursorLine-w.BufTopLine > w.Height {
		w.BufTopLine++
	}
	w.Mode = MODE_INSERT
	w.NeedRedraw = true
	// w.RedrawAllLines()
	w.MoveCursorToDisplayPosition()
}

func InsertHandleArrowRight(w *Window) {
	w.IncrCursorCol(1)
	w.MoveCursorToDisplayPosition()
}

func InsertHandleArrowLeft(w *Window) {
	w.IncrCursorCol(-1)
	w.MoveCursorToDisplayPosition()
}

func InsertHandleArrowUp(w *Window) {
	w.IncrCursorLine(-1)
	if w.CursorLine < w.BufTopLine && w.BufTopLine > 0 {
		w.BufTopLine--
		w.NeedRedraw = true
		// w.RedrawAllLines()
	}
	w.MoveCursorToDisplayPosition()
}

func InsertHandleArrowDown(w *Window) {
	oldLine := w.CursorLine
	w.IncrCursorLine(1)
	if w.CursorLine > w.BufTopLine+w.Height-1 && oldLine != w.CursorLine {
		w.BufTopLine++
		w.NeedRedraw = true
		// w.RedrawAllLines()
	}
	w.MoveCursorToDisplayPosition()
}

func NormalHandleArrowRight(w *Window) {
	w.IncrCursorCol(1)
	w.MoveCursorToDisplayPosition()
}

func NormalHandleArrowLeft(w *Window) {
	w.IncrCursorCol(-1)
	w.MoveCursorToDisplayPosition()
}

func NormalHandleArrowUp(w *Window) {
	w.IncrCursorLine(-1)
	if w.CursorLine < w.BufTopLine && w.BufTopLine > 0 {
		w.BufTopLine--
		w.NeedRedraw = true
		// w.RedrawAllLines()
	}
	w.MoveCursorToDisplayPosition()
}

func NormalHandleArrowDown(w *Window) {
	oldLine := w.CursorLine
	w.IncrCursorLine(1)
	if w.CursorLine > w.BufTopLine+w.Height && oldLine != w.CursorLine {
		w.BufTopLine++
		w.NeedRedraw = true
		// w.RedrawAllLines()
	}
	w.MoveCursorToDisplayPosition()
}

func CmdHandleDelete(w *Window) {
	w.Logger.Logln("Backspace Detected, content before deletion: %s", w.CmdBuf)
	if len(w.CmdBuf) != 0 {
		w.CmdBuf = DeleteByteAt(w.CmdBuf, w.CmdCursorCol-1)
		w.Logger.Logln("Content After deletion: %s", w.GetActiveLine())
		w.DisplayCmdLine()
		w.IncrCmdCursorCol(-1)
	}
}

func CmdHandleArrowRight(w *Window) {
	w.IncrCmdCursorCol(1)
	w.MoveCursorToCmdPosition()
}

func CmdHandleArrowLeft(w *Window) {
	w.IncrCmdCursorCol(-1)
	w.MoveCursorToCmdPosition()
}

func (w *Window) AgnosticHandleDelete(line []byte, index int) (newLine []byte) {
	w.Logger.Logln("Backspace Detected, content before deletion: %s", w.CmdBuf)
	if len(line) != 0 {
		newLine = DeleteByteAt(line, index)
		w.Logger.Logln("Content After deletion: %s", w.GetActiveLine())
	}
	return newLine
}
