package main

import (
	"errors"
	"os"
)

func Edit(w *Window, args ...string) error {
	var (
		f         *os.File
		fileFlags int
		filepath  string
	)
	if len(args) > 0 {
		filepath = args[0]
	} else {
		filepath = w.FileName
	}
	if filepath == "" {
		w.DisplayCmdMessage("Error, no file provided to open!")
		return errors.New("Err: no such file")
	}
	_, err := os.Stat(filepath)
	if err != nil {
		GlobalLogger.Logln("Opening new file: %s", filepath)
		fileFlags = os.O_RDONLY | os.O_CREATE
	} else {
		GlobalLogger.Logln("Opening existing file: %s", filepath)
		fileFlags = os.O_RDONLY
	}
	f, err = os.OpenFile(filepath, fileFlags, 0o644)
	if err != nil {
		GlobalLogger.Logln("Error opening file: %s", err)
		return err
	}
	w.LoadBuffer(f)
	w.FileName = f.Name()
	w.CursorCol, w.CursorLine = 1, 0
	w.TermTopLine = 1
	w.BufTopLine = 0
	w.NeedRedraw = true
	return nil
}

func Write(w *Window, args ...string) error {
	var filepath string
	if len(args) > 0 {
		filepath = args[0]
	} else {
		if w.FileName != "" {
			filepath = w.FileName
		} else {
			filepath = "tmt.txt"
		}
	}
	GlobalLogger.Logln("Opening and writing to file: %s", filepath)
	f, err := os.OpenFile(filepath, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	w.WriteBuffer(f)
	return nil
}
