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
		//		fmt.Printf("data: %s", string(data))
		//		fmt.Printf("Split index: %d\n", splitIndex)
		//		fmt.Printf("atEOF: %t\n", atEOF)
		err = nil
		if splitIndex == -1 {
			if atEOF {
				err = bufio.ErrFinalToken
				if len(data) != 0 {
					token = append(tokenBuffer, data[:]...)
				} else if len(tokenBuffer) > 0 {
					//fmt.Println("Copying data from token buffer")
					token = append(token, tokenBuffer...)
					//fmt.Printf("Token after appending: %s", string(token))
				}
				tokenBuffer = tokenBuffer[:0]

			} else {
				tokenBuffer = append(tokenBuffer, data...)
				buffering = true
			}
			advance = len(data)
		} else {
			//fmt.Println("got split match")
			if buffering {
				token = append(tokenBuffer, data[:splitIndex]...)
				tokenBuffer = tokenBuffer[:0]
				buffering = false
			} else {
				token = data[:splitIndex]
			}
			//fmt.Printf("%s\n", string(token))
			advance = splitIndex + 1
			//fmt.Printf("Advance is : %d\n", advance)
			//fmt.Printf("splitIndex is : %d\n", splitIndex)
		}
		if len(token) < 1 {
			token = nil
		}
		//fmt.Printf("token: %s\n", string(token))
		//fmt.Printf("tokenBuffer: %s\n", string(tokenBuffer))
		return advance, token, err
	}
}

func main() {
	/* Arg spec:
	First arg should be the character to split on (default to ,)
	second arg should be the file to tokenize with the split
	default to local testing file
	*/
	var splitChar []byte
	var filePath string
	if len(os.Args) > 1 {
		splitChar = []byte(os.Args[1])
	} else {
		splitChar = []byte(",")
	}

	if len(os.Args) > 2 {
		filePath = os.Args[2]
	} else {
		filePath = "./split.txt"
	}

	localSplitFunc := createSplitFunc(splitChar)

	fileDescriptor, err := os.Open(filePath)
	if err != nil {
		fmt.Errorf("error opening file %w", err)
	}

	scanner := bufio.NewScanner(fileDescriptor)
	scanner.Split(localSplitFunc)

	for scanner.Scan() {
		token := scanner.Text()
		if len(token) > 0 {
			fmt.Printf("%s\n", scanner.Text())
			//fmt.Printf("%s\n", strings.TrimSpace(scanner.Text()))
		}
	}

}
