package main

func NormalHandleUpMoveCmd(w *Window, ac *ActionContext) {
	count := ac.Count
	w.IncrCursorLine(-count)
	if w.CursorLine < w.BufTopLine && w.BufTopLine > 0 {
		w.BufTopLine -= count
		w.NeedRedraw = true
	}
	w.MoveCursorToDisplayPosition()
}

func NormalHandleLeftMoveCmd(w *Window, ac *ActionContext) {
	count := ac.Count
	IncrCursorColPtr(w.Buf.Lines[w.CursorLine], &w.CursorCol, -count)
	w.MoveCursorToDisplayPosition()
}

func NormalHandleRightMoveCmd(w *Window, ac *ActionContext) {
	count := ac.Count
	IncrCursorColPtr(w.Buf.Lines[w.CursorLine], &w.CursorCol, count)
	w.MoveCursorToDisplayPosition()
}

func NormalHandleDownMoveCmd(w *Window, ac *ActionContext) {
	count := ac.Count
	oldLine := w.CursorLine
	w.IncrCursorLine(count)
	w.MoveCursorToDisplayPosition()
	if w.CursorLine > w.BufTopLine+w.Height && oldLine != w.CursorLine {
		w.BufTopLine += count
		w.NeedRedraw = true
	}
}
