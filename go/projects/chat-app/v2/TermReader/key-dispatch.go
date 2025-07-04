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

func ParseMotion(w *Window, ka *KeyAction)

func ParseNextAction(w *Window, ka *KeyAction) *CommandEntry {
	var err error
	switch w.Mode {
	case MODE_NORMAL:
		entry, ok := NormalModeDispatchMap[string(ka.Value)]
		if !ok || (entry.MustBeFirst && len(w.ExpectingInputBuf) > 0) {
			w.ExpectingInputBuf = append(w.ExpectingInputBuf, ka.Value...)
			w.ActionCount, err = strconv.Atoi(string(w.ExpectingInputBuf))
			if err != nil {
				w.DisplayCmdMessage(fmt.Sprintf("Error %s is not a valid input count for NORMAL mode", string(w.ExpectingInputBuf)))
			}
			return nil
		} else {
			cmd := NormalModeDispatchMap[string(ka.Value)]
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

func PopulateActionContext(action *Action, buf []byte) {
}

func MakeNewActionContextPtr() *ActionContext {
	return &ActionContext{
		Count:     0,
		FullInput: make([]byte, 0, 10),
		Prefix:    make([]byte, 0, 10),
		Suffix:    make([]byte, 0, 10),
	}
}
