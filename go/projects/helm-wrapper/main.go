package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/dlclark/regexp2"
)

type WrapDirFS interface {
	fs.FS
	fs.StatFS
	fs.ReadFileFS
	fs.ReadDirFS
}

type Discard struct{}

func (d Discard) Write(b []byte) (int, error) {
	return len(b), nil
}

func logger(logs chan []byte, output io.Writer) {

}

func WrapLine(lineNumber int, line []byte, output io.Writer, logs io.Writer, nIndent int, indent int) {
	wrappingCheckRegex := regexp.MustCompile(`^.*{{\s*"{{(hub)?-?"\s*}}.*`)
	fmt.Fprintf(logs, "##########		Line: %d		##########\n", lineNumber)
	indentation := strings.Repeat(" ", (nIndent * indent))
	if !wrappingCheckRegex.Match(line) {
		regex := regexp2.MustCompile(`^(\s*)({{)(hub)?(-)?(.*?)(-)?(hub)?(}})(\s*$)`, 0)
		var newLine string
		if ok, _ := regex.MatchString(string(line)); ok {
			newLine, _ = regex.Replace(string(line), fmt.Sprintf(`%s${1}{{ "${2}${3}${4}" }}${5}{{ "${6}${7}${8}" }}${9}`, indentation), 0, -1)
		} else {
			newLine = fmt.Sprintf("%s%s", indentation, line)
		}
		hubRegex := regexp2.MustCompile(`(.*?)(\s+)(?<!{{\s*"\s*)({{hub-?)(?!\s*"\s*}})(\s+)(.*?)(\s+)(?<!{{\s*"\s*)(-?hub}})(?!\s*"\s*}})(\s+)(.*?)`, 0)
		newLine, _ = hubRegex.Replace(newLine, `${1} {{ "${3}" }} ${5} {{ "${7}" }} ${9}`, 0, -1)
		fmt.Fprintf(output, "%s\n", newLine)
		fmt.Fprintf(logs, "%s    ->     %s\n", line, newLine)
		return
	} else {
		fmt.Fprintf(output, "%s%s\n", indentation, line)
	}
	fmt.Fprintf(logs, "##########		UNCHANGED		##########\n")
}

func WrapHubLine(line []byte, output *bytes.Buffer) {
	regex := regexp.MustCompile(`^(\s*)({{hub-?)(.*)(-?hub}})(\s*$)`)
	newLine := regex.ReplaceAll(line, []byte(`${1}{{ "${2}" }}${3}{{ "${4}" }}${5}`))
	output.Write(newLine)
}

