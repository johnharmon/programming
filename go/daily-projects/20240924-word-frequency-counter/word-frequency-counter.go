package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
)

func createSplitFunc(splitChar []byte) func(data []byte, atEOF bool) (advance int, token []byte, err error) {
	tokenBuffer := make([]byte, 1024)
	buffering := false
	return func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		splitIndex := bytes.Index(data, splitChar)
		err = nil
		if splitIndex == -1 {
			if atEOF {
				err = bufio.ErrFinalToken
				if len(data) != 0 {
					token = append(tokenBuffer, data[:]...)
				} else if len(tokenBuffer) > 0 {
					token = append(token, tokenBuffer...)
				}
				tokenBuffer = tokenBuffer[:0]

			} else {
				tokenBuffer = append(tokenBuffer, data...)
				buffering = true
			}
			advance = len(data)
		} else {
			if buffering {
				token = append(tokenBuffer, data[:splitIndex]...)
				tokenBuffer = tokenBuffer[:0]
				buffering = false
			} else {
				token = data[:splitIndex]
			}
			advance = splitIndex + 1
		}
		if len(token) < 1 {
			token = nil
		}
		return advance, token, err
	}
}

func getWordFrequency(filepath string) (map[string]int, error) {
	wordFrequency := make(map[string]int)
	wordSplitFunc := createSplitFunc([]byte("\040"))
	fileObj, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %w", err)
	}
	defer fileObj.Close()
	scanner := bufio.NewScanner(fileObj)
	scanner.Split(wordSplitFunc)
	for scanner.Scan() {
		wordFrequency[scanner.Text()]++
	}
	if err := scanner.Err(); err != nil {
		scanError := fmt.Errorf("error from scanner: %w", err)
		return wordFrequency, scanError
	}
	return wordFrequency, nil
}

func main() {
	var filepath string
	filepath = "./words.txt"
	getWordFrequency(filepath)
}
