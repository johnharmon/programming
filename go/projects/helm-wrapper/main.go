package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
)

func logger(logs chan []byte, output io.Writer) {

}

func WrapLine(lineNumber int, line []byte, output *bytes.Buffer, logs io.Writer) {
	wrappingCheckRegex := regexp.MustCompile(`^.*{{\s*"{{(hub)?-?"\s*}}.*`)
	fmt.Fprintf(logs, "##########		Line: %d		##########\n", lineNumber)
	if !wrappingCheckRegex.Match(line) {
		regex := regexp.MustCompile(`^(\s*)({{)(hub)?(-)?(.*)(-)?(hub)?(}})(\s*$)`)
		newLine := regex.ReplaceAll(line, []byte(`${1}{{ "${2}${3}${4}" }}${5}{{ "${6}${7}${8}" }}${9}`))
		output.Write(newLine)
		fmt.Fprintf(logs, "%s    ->     %s\n", line, newLine)
		return
	}
	fmt.Fprintf(logs, "##########		UNCHANGED		##########\n")
}

func WrapHubLine(line []byte, output *bytes.Buffer) {
	regex := regexp.MustCompile(`^(\s*)({{hub-?)(.*)(-?hub}})(\s*$)`)
	newLine := regex.ReplaceAll(line, []byte(`${1}{{ "${2}" }}${3}{{ "${4}" }}${5}`))
	output.Write(newLine)
}

func IsMeta(line []byte) bool {
	regex := regexp.MustCompile(`^\s*#\s*meta\s*$`)
	return regex.Match(line)
}

func WrapIndentedLine(line []byte, output *bytes.Buffer) {
	wrappingCheckRegex := regexp.MustCompile(`^.*{{\s*"{{(hub)?-?"\s*}}.*`)
	if !wrappingCheckRegex.Match(line) {
		regex := regexp.MustCompile(`^(\s+)({{)(hub)?(-)?(.*)(-)?(hub)?(}})(\s*$)`)
		newLine := regex.ReplaceAll(line, []byte(`${1}{{ "${2}${3}${4}" }}${5}{{ "${6}${7}${8}" }}${9}`))
		output.Write(newLine)
	}
}

func WrapIndentedHubLine(line []byte, output *bytes.Buffer) {
	regex := regexp.MustCompile(`^(\s+)({{hub-?)(.*)(-?hub}})(\s*$)`)
	newLine := regex.ReplaceAll(line, []byte(`${1}{{ "${2}" }}${3}{{ "${4}" }}${5}`))
	output.Write(newLine)
}

func OpenFile(filePath string) (*os.File, error) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		pathErr := err.(*os.PathError)
		fmt.Printf("Error opening file: %s\n\"%s\"\n", pathErr.Path, pathErr.Err)
		fmt.Printf("")
		os.Exit(2)
	}
	if fileInfo.IsDir() {
		fmt.Printf("Error, pathspec: \"%s\" is a directory\n", filePath)
		os.Exit(3)
	}
	file, err := os.Open(filePath)
	if err != nil {
		pathErr := err.(*os.PathError)
		fmt.Printf("Error opening file: \"%s\"\n\"%s\"\n", pathErr.Path, pathErr.Err)
		os.Exit(4)
	}
	return file, nil
}

func main() {
	var (
		fileName string
		meta     bool
		metaOn   = false
		output   = &bytes.Buffer{}
	)
	flag.StringVar(&fileName, "file", "", "specify the relative path to the file to wrap")
	flag.StringVar(&fileName, "f", "", "specify the relative path to the file to wrap")
	flag.BoolVar(&meta, "meta", true, "Specify whether to use meta blocks to denote templating")
	flag.BoolVar(&meta, "m", true, "Specify whether to use meta blocks to denote templating")
	flag.Parse()
	if fileName == "" {
		fmt.Printf("You must specify a file to target\n")
		os.Exit(1)
	}
	logOutput := os.Stdout
	file, _ := OpenFile(fileName)
	scanner := bufio.NewScanner(file)
	lineNumber := 0
	if meta {
		fmt.Printf("Meta block tagging enabled\n")
		for scanner.Scan() {
			if IsMeta(scanner.Bytes()) {
				metaOn = !metaOn
				continue
			}
			if metaOn {
				WrapLine(lineNumber, scanner.Bytes(), output, logOutput)
			}
			lineNumber += 1
		}
	} else {
		for scanner.Scan() {
			WrapIndentedLine(scanner.Bytes(), output)
		}
	}
}
