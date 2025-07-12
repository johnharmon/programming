package main

import (
	"fmt"
	"os"
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

func ParseMotion(w *Window, ka *KeyAction) {}

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
	context.Count = state.CommandCount
	context.Suffix = make([]byte, len(state.Suffix))
	copy(context.Suffix, state.Suffix)
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
	GlobalLogger.Logln("Clearing action context")
	ac.Count = 0
	ac.FullInput = ac.FullInput[:0]
	ac.Prefix = ac.Prefix[:0]
	ac.Suffix = ac.Suffix[:0]
	return
}

func CleanParsingState(n *NormalModeParsingState) {
	GlobalLogger.Logln("Cleaning parsing state")
	clear(n.RawInput)
	n.ParsingCount = false
	n.ParsingMotion = false
	n.CommandIdentified = false
	n.Command = nil
	n.Motion = nil
	n.CommandCount = 1
	n.MotionCount = 1
	n.State = STATE_INITIAL_INPUT
	n.Suffix = []byte{}
	n.ExecReady = false
	ClearActionContext(&n.ActionContext)
	return
}

func NormalModeParseNextInput(w *Window, ka *KeyAction) (callAgain bool) {
	var command CommandEntry
	var ok bool
	switch w.NormPS.State {
	case STATE_INITIAL_INPUT:
		GlobalLogger.Logln("Initial input state entered for parsing")
		if ka.Value[0] >= 0x31 && ka.Value[0] <= 0x39 {
			w.NormPS.CommandCount = int(ka.Value[0] - 0x30)
			w.NormPS.RawInput = append(w.NormPS.RawInput, ka.Value[0])
			w.NormPS.State = STATE_PARSING_CMD_COUNT
			return false
		} else {
			command, ok = NormalModeDispatchMap[int(ka.Value[0])]
			if !ok {
				GlobalLogger.Logln("Command not found: %d", int(ka.Value[0]))
				CleanParsingState(&w.NormPS)
				return false
			}
			GlobalLogger.Logln("Command name: %s", command.Name)
			w.NormPS.State = STATE_CMD_IDENTIFIED
			w.NormPS.Command = &command
			return true
			//			if !command.SuffixRequired {
			//				actionContext := MakeNewActionContextPtr()
			//				PopulateActionContext(&w.NormPS, actionContext)
			//				command.ExecFunc(w, actionContext)
			//			}
		}
	case STATE_PARSING_CMD_COUNT:
		if ka.Value[0] >= 0x30 && ka.Value[0] <= 0x39 {
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
		GlobalLogger.Logln("Command Parser state: STATE_CMD_IDENTIFIED")
		fmt.Fprintf(os.Stderr, "SuffixRequired: %t", w.NormPS.Command.SuffixRequired)
		GlobalLogger.Logln("SuffixRequired: %t", w.NormPS.Command.SuffixRequired)
		if w.NormPS.Command.SuffixRequired {
			GlobalLogger.Logln("Command suffix required, waiting for more input")
			w.NormPS.State = STATE_PENDING_SUFFIX
			return true
		} else {
			w.NormPS.ActionContext = *MakeNewActionContextPtr()
			PopulateActionContext(&w.NormPS, &w.NormPS.ActionContext)
			GlobalLogger.Logln("normal mode function ready to execute: %s", w.NormPS.Command.Name)
			w.NormPS.ExecReady = true
			return false
			// command.ExecFunc(w, actionContext)
		}
	case STATE_PENDING_SUFFIX:
		if w.NormPS.Command.AcceptsMotion {
			w.NormPS.State = STATE_PARSING_MOTION
			return false
		} else if w.NormPS.Command.AcceptsSpecialSuffix {
			GlobalLogger.Logln("Setting state to STATE_PARSING_SPECIAL_SUFFIX")
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
		w.NormPS.Suffix = ka.Value
		PopulateActionContext(&w.NormPS, &w.NormPS.ActionContext)
		w.NormPS.ExecReady = true
		return false

	default:
		return false
	}
	return false
}
