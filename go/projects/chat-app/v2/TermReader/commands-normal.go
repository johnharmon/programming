package main

import (
	"bytes"
	"fmt"
	"os"
)

func NormalHandleUpMoveCmd(w *Window, ac *ActionContext) {
	GlobalLogger.Logln("NormalHandleUpMoveCmd invoked")
	count := ac.Count
	w.IncrCursorLine(-count)
	if w.CursorLine < w.BufTopLine && w.BufTopLine > 0 {
		w.BufTopLine -= count
		w.NeedRedraw = true
	}
	w.MoveCursorToDisplayPosition()
}

func NormalHandleLeftMoveCmd(w *Window, ac *ActionContext) {
	GlobalLogger.Logln("NormalHandleLeftMoveCmd invoked, count: %d", ac.Count)
	count := ac.Count
	IncrTwoCursorColPtr(w.Buf.Lines[w.CursorLine], &w.CursorCol, &w.DesiredCursorCol, -count)
	w.MoveCursorToDisplayPosition()
}

func NormalHandleRightMoveCmd(w *Window, ac *ActionContext) {
	GlobalLogger.Logln("NormalHandleRightMoveCmd invoked")
	count := ac.Count
	IncrTwoCursorColPtr(w.Buf.Lines[w.CursorLine], &w.CursorCol, &w.DesiredCursorCol, count)
	w.MoveCursorToDisplayPosition()
}

func NormalHandleDownMoveCmd(w *Window, ac *ActionContext) {
	GlobalLogger.Logln("NormalHandleDownMoveCmd invoked")
	count := ac.Count
	oldLine := w.CursorLine
	w.IncrCursorLine(count)
	w.MoveCursorToDisplayPosition()
	if w.CursorLine > w.BufTopLine+w.Height && oldLine != w.CursorLine {
		w.BufTopLine += count
		w.NeedRedraw = true
	}
}

func NormalModeSwitchToInsert(w *Window, ac *ActionContext) {
	w.Mode = MODE_INSERT
}

func NormalHandleForwardFind(w *Window, ac *ActionContext) {
	findBytes := ac.Suffix[0]
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
}
