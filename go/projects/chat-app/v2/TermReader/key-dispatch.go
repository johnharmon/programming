package main

import (
	"fmt"
	"strconv"
)

//func InitNormalModeDispatchMap() {
//	dispatchMap := NormalModeDispatchMap
//	dispatchMap["h"] = Action{
//		ExecFunc: NormalHandleDownMove,
//	}
//}

//func (w *Window) ParseNextAction(ka *KeyAction) (actionReady bool) {
//}

//func MakeActionParser() (func (*Window, *KeyAction) *CommandEntry){
//	var(
//		motionParseed = false
//		actionCount = 0
//		motionCount = 0
//
//	)
//	return func(w *Window, ka *KeyAction) *CommandEntry {
//		var err error
//		switch w.Mode {
//		case MODE_NORMAL:
//			entry, ok := NormalModeDispatchMap[string(ka.Value)]
//			if !ok || (entry.MustBeFirst && len(w.ExpectingInputBuf) > 0) {
//				w.ExpectingInputBuf = append(w.ExpectingInputBuf, ka.Value...)
//				w.ActionCount, err = strconv.Atoi(string(w.ExpectingInputBuf))
//				if err != nil {
//					w.DisplayCmdMessage(fmt.Sprintf("Error %s is not a valid input count for NORMAL mode", string(w.ExpectingInputBuf)))
//				}
//				return nil
//			} else {
//				cmd := NormalModeDispatchMap[string(ka.Value)]
//				return &cmd
//			}
//			//			awc := &ActionWithContext{Action: &action, ActionContext: MakeNewActionContextPtr()}
//			//			PopulateActionContext(&action, w.ExpectingInputBuf)
//			//			return awc
//
//		default:
//			return nil
//		}
//}

func ParseMotion(w *Window, ka *KeyAction)

func ParseNextAction(w *Window, ka *KeyAction) *CommandEntry {
	var err error
	switch w.Mode {
	case MODE_NORMAL:
		entry, ok := NormalModeDispatchMap[int(ka.Value[0])]
		if !ok || (entry.MustBeFirst && len(w.ExpectingInputBuf) > 0) {
			w.ExpectingInputBuf = append(w.ExpectingInputBuf, ka.Value...)
			w.ActionCount, err = strconv.Atoi(string(w.ExpectingInputBuf))
			if err != nil {
				w.DisplayCmdMessage(fmt.Sprintf("Error %s is not a valid input count for NORMAL mode", string(w.ExpectingInputBuf)))
			}
			return nil
		} else {
			cmd := NormalModeDispatchMap[int(ka.Value[0])]
			return &cmd
		}
		//			awc := &ActionWithContext{Action: &action, ActionContext: MakeNewActionContextPtr()}
		//			PopulateActionContext(&action, w.ExpectingInputBuf)
		//			return awc

	default:
		return nil
	}
}

/*
	func MakeActionParser() func(*Window, *KeyAction) *Action {
		var ac *ActionContext = MakeNewActionContextPtr()
		return func(w *Window, ka *KeyAction) *Action {
		}
	}

	func MakeNewActionContext() ActionContext {
		return ActionContext{
			Count:     0,
			FullInput: make([]byte, 0, 10),
			Prefix:    make([]byte, 0, 10),
			Suffix:    make([]byte, 0, 10),
		}
	}
*/

func PopulateActionContext(state *NormalModeParsingState, context *ActionContext) {
}

func MakeNewActionContextPtr() *ActionContext {
	return &ActionContext{
		Count:     0,
		FullInput: make([]byte, 0, 10),
		Prefix:    make([]byte, 0, 10),
		Suffix:    make([]byte, 0, 10),
	}
}

func ClearActionContext(ac *ActionContext) {
	ac.Count = 0
	ac.FullInput = ac.FullInput[:0]
	ac.Prefix = ac.Prefix[:0]
	ac.Suffix = ac.Suffix[:0]
	return
}

func CleanParsingState(n *NormalModeParsingState) {
	clear(n.RawInput)
	n.ParsingCount = false
	n.ParsingMotion = false
	n.CommandIdentified = false
	n.Command = nil
	n.Motion = nil
	n.CommandCount = 1
	n.MotionCount = 1
	n.State = 0
	n.Suffix = ""
	n.ExecReady = false
	ClearActionContext(&n.ActionContext)
	return
}

func NormalModeParseNextInput(w *Window, ka *KeyAction) (callAgain bool) {
	var command CommandEntry
	var ok bool
	switch w.NormPS.State {
	case STATE_INITIAL_INPUT:
		if ka.Value[0] >= 0x31 && ka.Value[0] <= 0x39 {
			w.NormPS.CommandCount = int(ka.Value[0] - 0x30)
			w.NormPS.RawInput = append(w.NormPS.RawInput, ka.Value[0])
			w.NormPS.State = STATE_PARSING_CMD_COUNT
		} else {
			command, ok = NormalModeDispatchMap[int(ka.Value[0])]
			if !ok {
				CleanParsingState(&w.NormPS)
				return false
			}
			w.NormPS.State = STATE_CMD_IDENTIFIED
			return true
			//			if !command.SuffixRequired {
			//				actionContext := MakeNewActionContextPtr()
			//				PopulateActionContext(&w.NormPS, actionContext)
			//				command.ExecFunc(w, actionContext)
			//			}
		}
	case STATE_PARSING_CMD_COUNT:
		if ka.Value[0] >= 0x31 && ka.Value[0] <= 0x39 {
			w.NormPS.CommandCount = (w.NormPS.CommandCount * 10) + int((ka.Value[0] - 0x30))
		} else {
			command, ok = NormalModeDispatchMap[int(ka.Value[0])]
			if !ok {
				CleanParsingState(&w.NormPS)
				return false
			} else {
				w.NormPS.Command = &command
				w.NormPS.State = STATE_CMD_IDENTIFIED
				return true
			}
		}
	case STATE_CMD_IDENTIFIED:
		if command.SuffixRequired {
			w.NormPS.State = STATE_PENDING_SUFFIX
			// return true
		} else {
			w.NormPS.ActionContext = *MakeNewActionContextPtr()
			PopulateActionContext(&w.NormPS, &w.NormPS.ActionContext)
			GlobalLogger.Logln("normal mode function ready to execute: %s", command.Name)
			w.NormPS.ExecReady = true
			return false
			// command.ExecFunc(w, actionContext)
		}
	case STATE_PENDING_SUFFIX:
		if command.AcceptsMotion {
			w.NormPS.State = STATE_PARSING_MOTION
			return false
		} else if command.AcceptsSpecialSuffix {
			w.NormPS.State = STATE_PARSING_SPECIAL_SUFFIX
			return false
		}
	case STATE_PARSING_MOTION:
		if ka.Value[0] >= 0x31 && ka.Value[0] <= 0x39 {
			w.NormPS.MotionCount = (w.NormPS.MotionCount * 10) + int((ka.Value[0] - 0x30))
		} else {
			motion := MotionDispatchMap[string(ka.Value)]
			w.NormPS.Motion = &motion
			w.NormPS.ExecReady = true
			return false
			// intIndex, intOffset := motion.MotionFunc(w.Buf.Lines[w.CursorLine], w.CursorCol-1)
		}
	case STATE_PARSING_SPECIAL_SUFFIX:
		w.NormPS.Suffix = string(ka.Value[0])
		w.NormPS.ExecReady = true
		return false

	default:
		return false
	}
	return false
}
