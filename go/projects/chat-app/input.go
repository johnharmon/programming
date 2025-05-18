package main

import (
	"bytes"
	"io"
)

type InputScanner struct {
	Remaining   []byte
	LastMessage []byte
	LastChunk   []byte
	Input       io.Reader
	Delimiter   byte
}

func (i *InputScanner) Scan() {
	buf := make([]byte, 1024)
	i.Input.Read(buf)
	i.ScanInput(buf)
}

func (i *InputScanner) ScanInput(input []byte) (message []byte, messageRead bool) {
	remaining := []byte{}

	delimIndex := bytes.IndexByte(input, i.Delimiter)
	if delimIndex != -1 {
		message = append(i.Remaining, input[:delimIndex]...)
		remaining = input[delimIndex+1:]
		messageRead = true
	} else {
		message = append(i.Remaining, input...)
		remaining = input[:0]
		messageRead = false
	}
	i.Remaining = remaining
	if len(message) > 0 {
		i.LastMessage = message
	}
	i.LastChunk = input
	return message, messageRead
}

func ReadUntil(input io.Reader, delim byte) (message []byte, err error) {
	return message, err
}

func ScanInput(input []byte, delim byte) (message []byte, remaining []byte, isTerminated, err error) {
	delimIndex := bytes.IndexByte(input, delim)
	if delimIndex != -1 {
		message = input[:delimIndex]
		remaining = input[delimIndex+1:]
	} else {
		message = input
		remaining = input[:0]
	}
	return message, remaining, isTerminated, err
}

func NewInputScanner(input io.Reader) {
	is := InputScanner{}
	is.Remaining = []byte
	is.LastMessage = []byte
	is.LastChunk = []byte
	is.Input = os.Stdin
	is.Delimiter = []byte("\n")
	return &is
}

func RunScanner() {
	scanner := NewInputScanner()
	for {
		scanner.Scan()
	}
}

func main() {
	RunScanner()
}
