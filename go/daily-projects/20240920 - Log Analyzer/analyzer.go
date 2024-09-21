package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"sort"
)

func AnalyzeFile(filePath string) (map[string]int, error) {
	debug := false
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("error opening the file: %w", err)
	}
	messageCounts := map[string]int{
		"INFO":  0,
		"WARN":  0,
		"ERROR": 0,
		"FATAL": 0,
		"DEBUG": 0,
	}

	scanner := bufio.NewScanner(file)

	if debug {
		fmt.Printf("Starting to scan %s\n", filePath)
	}

	for scanner.Scan() {

		if debug {
			fmt.Println("Scanning line")
		}

		line := scanner.Bytes()
		headerEnd := bytes.Index(line, []byte(":"))
		if headerEnd == -1 {
			continue
		}
		lineStart := line[:headerEnd]
		lineStartString := string(lineStart)

		if debug {
			fmt.Println(lineStartString)
		}

		switch {
		case lineStartString == "INFO":
			messageCounts["INFO"]++
		case lineStartString == "WARN":
			messageCounts["WARN"]++
		case lineStartString == "DEBUG":
			messageCounts["DEBUG"]++
		case lineStartString == "FATAL":
			messageCounts["FATAL"]++
		case lineStartString == "ERROR":
			messageCounts["ERROR"]++
		}
	}
	if err := scanner.Err(); err != nil {
		return messageCounts, err
	} else {
		return messageCounts, nil
	}
}

func PrintAnalysis(messageCounts map[string]int) {
	keysAlphabetical := make([]string, 0, len(messageCounts))
	for key := range messageCounts {
		keysAlphabetical = append(keysAlphabetical, key)
	}
	sort.Strings(keysAlphabetical)
	for _, key := range keysAlphabetical {
		fmt.Printf("%s Messages: %d\n", key, messageCounts[key])
	}
}

func main() {
	var filePath string

	if len(os.Args) > 1 {
		filePath = os.Args[1]
	} else {
		filePath = "log.txt"
	}
	fmt.Printf("Analyzing file: %s\n", filePath)
	messageCounts, err := AnalyzeFile(filePath)
	if err != nil {
		fmt.Printf("Error analyzing file: %s\n", err)
		fmt.Printf("However, before running into an error, here are the counts we got for each message level:\n")
	}
	PrintAnalysis(messageCounts)
}
