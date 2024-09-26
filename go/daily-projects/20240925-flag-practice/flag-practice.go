package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

const NewLineByte = byte(012)

func writeTaskFile(taskFile *os.File, content []byte) {
	taskFile.Seek(0, 0)
	taskFile.Truncate(0)
	taskFile.Write(content)
}

func addLine(content []byte, newLine []byte) (newContent []byte) {
	if content[len(content)-1] != NewLineByte {
		newLine = append([]byte{NewLineByte}, newLine...)
	}
	content = append(content, newLine...)
	return content
}

func createTaskFile(taskFile string) (taskFileObj *os.File, funcErr error) { // returns an open file, be sure to close it :)
	fileInfo, err := os.Stat(taskFile)
	if err != nil {
		if os.IsNotExist(err) {
			_, createErr := os.Create(taskFile)
			if createErr != nil {
				return nil, fmt.Errorf("error creating new task file: %w", createErr)
			}
			fileInfo, err = os.Stat(taskFile)
			if err != nil {
				return nil, fmt.Errorf("error getting file stat: %w", err)
			}
		} else {
			return nil, fmt.Errorf("error getting file stat: %w", err)
		}
	}
	taskFileObj, openErr := os.OpenFile(taskFile, os.O_RDWR|os.O_CREATE, fileInfo.Mode())
	if openErr != nil {
		return nil, fmt.Errorf("error opening created file: %w", openErr)
	}
	return taskFileObj, nil
}

func processLine(line []byte) (newLine []byte, doPrint bool) {
	lineString := string(line)
	lineString = strings.TrimSpace(lineString)
	if lineString != "" { // only care about non empty lines, will NOT write empty lines back to the file
		newLine = append(line, NewLineByte)                                            // append new line
		if strings.HasPrefix(lineString, "//") || strings.HasPrefix(lineString, "#") { // keep comments in file, but do not print them
			return newLine, false
		} else {
			return newLine, true
		}
	} else {
		return nil, false
	}
}

func printTasks(taskFile *os.File) error {
	scanner := bufio.NewScanner(taskFile)
	taskNumber := 0
	newLines := []byte{}

	for scanner.Scan() {
		newLine, doPrint := processLine(scanner.Bytes())
		if newLine != nil {
			newLines = append(newLines, newLine...)
		}
		if doPrint {
			taskNumber++
			fmt.Printf("%d.  %s", taskNumber, newLine)
		}
	}
	scanError := scanner.Err()
	if scanError != nil {
		return fmt.Errorf("scanning file encountered error: %s", scanError)
	}
	taskFile.Seek(0, 0)
	taskFile.Truncate(0)
	taskFile.Write(newLines)
	return nil
}

func cleanExtraLines(content []byte) []byte {
	var previousByte byte
	cleanedBytes := []byte{}
	for _, val := range content {
		if previousByte != NewLineByte {
			cleanedBytes = append(cleanedBytes, val)
		}
		previousByte = val
	}
	return cleanedBytes
}

func cleanFile(taskFile *os.File) []byte {
	scanner := bufio.NewScanner(taskFile)
	newLines := []byte{}
	for scanner.Scan() {
		newLine, _ := processLine(scanner.Bytes())
		if newLine != nil {
			newLines = append(newLines, newLine...)
		}
	}
	return newLines
}

func addTask(taskFile *os.File, taskString string) (content []byte) {
	taskBytes := []byte(taskString)
	cleanedFile := cleanFile(taskFile)
	newFileBytes := addLine(cleanedFile, taskBytes)
	return newFileBytes
}

func main() {
	taskFile := "./.taskfile"
	taskFileObj, err := createTaskFile(taskFile)
	if err != nil {
		fmt.Printf("Error creating task file: %s\n", err)
		defer taskFileObj.Close()
		os.Exit(1)
	}
	defer taskFileObj.Close()
	var (
		list    bool
		add     string
		remove  int
		sremove string
	)

	flag.BoolVar(&list, "list", false, "List the tasks on the list") //Create list boolean flag
	flag.StringVar(&add, "add", "", "Add a task to the list")
	flag.IntVar(&remove, "remove", 0, "Remove a task from the list (base 1 indexed)")
	flag.StringVar(&sremove, "remove_name", "", "Remove a task by name")
	flag.Parse()
	//	for index, value := range flag.Args() {
	//		fmt.Printf("Argument #%d was: %s\n", index, value)
	//	}
	//	flag.VisitAll(func(f *flag.Flag) {
	//		fmt.Printf("%s: %s\n", f.Name, f.Value)
	//	})

	if list {
		printTasks(taskFileObj)
	} else if add != "" {
		addTask(taskFileObj, add)
	}

}
