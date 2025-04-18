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
	InputFile  string
	OutputFile string
	Meta       bool
	Verbose    bool
	Output     io.Writer
	LogOutput  io.Writer
	Indent     int
	Directory  string
	Recurse    bool
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

func OpenDirectory(dirPath string) ([]*os.File, int, error) {

	dirInfo, err := os.Stat(dirPath)
	if err != nil {
		if fOptions&os.O_CREATE == 0 {
			pathErr := err.(*os.PathError)
			returnErr := fmt.Errorf("error opening dir: %s\n\"%s\"\n", pathErr.Path, pathErr.Err)
			return &os.File{}, 2, returnErr
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

func ProcessFile(fsys WrapDirFS, fileName string, dirPath string, env *Env) (errs []error) {
	outputFileName, err := GetOutputName(fileName)
	if err != nil {
		return append(errs, err)
	}
	inFilePath := filepath.Join(dirPath, fileName)
	outFilePath := filepath.Join(dirPath, outputFileName)
	fmt.Fprintf(env.LogO, "Attempting to open %s for reading...\n", inFilePath)
	inputFile, err := fsys.Open(fileName)
	if err != nil {
		fmt.Fprintf(env.LogO, "Error opening %s for reading: \"%s\"\n", inFilePath, err)
		return append(errs, err)
	}
	outputFile, err := OpenOutputFile(outFilePath, os.O_APPEND|os.O_TRUNC|os.O_CREATE)
	fmt.Fprintf(env.LogO, "Attempting to open %s for writing\n", outFilePath)
	if err != nil {
		fmt.Fprintf(env.LogO, "Error opening %s for reading: \"%s\"\n", outFilePath, err)
		return append(errs, err)
	}
	defer outputFile.Close()
	defer inputFile.Close()

	var (
		meta       = false
		metaOn     = false
		lineNumber = 0
		nIndent    = 0
		tIndent    = 0
		indent     = env.Config.Indent
		logOutput  io.Writer
	)

	if err != nil {
		errs = append(errs, fmt.Errorf("%s\n", err))
		os.Exit(1003)
	}
	scanner := bufio.NewScanner(inputFile)
	fmt.Fprintf(env.LogO, "Meta block tagging enabled\n")
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

func TraverseDirectory(dirPath string, env *Env, errW io.Writer) error {
	absPath, _ := filepath.Abs(dirPath)
	fmt.Fprintf(env.ErrO, "Scanning directory: %s\n", absPath)
	root, ok := os.DirFS(dirPath).(WrapDirFS)
	if !ok {
		fmt.Fprintf(env.ErrO, "Error , could not type assert %T to WrapDirFS\n", root)

	}
	stat, err := root.Stat(".")
	if err != nil {
		fmt.Fprintf(env.ErrO, "%s\n", absPath)
		fmt.Fprintf(env.ErrO, "%s\n", err)
		os.Exit(1000)
	}
	isDir := stat.IsDir()
	dirEntries, err := root.ReadDir(".")
	if isDir && err == nil {
		if filepath.Base(dirPath) == "templates" {
			ProcessTemplateDir(dirPath, root, dirEntries, env, env.ErrO)
		} else {
			ProcessDir(dirPath, dirEntries, env, env.ErrO)
		}
	} else {
		fmt.Fprintf(env.ErrO, "%s: not a directory\n", absPath)
		os.Exit(1001)
	}

	return nil
}

func ProcessDir(dirPath string, dirEntries []fs.DirEntry, env *Env, errW io.Writer) {
	absPath, err := filepath.Abs(dirPath)
	if err != nil {
		fmt.Fprintf(env.ErrO, "Error opening directory: \"%s\"\n%s\n", absPath, err)
		return
	}
	for _, entry := range dirEntries {
		if entry.IsDir() {
			TraverseDirectory(filepath.Join(dirPath, entry.Name()), env, env.ErrO)
		}
	}
	return
}

func ProcessTemplateDir(dirPath string, root WrapDirFS, dirEntries []fs.DirEntry, env *Env, errW io.Writer) {
	absPath, err := filepath.Abs(dirPath)
	if err != nil {
		fmt.Fprintf(env.ErrO, "%s\n", dirPath)
		fmt.Fprintf(env.ErrO, "Error opening directory: \"%s\"\n%s\n", absPath, err)
		return
	}
	fmt.Fprintf(env.LogO, "Template director %s identified\n", absPath)
	nameRegex := regexp.MustCompile(`^\.[A-Za-z0-9\-]+\.template$`)
	for _, entry := range dirEntries {
		if entry.IsDir() {
			TraverseDirectory(filepath.Join(dirPath, entry.Name()), env, env.ErrO)
		} else if nameRegex.MatchString(entry.Name()) {
			errs := ProcessFile(root, dirPath, entry.Name(), env)
			if errs != nil {
				fmt.Fprintf(env.ErrO, "Error(s) processing files:\n")
				for _, err := range errs {
					fmt.Fprintf(env.ErrO, "\t%s\n", err)
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
	flag.StringVar(&config.InputFile, "file", "", "specify the relative path to the file to wrap")
	flag.StringVar(&config.InputFile, "f", "", "specify the relative path to the file to wrap")
	flag.StringVar(&config.Directory, "directory", "", "speficy the directory to template")
	flag.StringVar(&config.Directory, "d", "", "speficy the directory to template")
	flag.StringVar(&config.OutputFile, "o", "", "specify the relative path to the file to wrap")
	flag.StringVar(&config.OutputFile, "output", "", "specify the relative path to the file to wrap")
	flag.IntVar(&config.Indent, "indent", 2, "Number of spaces for an indent")
	flag.IntVar(&config.Indent, "i", 2, "Number of spaces for an indent")
	flag.BoolVar(&config.Meta, "meta", true, "Specify whether to use meta blocks to denote templating")
	flag.BoolVar(&config.Meta, "m", true, "Specify whether to use meta blocks to denote templating")
	flag.BoolVar(&config.Verbose, "verbose", false, "Specifies whether to use verbose output")
	flag.BoolVar(&config.Verbose, "v", false, "Specifies whether to use verbose output")
	flag.BoolVar(&config.Recurse, "recurse", false, "Specifies whether to use verbose output")
	flag.BoolVar(&config.Recurse, "r", false, "Specifies whether to use verbose output")
	flag.Parse()
	return config

}

type Env struct {
	Config     *FlagConfig
	ErrO       io.Writer
	LogO       io.Writer
	InputFile  io.Reader
	OutputFile io.Writer
}

func ValidateFlags(config *FlagConfig) (errs []error) {
	if config.Recurse {
		if config.Directory == "" {
			errs = append(errs, fmt.Errorf("Error: You must specify a directory when using the recurse option"))
		}
		if config.InputFile != "" {
			errs = append(errs, fmt.Errorf("Warning: specifying an input file with the recurse option will cause the flag value to be ignored"))
		}
	}
	return errs
}

func HandleErrors(env *Env, errs []error) {
	for idx, err := range errs {
		env.ErrO.Write([]byte(fmt.Sprintf("Err NO #%d: %s\n", idx, err.Error())))
		os.Exit(1)
	}
}

func SetAndValidateFlags(env *Env) (config *FlagConfig) {
	config = SetFlags()
	env.Config = config
	errs := ValidateFlags(config)
	if errs != nil {
		HandleErrors(env, errs)
	}
	return config
}

func SetEnvAndFlags() (env *Env) {
	config := SetFlags()
	env = &Env{}
	env.Config = config
	env.ErrO = os.Stderr
	if config.Verbose {
		env.LogO = os.Stdout
	} else {
		env.LogO = io.Discard
	}
	return env
}

func main() {
	env := SetEnvAndFlags()
	config := env.Config
	var (
		logOutput io.Writer
	)
	if config.InputFile == "" {
		fmt.Printf("You must specify a file to target\n")
		os.Exit(1)
	}
	fmt.Printf("Input file: %s\n", config.InputFile)
	if config.Verbose {
		logOutput = os.Stdout
	} else {
		logOutput = Discard{}
	}
	inputFile, ec, err := OpenFile(config.InputFile, os.O_RDONLY)
	if err != nil {
		fmt.Printf("Error opening input file:\n")
		fmt.Printf("ERROR: %s\n", err)
		os.Exit(ec)
	}

	var outputFile io.Writer
	if config.OutputFile != "" {
		outputFile, ec, err = OpenFile(config.OutputFile, os.O_WRONLY|os.O_TRUNC|os.O_CREATE)
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
	if config.Directory == "" || config.InputFile == "" {
		fmt.Fprintf(os.Stderr, "Error: you must provide a file or directory to operate on")
		os.Exit(1005)

	}
	if config.Verbose {
		logOutput = os.Stdout
	} else {
		logOutput = io.Discard
	}
	meta := config.Meta
	metaOn := false
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
				WrapLine(lineNumber, scanner.Bytes(), outputFile, logOutput, nIndent, config.Indent)
			} else {
				outputFile.Write(line)
				outputFile.Write([]byte("\n"))
			}
			lineNumber++
		}
	} else {
		for scanner.Scan() {
			WrapIndentedLine(scanner.Bytes(), outputFile)
		}
	}
}
