package main

//func InitNormalModeDispatchMap() {
//	dispatchMap := NormalModeDispatchMap
//	dispatchMap["h"] = Action{
//		ExecFunc: NormalHandleDownMove,
//	}
//}
/*

func (w *Window) ParseNextAction(ka *KeyAction) (actionReady bool) {
}

func ParseNextAction(w *Window, ka *KeyAction) *Action {
	switch w.Mode {
	case MODE_NORMAL:
		action, ok := NormalModeDispatchMap[ka.Action]
		switch {
		case len(
		}
		if !ok || (ok && action.MustBeFirst) {
			w.ExpectingInputBuf = append(w.ExpectingInputBuf, ka.Value...)
			return nil
		} else {
			w.ExpectingInputBuf = append(w.ExpectingInputBuf, ka.Value...)
			action.Context = MakeNewActionContextPtr()
			action.Context.Count = len(w.ExpectingInputBuf)
		}
	default:
		return nil
	}
}

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

func MakeNewActionContextPtr() *ActionContext {
	return &ActionContext{
		Count:     0,
		FullInput: make([]byte, 0, 10),
		Prefix:    make([]byte, 0, 10),
		Suffix:    make([]byte, 0, 10),
	}
}
*/
