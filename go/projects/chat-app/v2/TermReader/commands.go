package main

func Edit(w *Window, args ...string) error {
	cmdRaw := w.CmdBuf
	cmdArgs := ProcessCmdArgs(cmdRaw)
	return nil
}