func IsMeta(line []byte) (meta bool, nIndent int) {
	regex := regexp.MustCompile(`^\s*#\s*meta\s*([0-9]*)\s*$`)
	meta = false
	nIndent = 0
	matches := regex.FindSubmatch(line)
	var err error
	matchCount := len(matches)
	if matchCount > 0 {
		meta = true
		if len(matches) > 1 {
			if matchCount > 0 {
				nIndent, err = strconv.Atoi(string(matches[1]))
				if err != nil {
					nIndent = 0
				}
			}
		}
	}
	return meta, nIndent
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
				returnErr := fmt.Errorf("error opening file:\nPath: %s\nError: \"%s\"\n", pathErr.Path, pathErr.Err)
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
		fmt.Printf("Error opening file: Path: \"%s\"\n Error: \"%s\"\n", pathErr.Path, pathErr.Err)
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

func OpenDirectory(dirPath string) ([]*os.File, error) {

	dirInfo, err := os.Stat(dirPath)
	if err != nil {
		if fOptions&os.O_CREATE == 0 {
			pathErr := err.(*os.PathError)
			returnErr := fmt.Errorf("error opening dir: %s\n\"%s\"\n", pathErr.Path, pathErr.Err)
			return &os.dir{}, 2, returnErr
		} else {
			fdir, err := os.Opendir(dirPath, fOptions, 0644)
			if err != nil {
				pathErr := err.(*os.PathError)
				returnErr := fmt.Errorf("error opening dir:\nPath: %s\nError: \"%s\"\n", pathErr.Path, pathErr.Err)
				return &os.File{}, 2, returnErr
			}
			return fdir, 0, nil
		}
	}
	if dirInfo.IsDir() {
		returnErr := fmt.Errorf("Error, pathspec: \"%s\" is a directory\n", dirPath)
		return &os.File{}, 3, returnErr
	}
	dir, err := os.Opendir(dirPath, fOptions, 0644)
	if err != nil {
		pathErr := err.(*os.PathError)
		returnErr := fmt.Errorf("error opening dir: %s\n\"%s\"\n", pathErr.Path, pathErr.Err)
		fmt.Printf("Error opening dir: Path: \"%s\"\n Error: \"%s\"\n", pathErr.Path, pathErr.Err)
		return &os.File{}, 4, returnErr
	}
	return dir, 0, nil
	return []*os.File{}, nil
}

func ProcessDirEntry(path string, d fs.DirEntry, err error) error {
	return nil
}

func IsTemplate(parent os.DirEntry, file os.DirEntry) bool {
	return false
}

func ProcessFile(fsys WrapDirFS, fileName string, dirPath string) (errs []error) {
	file, err := fsys.Open(fileName)
	if err != nil {
		errs = append(errs, fmt.Errorf("%s\n", err))
		os.Exit(1003)
	}
	scanner := bufio.NewScanner(file)
	return errs
}

func TraverseDirectory(dirPath string, errW io.Writer) (files []*os.File, directories []*os.File) {
	absPath, _ := filepath.Abs(dirPath)
	root := os.DirFS(dirPath).(WrapDirFS)
	stat, err := root.Stat(".")
	if err != nil {
		fmt.Fprintf(errW, "%s\n", absPath)
		fmt.Fprintf(errW, "%s\n", err)
		os.Exit(1000)
	}
	isDir := stat.IsDir()
	if isDir {
		TraverseDirectory(filepath.Join(dirPath, stat.Name()), errW)
	} else {
		fmt.Fprintf(errW, "%s: not a directory\n", absPath)
		os.Exit(1001)

	}
	if filepath.Base(dirPath) == "templates" {
		dirEntries, err := root.ReadDir(".")
		if err != nil {
			fmt.Fprintf(errW, "%s\n", absPath)
			fmt.Fprintf(errW, "%s\n", err)
			os.Exit(1002)
		}
		nameRegex := regexp.MustCompile(`^\.[A-Za-z0-9\-]+\.template$`)
		for _, entry := range dirEntries {
			name := entry.Name()
			if entry.IsDir() {
				TraverseDirectory(filepath.Join(dirPath, name), errW)
			} else if nameRegex.MatchString(name) {
				errs := ProcessFile(root, name, dirPath)
				for _, err := range errs {
					fmt.Fprintf(errW, "%s\n", err)
				}
			}
			//more processing logic later
		}
	}

	return files, directories
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
		indent         int
		directory      string
	)
	flag.StringVar(&inputFileName, "file", "", "specify the relative path to the file to wrap")
	flag.StringVar(&inputFileName, "f", "", "specify the relative path to the file to wrap")
	flag.StringVar(&directory, "directory", "", "speficy the directory to template")
	flag.StringVar(&directory, "d", "", "speficy the directory to template")
	flag.StringVar(&outputFileName, "o", "", "specify the relative path to the file to wrap")
	flag.StringVar(&outputFileName, "output", "", "specify the relative path to the file to wrap")
	flag.IntVar(&indent, "indent", 2, "Number of spaces for an indent")
	flag.IntVar(&indent, "i", 2, "Number of spaces for an indent")
	flag.BoolVar(&meta, "meta", true, "Specify whether to use meta blocks to denote templating")
	flag.BoolVar(&meta, "m", true, "Specify whether to use meta blocks to denote templating")
	flag.BoolVar(&verbose, "verbose", false, "Specifies whether to use verbose output")
	flag.BoolVar(&verbose, "v", false, "Specifies whether to use verbose output")
	flag.Parse()
	if inputFileName == "" {
		fmt.Printf("You must specify a file to target\n")
		os.Exit(1)
	}
	fmt.Printf("Input file: %s\n", inputFileName)
	if verbose {
		logOutput = os.Stdout
	} else {
		logOutput = Discard{}
	}
	inputFile, ec, err := OpenFile(inputFileName, os.O_RDONLY)
	if err != nil {
		fmt.Printf("Error opening input file:\n")
		fmt.Printf("ERROR: %s\n", err)
		os.Exit(ec)
	}

	var outputFile io.Writer
	if outputFileName != "" {
		outputFile, ec, err = OpenFile(outputFileName, os.O_WRONLY|os.O_TRUNC|os.O_CREATE)
		if err != nil {
			fmt.Printf("Error opening output file:\n")
			fmt.Printf("ERROR: %s\n", err)
			os.Exit(ec)
		}
	} else {
		outputFile = os.Stdout
	}
	scanner := bufio.NewScanner(inputFile)
	lineNumber := 0
	var nIndent = 0
	var tIndent = 0
	if meta {
		fmt.Printf("Meta block tagging enabled\n")
		for scanner.Scan() {
			line := scanner.Bytes()
			if meta, tIndent = IsMeta(line); meta {
				nIndent = tIndent
				metaOn = !metaOn
				continue
			}
			if metaOn {
				WrapLine(lineNumber, scanner.Bytes(), outputFile, logOutput, nIndent, indent)
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
