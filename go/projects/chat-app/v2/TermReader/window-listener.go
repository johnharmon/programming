package main

import (
	"fmt"
	"os"
)

func GetModeString(mode int) string {
	switch mode {
	case 0:
		return "NORMAL"
	case 1:
		return "INSERT"
	case 2:
		return "VISUAL"
	case 3:
		return "CMD"
	default:
		return "UNKNOWN"
	}
}

func (w *Window) Listen() {
	// redrawHandler := w.MakeRedrawHandler()
	gl := GlobalLogger
	//	expectingInput := false
	//	var expectingInputFunc func(*Window, *KeyAction) bool
	var callAgain bool
	w.Logger = gl
	w.Out.Write(TERM_CLEAR_SCREEN)
	w.DisplayStatusLine()
	w.MoveCursorToPosition(1, 1)
	var ka *KeyAction
	for {
		GlobalLogger.Logln("########## START OF WINDOW LISTEN LOOP ##########")
		GlobalLogger.Logln("Window Mode: %s", GetModeString(w.Mode))
		ka = <-w.EventChan
		gl.Logln("Window received *KeyAction: %s", ka.String())
	ModeSwitch:
		switch w.Mode {
		case MODE_INSERT:
			if ka.Action == "Print" && len(ka.Value) == 1 {
				gl.Logln("Raw write triggered for %s", ka.String())
				w.WriteRaw(ka.Value)
				IncrTwoCursorColPtr(w.Buf.Lines[w.CursorLine], &w.CursorCol, &w.DesiredCursorCol, 1)
				GlobalLogger.Logln("New Cursor Col after ptr mutation: %d", w.CursorCol)
				// w.IncrCursorCol(1)
				w.RedrawLine(w.Buf.ActiveLine)
				w.MoveCursorToDisplayPosition()
				// w.KeyActionReturner <- ka
			} else {
				switch ka.Action {
				case "Backspace":
					w.Logger.Logln("Backspace Detected, content before deletion: %s", w.GetActiveLine())
					w.Buf.Lines[w.Buf.ActiveLine] = DeleteByteAt(w.Buf.Lines[w.Buf.ActiveLine], w.CursorCol-1)
					w.IncrCursorLine(-1)
					w.Logger.Logln("Content After deletion: %s", w.GetActiveLine())
					w.RedrawLine(w.CursorLine)
					// w.IncrCursorCol(-1)
					IncrCursorColPtr(w.Buf.Lines[w.CursorLine], &w.CursorCol, -1)
					GlobalLogger.Logln("New Cursor Col after ptr mutation: %d", w.CursorCol)
				case "Delete":
					InsertHandleDelete(w)
				case "ArrowRight":
					InsertHandleArrowRight(w)
				case "ArrowLeft":
					InsertHandleArrowLeft(w)
				case "ArrowUp":
					InsertHandleArrowUp(w)
				case "ArrowDown":
					InsertHandleArrowDown(w)
				case "Enter":
					InsertHandleEnter(w)
				case "Escape":
					w.Mode = MODE_NORMAL
					GlobalLogger.Logln("Setting mode to normal")
					break ModeSwitch

				}
			}
		case MODE_NORMAL:
			callAgain = NormalModeParseNextInput(w, ka)
			for callAgain {
				GlobalLogger.Logln("callAgain invoked on the normal mode command parser")
				callAgain = NormalModeParseNextInput(w, ka)
			}

			if w.NormPS.ExecReady {
				fmt.Fprintf(os.Stderr, "Executing command: +%v\n", w.NormPS.Command)
				w.NormPS.Command.ExecFunc(w, &w.NormPS.ActionContext)
				CleanParsingState(&w.NormPS)
				ClearActionContext(&w.NormPS.ActionContext)
			}
			/*
				if expectingInput {
					expectingInput = expectingInputFunc(w, ka)
				} else if ka.Action == "Print" {
					switch ka.Value[0] {
					case CHAR_h:
						NormalHandleLeftMove(w, 1)
					case CHAR_j:
						NormalHandleDownMove(w, 1)
					case CHAR_k:
						NormalHandleUpMove(w, 1)
					case CHAR_l:
						NormalHandleRightMove(w, 1)
					case CHAR_i:
						w.Mode = MODE_INSERT
					case CHAR_f:
						expectingInput = true
						expectingInputFunc = NormalHandleForwardFind
					case CHAR_COLON:
						GlobalLogger.Logln("Setting mode to cmd")
						w.PrevCursorCol = w.CursorCol
						w.CursorCol = w.CmdCursorCol
						w.Mode = MODE_CMD
						// w.CursorCol = 2
						w.CmdBuf[0] = ':'
					default:
						break
					}
				} else {
					switch ka.Action {
					case "ArrowRight":
						NormalHandleArrowRight(w)
					case "ArrowLeft":
						NormalHandleArrowLeft(w)
					case "ArrowUp":
						NormalHandleArrowUp(w)
					case "ArrowDown":
						NormalHandleArrowDown(w)
					case "Enter":
						NormalHandleEnter(w)
					}
				}
			*/
		case MODE_CMD:
			GlobalLogger.Logln("CMD INPUT PARSING STARTED")
			switch {
			case ka.Action == "Escape":
				w.Mode = MODE_NORMAL
				w.CmdCursorCol = w.CursorCol
				w.CursorCol = w.PrevCursorCol
				w.MoveCursorToDisplayPosition()
			case ka.Action == "Enter":
				_ = w.ProcessCmd()
				w.Mode = MODE_NORMAL
			case ka.Action == "Delete":
				CmdHandleDelete(w)
			case ka.Action == "ArrowRight":
				CmdHandleArrowRight(w)
			case ka.Action == "ArrowLeft":
				CmdHandleArrowLeft(w)
			case ka.PrintRaw && len(ka.Value) == 1:
				w.WriteToCmd(ka.Value)
				IncrCursorColPtr(w.CmdBuf, &w.CmdCursorCol, 1)

			}
		case MODE_VISUAL:
			continue
		}
		switch {
		case w.Mode == MODE_CMD:
			w.DisplayCmdLine()
		default:
			w.DisplayStatusLine()
			// w.MoveCursorToDisplayPosition()
		}
		switch {
		case w.Mode == MODE_NORMAL || w.Mode == MODE_INSERT || w.Mode == MODE_VISUAL:
			w.MoveCursorToDisplayPosition()
		case w.Mode == MODE_CMD:
			w.DisplayCmdLine()
			w.MoveCursorToCmdPosition()
		}
		// w.RedrawAllLines()
		if w.NeedRedraw {
			w.RedrawAllLines()
			w.NeedRedraw = false
			w.MoveCursorToDisplayPosition()
		}
		if ka.FromPool {
			w.KeyActionReturner <- ka
		}
		GlobalLogger.Logln("########## END OF WINDOW LISTEN LOOP ##########")
	}
}
