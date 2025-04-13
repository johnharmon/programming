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

type FlagConfig struct {
	inputFile  string
	outputFile string
	meta       bool
	verbose    bool
	output     io.Writer
	logOutput  io.Writer
	indent     int
	directory  string
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

func ValidateInputName(inputFileName string) (matches [][]string, validationRegex *regexp.Regexp) {
	validationRegex = regexp.MustCompile(`^\.([A-Za-z0-9\-]+)\.template$`)
	matches = validationRegex.FindAllStringSubmatch(inputFileName, -1)
	return matches, validationRegex
}

func GetOutputName(inputFileName string) (outputFileName string, err error) {
	matches, regex := ValidateInputName(inputFileName)
	if matches == nil {
		return "", fmt.Errorf("Error: not a valid file name for templating")
	}
	outputFileName = regex.ReplaceAllString(inputFileName, "{$1}.yaml")
	return outputFileName, nil
}

func ProcessFile(fsys WrapDirFS, fileName string, dirPath string, config *FlagConfig) (errs []error) {
	outputFileName, err := GetOutputName(fileName)
	if err != nil {
		return append(errs, err)
	}
	outputFile, err := OpenOutputFile(filepath.Join(dirPath, outputFileName), os.O_APPEND|os.O_TRUNC|os.O_CREATE)
	if err != nil {
		return append(errs, err)
	}
	defer outputFile.Close()
	inputFile, err := fsys.Open(fileName)
	if err != nil {
		return append(errs, err)
	}
	defer inputFile.Close()

	var (
		meta       = false
		metaOn     = false
		lineNumber = 0
		nIndent    = 0
		tIndent    = 0
		indent     = config.indent
		logOutput  io.Writer
	)

	if err != nil {
		errs = append(errs, fmt.Errorf("%s\n", err))
		os.Exit(1003)
	}
	scanner := bufio.NewScanner(inputFile)
	fmt.Fprintf(logOutput, "Meta block tagging enabled\n")
	for scanner.Scan() {
		line := scanner.Bytes()
		if meta, tIndent = IsMeta(line); meta {
			metaOn = !metaOn
			if metaOn {
				nIndent = tIndent
			}
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

	return errs
}

func TraverseDirectory(dirPath string, config *FlagConfig, errW io.Writer) error {
	absPath, _ := filepath.Abs(dirPath)
	root, ok := os.DirFS(dirPath).(WrapDirFS)
	if !ok {
		fmt.Fprintf(errW, "Error , could not type assert %T to WrapDirFS\n", root)

	}
	stat, err := root.Stat(".")
	if err != nil {
		fmt.Fprintf(errW, "%s\n", absPath)
		fmt.Fprintf(errW, "%s\n", err)
		os.Exit(1000)
	}
	isDir := stat.IsDir()
	if isDir {
		if filepath.Base(dirPath) == "templates" {
			ProcessTemplateDir(dirPath, root, config, errW)
		} else {
			ProcessDir(dirPath, root, config, errW)
		}
	} else {
		fmt.Fprintf(errW, "%s: not a directory\n", absPath)
		os.Exit(1001)
	}

	return nil
}

func ProcessDir(dirPath string, root WrapDirFS, config *FlagConfig, errW io.Writer) {
	dirEntries, err := root.ReadDir(".")
	absPath, err := filepath.Abs(dirPath)
	if err != nil {
		fmt.Fprintf(errW, "Error opening directory: \"%s\"\n%s\n", absPath, err)
		return
	}
	for _, entry := range dirEntries {
		if entry.IsDir() {
			TraverseDirectory(filepath.Join(dirPath, entry.Name()), config, errW)
		}
	}
	return
}

func ProcessTemplateDir(dirPath string, root WrapDirFS, config *FlagConfig, errW io.Writer) {
	dirEntries, err := root.ReadDir(".")
	absPath, err := filepath.Abs(dirPath)
	if err != nil {
		fmt.Fprintf(errW, "%s\n", dirPath)
		fmt.Fprintf(errW, "Error opening directory: \"%s\"\n%s\n", absPath, err)
		return
	}
	nameRegex := regexp.MustCompile(`^\.[A-Za-z0-9\-]+\.template$`)
	for _, entry := range dirEntries {
		if entry.IsDir() {
			TraverseDirectory(dirPath, config, errW)
		} else if nameRegex.MatchString(entry.Name()) {
			errs := ProcessFile(root, dirPath, entry.Name(), config)
			if errs != nil {
				fmt.Fprintf(errW, "Error(s) processing files:\n")
				for _, err := range errs {
					fmt.Fprintf(errW, "\t%s\n", err)
				}
				return
			}
		} else {
			continue
		}
	}

}

func SetFlags() (config *FlagConfig) {
	config = &FlagConfig{}
	flag.StringVar(&config.inputFile, "file", "", "specify the relative path to the file to wrap")
	flag.StringVar(&config.inputFile, "f", "", "specify the relative path to the file to wrap")
	flag.StringVar(&config.directory, "directory", "", "speficy the directory to template")
	flag.StringVar(&config.directory, "d", "", "speficy the directory to template")
	flag.StringVar(&config.outputFile, "o", "", "specify the relative path to the file to wrap")
	flag.StringVar(&config.outputFile, "output", "", "specify the relative path to the file to wrap")
	flag.IntVar(&config.indent, "indent", 2, "Number of spaces for an indent")
	flag.IntVar(&config.indent, "i", 2, "Number of spaces for an indent")
	flag.BoolVar(&config.meta, "meta", true, "Specify whether to use meta blocks to denote templating")
	flag.BoolVar(&config.meta, "m", true, "Specify whether to use meta blocks to denote templating")
	flag.BoolVar(&config.verbose, "verbose", false, "Specifies whether to use verbose output")
	flag.BoolVar(&config.verbose, "v", false, "Specifies whether to use verbose output")
	flag.Parse()
	return config

}

func main() {
	var (
		logOutput io.Writer
	)
	config := SetFlags()
	if config.inputFile == "" {
		fmt.Printf("You must specify a file to target\n")
		os.Exit(1)
	}
	fmt.Printf("Input file: %s\n", config.inputFile)
	if config.verbose {
		logOutput = os.Stdout
	} else {
		logOutput = Discard{}
	}
	inputFile, ec, err := OpenFile(config.inputFile, os.O_RDONLY)
	if err != nil {
		fmt.Printf("Error opening input file:\n")
		fmt.Printf("ERROR: %s\n", err)
		os.Exit(ec)
	}

	var outputFile io.Writer
	if config.outputFile != "" {
		outputFile, ec, err = OpenFile(config.outputFile, os.O_WRONLY|os.O_TRUNC|os.O_CREATE)
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
	if config.directory == "" || config.inputFile == "" {
		fmt.Fprintf(os.Stderr, "Error: you must provide a file or directory to operate on")
		os.Exit(1005)

	}
	if config.verbose {
		logOutput = os.Stdout
	} else {
		logOutput = io.Discard
	}
	if meta {
		fmt.Printf("Meta block tagging enabled\n")
		for scanner.Scan() {
			line := scanner.Bytes()
			if meta, tIndent = IsMeta(line); meta {
				nIndent = tIndent
				metaOn = !metaOn
				if metaOn {
					nIndent = tIndent
				}
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
