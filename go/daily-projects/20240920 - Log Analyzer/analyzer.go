package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"sort"
)

func AnalyzeFile(filepath string) (map[string]int, error) {
	debug := false
	file, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("error opening the file: %w", err)
	}
	message_counts := map[string]int{
		"INFO":  0,
		"WARN":  0,
		"ERROR": 0,
		"FATAL": 0,
		"DEBUG": 0,
	}

	scanner := bufio.NewScanner(file)

	if debug {
		fmt.Printf("Starting to scan %s\n", filepath)
	}

	for scanner.Scan() {

		if debug {
			fmt.Println("Scanning line")
		}

		line := scanner.Bytes()
		header_end := bytes.Index(line, []byte("\072"))
		line_start := line[:header_end]
		line_start_string := string(line_start)

		if debug {
			fmt.Println(line_start_string)
		}

		switch {
		case line_start_string == "INFO":
			message_counts["INFO"]++
		case line_start_string == "WARN":
			message_counts["WARN"]++
		case line_start_string == "DEBUG":
			message_counts["DEBUG"]++
		case line_start_string == "FATAL":
			message_counts["FATAL"]++
		case line_start_string == "ERROR":
			message_counts["ERROR"]++
		}
		if err := scanner.Err(); err != nil {
			return message_counts, err
		}
	}
	return message_counts, nil
}

func PrintAnalysis(message_counts map[string]int) {
	keys_alphabetical := make([]string, 0, len(message_counts))
	for key := range message_counts {
		keys_alphabetical = append(keys_alphabetical, key)
	}
	sort.Strings(keys_alphabetical)
	for _, key := range keys_alphabetical {
		fmt.Printf("%s Messages: %d\n", key, message_counts[key])
	}
}

func main() {
	var filepath string

	if len(os.Args) > 1 {
		filepath = os.Args[1]
	} else {
		filepath = "log.txt"
	}
	fmt.Printf("Analyzing file: %s\n", filepath)
	message_counts, err := AnalyzeFile(filepath)
	if err != nil {
		fmt.Printf("Error analyzing file: %s\n", err)
		fmt.Printf("However, before running into an error, here are the counts we got for each message level:\n")
		PrintAnalysis(message_counts)
	} else {
		PrintAnalysis(message_counts)
	}
}
