package main

import (
	"errors"
	"os"
)

func Edit(w *Window, args ...string) error {
	if len(args) > 0 {
		filepath := args[0]
		var f *os.File
		if filepath == "" {
			w.DisplayCmdMessage("Error, no file provided to open!")
			return errors.New("err: no such file")
		}
		_, err := os.Stat(filepath)
		if err != nil {
			GlobalLogger.Logln("Opening new file: %s", filepath)
			f, err = os.OpenFile(filepath, os.O_RDWR|os.O_CREATE, 0o644)
			if err != nil {
				GlobalLogger.Logln("Error opening file: %s", err)
				return err
			}
		} else {
			f, err = os.OpenFile(filepath, os.O_RDWR, 0o644)
			GlobalLogger.Logln("Opening existing file: %s", filepath)
			if err != nil {
				GlobalLogger.Logln("Error opening file: %s", err)
				return err
			}
		}
		w.LoadBuffer(f)
		w.CursorCol, w.CursorLine = 1, 0
		w.TermTopLine = 1
		w.BufTopLine = 0
		w.RedrawAllLines()
	}
	return nil
}
