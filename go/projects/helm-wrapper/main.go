package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"

	"github.com/dlclark/regexp2"
)

type Discard struct{}

func (d Discard) Write(b []byte) (int, error) {
	return len(b), nil
}

func logger(logs chan []byte, output io.Writer) {

}

func WrapLine(lineNumber int, line []byte, output io.Writer, logs io.Writer) {
	wrappingCheckRegex := regexp.MustCompile(`^.*{{\s*"{{(hub)?-?"\s*}}.*`)
	fmt.Fprintf(logs, "##########		Line: %d		##########\n", lineNumber)
	if !wrappingCheckRegex.Match(line) {
		regex := regexp2.MustCompile(`^(\s*)({{)(hub)?(-)?(.*?)(-)?(hub)?(}})(\s*$)`, 0)
		newLine, _ := regex.Replace(string(line), `${1}{{ "${2}${3}${4}" }}${5}{{ "${6}${7}${8}" }}${9}`, 0, -1)
		hubRegex := regexp2.MustCompile(`(.*?)(\s+)(?<!{{\s*"\s*)({{hub-?)(?!\s*"\s*}})(\s+)(.*?)(\s+)(?<!{{\s*"\s*)(-?hub}})(?!\s*"\s*}})(\s+)(.*?)`, 0)
		newLine, _ = hubRegex.Replace(newLine, `${1} {{ "${3}" }} ${5} {{ "${7}" }} ${9}`, 0, -1)
		fmt.Fprintf(output, "%s\n", newLine)
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

func OpenFile(filePath string, fOptions int) (ffile *os.File, exitCode int, ferr error) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		if fOptions&os.O_CREATE == 0 {
			pathErr := err.(*os.PathError)
			returnErr := fmt.Errorf("error opening file: %s\n\"%s\"\n", pathErr.Path, pathErr.Err)
			return &os.File{}, 2, returnErr
		} else {
			ffile, err := os.OpenFile(filePath, fOptions, 0644)
			if err != nil {
				pathErr := err.(*os.PathError)
				returnErr := fmt.Errorf("error opening file: %s\n\"%s\"\n", pathErr.Path, pathErr.Err)
				return &os.File{}, 2, returnErr
			}
			return ffile, 0, nil
		}
	}
	if fileInfo.IsDir() {
		returnErr := fmt.Errorf("Error, pathspec: \"%s\" is a directory\n", filePath)
		return &os.File{}, 3, returnErr
	}
	file, err := os.OpenFile(filePath, fOptions, 0644)
	if err != nil {
		pathErr := err.(*os.PathError)
		returnErr := fmt.Errorf("error opening file: %s\n\"%s\"\n", pathErr.Path, pathErr.Err)
		fmt.Printf("Error opening file: \"%s\"\n\"%s\"\n", pathErr.Path, pathErr.Err)
		return &os.File{}, 4, returnErr
	}
	return file, 0, nil
}

func OpenOutputFile(filePath string, fOptions int) (*os.File, error) {
	file, _, err := OpenFile(filePath, fOptions)
	if err != nil {
		return nil, fmt.Errorf("error opening output file: %s\n", err)
	}
	return file, nil

}

func main() {
	var (
		inputFileName  string
		outputFileName string
		meta           bool
		verbose        bool
		metaOn         = false
		output         = &bytes.Buffer{}
		logOutput      io.Writer
	)
	flag.StringVar(&inputFileName, "file", "", "specify the relative path to the file to wrap")
	flag.StringVar(&inputFileName, "f", "", "specify the relative path to the file to wrap")
	flag.StringVar(&outputFileName, "o", "", "specify the relative path to the file to wrap")
	flag.StringVar(&outputFileName, "output", "", "specify the relative path to the file to wrap")
	flag.BoolVar(&meta, "meta", true, "Specify whether to use meta blocks to denote templating")
	flag.BoolVar(&meta, "m", true, "Specify whether to use meta blocks to denote templating")
	flag.BoolVar(&verbose, "verbose", true, "Specifies whether to use verbose output")
	flag.BoolVar(&verbose, "v", true, "Specifies whether to use verbose output")
	flag.Parse()
	if inputFileName == "" {
		fmt.Printf("You must specify a file to target\n")
		os.Exit(1)
	}
	if verbose {
		logOutput = os.Stdout
	} else {
		logOutput = Discard{}
	}
	inputFile, ec, err := OpenFile(inputFileName, os.O_RDONLY)
	if err != nil {
		fmt.Printf("ERROR: %s\n", err)
		os.Exit(ec)
	}
	outputFile, ec, err := OpenFile(outputFileName, os.O_WRONLY|os.O_TRUNC|os.O_CREATE)
	if err != nil {
		fmt.Printf("ERROR: %s\n", err)
		os.Exit(ec)
	}
	scanner := bufio.NewScanner(inputFile)
	lineNumber := 0
	if meta {
		fmt.Printf("Meta block tagging enabled\n")
		for scanner.Scan() {
			line := scanner.Bytes()
			if IsMeta(line) {
				metaOn = !metaOn
				continue
			}
			if metaOn {
				WrapLine(lineNumber, scanner.Bytes(), outputFile, logOutput)
			} else {
				outputFile.Write(line)
				outputFile.Write([]byte("\n"))
			}
			lineNumber++
		}
	} else {
		for scanner.Scan() {
			WrapIndentedLine(scanner.Bytes(), output)
		}
	}
}
