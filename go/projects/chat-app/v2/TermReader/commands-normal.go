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

func NormalHandleDeleteCharNoCursorMove(w *Window, ac *ActionContext) {
	for i := 0; i < ac.Count; i++ {
		if len(w.GetActiveLine()) == 0 {
			return
		}
		bytesDeleted, lLen, _ := DeleteCharacterAt(w.Buf.Lines[w.CursorLine], w.CursorCol-1)
		w.Buf.Lines[w.CursorLine] = w.Buf.Lines[w.CursorLine][:len(w.Buf.Lines[w.CursorLine])-bytesDeleted]

		w.Logger.Logln("Content After deletion: %s", w.GetActiveLine())
		// w.RedrawLine(w.CursorLine)
		w.NeedRedraw = true
		if w.CursorCol > lLen {
			w.IncrCursorCol(-1)
		}
	}
}

func NormalHandleDeleteChar(w *Window, ac *ActionContext) {
	// w.Logger.Logln("Backspace Detected, content before deletion: %s", w.GetActiveLine())
	if len(w.GetActiveLine()) == 0 {
		return
		//		w.Buf.Lines = DeleteLineAt(w.Buf.Lines, w.CursorLine, 1)
		//		w.IncrCursorLine(-1)
		//		numDisplayedLines := len(w.Buf.Lines) - w.BufTopLine
		//		if numDisplayedLines < w.Height && w.BufTopLine > 0 {
		//			w.BufTopLine--
		//		}
		//		w.NeedRedraw = true
		//		// w.RedrawAllLines()
		//		w.IncrCursorCol(len(w.GetActiveLine()))
	} else {
		// w.Buf.Lines[w.CursorLine] = DeleteByteAt(w.GetActiveLine(), w.CursorCol-1)
		bytesDeleted, _, _ := DeleteCharacterAt(w.Buf.Lines[w.CursorLine], w.CursorCol-1)
		w.Buf.Lines[w.CursorLine] = w.Buf.Lines[w.CursorLine][:len(w.Buf.Lines[w.CursorLine])-bytesDeleted]

		w.Logger.Logln("Content After deletion: %s", w.GetActiveLine())
		// w.RedrawLine(w.CursorLine)
		w.NeedRedraw = true
		w.IncrCursorCol(-1)
	}
}

func NormalModeHandleDeleteCmd(w *Window, ac *ActionContext) {
	switch {
	case string(ac.Suffix) == "d":
		w.Logger.Logln("Deleting cursor line: %d", w.CursorLine)
		w.Buf.Lines = DeleteLineAt(w.Buf.Lines, w.CursorLine, ac.Count)
		if w.CursorLine >= len(w.Buf.Lines) && len(w.Buf.Lines) > 0 {
			w.CursorLine = len(w.Buf.Lines) - 1
		}
		w.NeedRedraw = true
		//		for i := 1; i <= ac.Count; i++ {
		//			if w.CursorLine >= len(w.Buf.Lines) {
		//				w.CursorLine = len(w.Buf.Lines) - 1
		//			}
		//			w.Buf.Lines = DeleteLineAt(w.Buf.Lines, w.CursorLine, 1)
		//			if w.CursorLine >= len(w.Buf.Lines) {
		//				w.CursorLine = len(w.Buf.Lines) - 1
		//			}
		//		}
	}
}

func NormalModeSwitchToCmd(w *Window, ac *ActionContext) {
	GlobalLogger.Logln("Switching mode to command")
	w.CmdBuf[0] = ':'
	w.Mode = MODE_CMD
}

func NormalHandleMoveToLineStart(w *Window, ac *ActionContext) {
	w.IncrCursorCol(-w.CursorCol)
	w.MoveCursorToDisplayPosition()
}

func NormalHandleAppendToEOL(w *Window, ac *ActionContext) {
	curBytePosition, _ := GetNthChar(w.Buf.Lines[w.CursorLine], w.CursorCol)
	charsFromPos := Utf8Len(w.Buf.Lines[w.CursorLine][curBytePosition:])
	w.IncrCursorCol(charsFromPos + 1)
	w.MoveCursorToDisplayPosition()
	NormalModeSwitchToInsert(w, ac)
}
