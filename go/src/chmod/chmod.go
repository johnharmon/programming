package main

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"syscall"
)

func ErrorWrap(err error, message string) error {
	return fmt.Errorf("%s: %w", message, err)
}

func GetFileInfo(filepath string) (*syscall.Stat_t, *os.FileInfo, error) {
	var result_error error
	file_info, err := os.Stat(filepath)
	if err != nil {
		result_error = fmt.Errorf("error: %s", err)
	}
	fileInfo := file_info.Sys().(*syscall.Stat_t)
	return fileInfo, &file_info, result_error
}

func Chmod(filepath string, mode os.FileMode) error {
	fs.FileMode
	return nil
}

func ConvertFileMode(mode string) (os.FileMode, error) {
	return 0, nil
}

func main() {
	my_err := errors.New("beep")
	fmt.Println(my_err)
	fmt.Printf("Placeholder\n")
}

/*
var (
	patherr *os.PathError
	patherrptr = &patherr
)
_, err := someFunc()

if errors.As(err, patherrptr) {
	fmt.Println("Error is a path error")
}
*/
