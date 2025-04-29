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
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/dlclark/regexp2"
)

type Env struct {
	Config     *FlagConfig
	ErrO       io.Writer
	LogO       io.Writer
	InputFile  io.Reader
	OutputFile io.Writer
}

type Logger struct {
	Info  io.Writer
	Error io.Writer
	Warn  io.Writer
}

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
	LogFile    string
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

func WrapLine(lineNumber int, line []byte, output io.Writer, logs io.Writer, nIndent int, indent int, env *Env) {
	if env.Config.Verbose {
		fmt.Fprintf(env.LogO, "Executing function: \"WrapLine\"\n")
	}
	wrappingCheckRegex := regexp.MustCompile(`^.*{{\s*"{{(hub)?-?"\s*}}.*`)
	fmt.Fprintf(env.LogO, "##########		Line: %d		##########\n", lineNumber)
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
		fmt.Fprintf(logs, "%s\n", newLine)
		fmt.Fprintf(env.LogO, "%s    ->     %s\n", line, newLine)
		return
	} else {
		fmt.Fprintf(output, "%s%s\n", indentation, line)
		fmt.Fprintf(logs, "%s%s\n", indentation, line)
	}
	fmt.Fprintf(env.LogO, "##########		UNCHANGED		##########\n")
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

func WrapIndentedLine(line []byte, output io.Writer) {
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

func OpenFile(filePath string, fOptions int, env *Env) (ffile *os.File, exitCode int, ferr error) {
	if env.Config.Verbose {
		fmt.Fprintf(env.LogO, "Filepath: %s\n", filePath)
	}
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		if fOptions&os.O_CREATE == 0 {
			pathErr := err.(*os.PathError)
			returnErr := fmt.Errorf("error opening file: %s\n\"%s\"\n", pathErr.Path, pathErr.Err)
			return &os.File{}, 2, returnErr
		} else {
			ffile, err := os.OpenFile(filePath, fOptions, 0o644)
			if err != nil {
				pathErr := err.(*os.PathError)
				returnErr := fmt.Errorf("error opening file:\nPath: %s\nError: \"%s\"\n", pathErr.Path, pathErr.Err)
				return &os.File{}, 2, returnErr
			}
			return ffile, 0, nil
		}
	} else if fileInfo.IsDir() {
		returnErr := fmt.Errorf("Error, pathspec: \"%s\" is a directory\n", filePath)
		return &os.File{}, 3, returnErr
	} else {
		ffile, err = os.OpenFile(filePath, fOptions, 0o644)
		if err != nil {
			pathErr := err.(*os.PathError)
			returnErr := fmt.Errorf("error opening file: %s\n\"%s\"\n", pathErr.Path, pathErr.Err)
			fmt.Printf("Error opening file: Path: \"%s\"\n Error: \"%s\"\n", pathErr.Path, pathErr.Err)
			return &os.File{}, 4, returnErr
		}
	}
	return ffile, 0, nil
}

func OpenOutputFile(filePath string, fOptions int, env *Env) (*os.File, error) {
	file, _, err := OpenFile(filePath, fOptions, env)
	if err != nil {
		return nil, fmt.Errorf("error opening output file: %s\n", err)
	}
	return file, nil
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
	outputFileName = regex.ReplaceAllString(inputFileName, "${1}.yaml")
	return outputFileName, nil
}

func ProcessTemplate(input io.Reader, output io.Writer, env *Env) error {
	templateOutput := &bytes.Buffer{}
	if env.Config.Verbose {
		fmt.Fprintf(env.LogO, "Executing function \"ProcessTemplate\"\n")
	}

	var (
		meta       = false
		metaOn     = false
		lineNumber = 0
		nIndent    = 0
		tIndent    = 0
		indent     = env.Config.Indent
	)
	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		if env.Config.Verbose {
			fmt.Fprintf(env.LogO, "Scanning line number: %d\n", lineNumber)
		}
		line := scanner.Bytes()
		if meta, tIndent = IsMeta(line); meta {
			metaOn = !metaOn
			if env.Config.Verbose {
				fmt.Fprintf(env.LogO, "Meta tag block found, meta set to %t\n", metaOn)
			}
			if metaOn {
				nIndent = tIndent
			}
			lineNumber++
			continue
		}
		if metaOn {
			line := scanner.Bytes()
			if env.Config.Verbose {
				fmt.Fprintf(env.LogO, "Wrapping line: %s\n", line)
			}
			WrapLine(lineNumber, line, output, templateOutput, nIndent, indent, env)
		} else {
			output.Write(line)
			output.Write([]byte("\n"))
			if env.Config.Verbose {
				fmt.Fprint(templateOutput, string(line))
				fmt.Fprint(templateOutput, "\n")
			}
		}
		lineNumber++
	}
	if env.Config.Verbose {
		fmt.Fprintf(os.Stdout, "Contents of buffer for operations on file\n")
		templateOutput.WriteTo(os.Stdout)
	}
	err := scanner.Err()
	if err != nil {
		fmt.Fprintf(env.ErrO, "Error when reading file: %v \n", err)
	}
	return err
}

func TraverseDirectory(dirPath string, env *Env) error {
	absPath, _ := filepath.Abs(dirPath)
	fmt.Fprintf(env.ErrO, "Scanning directory: %s\n", absPath)
	root, ok := os.DirFS(dirPath).(WrapDirFS)
	if !ok {
		fmt.Fprintf(env.ErrO, "Error, could not type assert %T to WrapDirFS\n", root)
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
			ProcessTemplateDir(dirPath, root, dirEntries, env)
		} else {
			ProcessDir(dirPath, dirEntries, env, env.ErrO)
		}
	} else {
		fmt.Fprintf(env.ErrO, "%s: not a directory\n", absPath)
		os.Exit(1001)
	}

	return nil
}

func ProcessTemplateFile(fp string, env *Env) (errs []error) {
	if env.Config.Verbose {
		fmt.Fprintf(env.LogO, "Executing function: \"ProcessTemplateFile\"\n")
	}
	inputFileName := filepath.Base(fp)
	dirpath := filepath.Dir(fp)
	inputFile, _, err := OpenFile(fp, os.O_RDONLY, env)
	if err != nil {
		fmt.Fprintf(env.ErrO, "Error opening template file: %s\n", fp)
		errs := append(errs, err)
		return errs
	}
	defer inputFile.Close()
	outputFileName, err := GetOutputName(inputFileName)
	if err != nil {
		fmt.Fprintf(env.LogO, "Non-template file: \"%s\" encountered, skipping...\n", fp)
		return errs
	}
	outputFilePath := filepath.Join(dirpath, outputFileName)
	if env.Config.Verbose {
		fmt.Fprintf(env.LogO, "Using %s as output file\n", outputFilePath)
	}
	outputFile, _, err := OpenFile(outputFilePath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, env)
	if env.Config.Verbose {
		fmt.Fprintf(env.LogO, "Output file is: %s\n", outputFile.Name())
	}
	err = ProcessTemplate(inputFile, outputFile, env)
	if err != nil {
		errs = append(errs, err)
	}
	defer outputFile.Close()
	if env.Config.Verbose {
		fmt.Fprintf(env.ErrO, "Closing output file: %s\n", outputFilePath)
	}
	return errs
}

func ProcessDir(dirPath string, dirEntries []fs.DirEntry, env *Env, errW io.Writer) {
	absPath, err := filepath.Abs(dirPath)
	if err != nil {
		fmt.Fprintf(env.ErrO, "Error opening directory: \"%s\"\n%s\n", absPath, err)
		return
	}
	for _, entry := range dirEntries {
		if entry.IsDir() {
			TraverseDirectory(filepath.Join(dirPath, entry.Name()), env)
		}
	}
}

func ProcessTemplateDir(dirPath string, root WrapDirFS, dirEntries []fs.DirEntry, env *Env) {
	absPath, err := filepath.Abs(dirPath)
	if err != nil {
		fmt.Fprintf(env.ErrO, "%s\n", dirPath)
		fmt.Fprintf(env.ErrO, "Error opening directory: \"%s\"\n%s\n", absPath, err)
		return
	}
	fmt.Fprintf(env.LogO, "Template director %s identified\n", absPath)
	nameRegex := regexp.MustCompile(`^\.[A-Za-z0-9\-]+\.template$`)
	for _, entry := range dirEntries {
		entryName := entry.Name()
		if entry.IsDir() {
			TraverseDirectory(filepath.Join(dirPath, entryName), env)
		} else if nameRegex.MatchString(entry.Name()) {
			errs := ProcessTemplateFile(filepath.Join(dirPath, entryName), env)
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
	flag.StringVar(&config.LogFile, "l", "", "Sepcifies the destination output file for logs")
	flag.StringVar(&config.LogFile, "log", "", "Sepcifies the destination output file for logs")
	flag.StringVar(&config.LogFile, "logs", "", "Sepcifies the destination output file for logs")
	flag.Parse()
	return config
}

func ValidateFlags(config *FlagConfig) (errs []error) {
	if config.Recurse {
		if config.Directory == "" {
			errs = append(errs, fmt.Errorf("Error: You must specify a directory when using the recurse option"))
		} else {
			if config.InputFile != "" {
				errs = append(errs, fmt.Errorf("Warning: specifying an input file with the recurse option will cause the flag value to be ignored"))
			}
			if config.OutputFile != "" {
				errs = append(errs, fmt.Errorf("Warning: specifying an output file with the recurse option will cause the flag value to be ignored"))
			}
		}
	} else {
		if config.InputFile == "" {
			errorf := "Error: you must specify an input file when the recurse option is not specified\n"
			SHandleErrors(os.Stderr, errorf, nil, 1)
		}
	}
	return errs
}

func HandleErrorsWithTrace(callingFunction string, env *Env, errs ...error) {
	fmt.Fprintf(env.ErrO, "Error(s) from function: %s\n", callingFunction)
	HandleErrorsWithEnv(env, errs)
}

func HandleErrorsWithEnv(env *Env, errs []error) {
	for idx, err := range errs {
		fmt.Fprintf(env.ErrO, "Err NO #%d: %s\n", idx, err.Error())
		os.Exit(1)
	}
}

func SHandleErrors(errO io.Writer, errorF string, err error, ec int) {
	if err != nil {
		fmt.Fprintf(errO, errorF, err)
	} else {
		fmt.Fprint(errO, errorF)
	}
	os.Exit(ec)
}

func HandleErrors(errs ...error) {
	for idx, err := range errs {
		fmt.Fprintf(os.Stdout, "Err NO #%d: %s\n", idx, err.Error())
		os.Exit(1)
	}
}

func SetAndValidateFlags() (config *FlagConfig) {
	config = SetFlags()
	errs := ValidateFlags(config)
	if errs != nil {
		HandleErrors(errs...)
	}
	return config
}

func NewDefaultEnv() (env *Env) {
	env = &Env{}
	env.ErrO = os.Stderr
	env.LogO = io.Discard
	return env
}

func SetEnvAndFlags() (env *Env) {
	env = NewDefaultEnv()
	config := SetAndValidateFlags()
	env.Config = config
	if env.Config.Verbose {
		DumpFlags(config)
	}
	env.ErrO = os.Stderr
	if config.Verbose {
		if config.LogFile != "" {
			logfile, _, err := OpenFile(config.LogFile, os.O_CREATE|os.O_TRUNC, env)
			if err != nil {
				HandleErrorsWithTrace("SetEnvAndFlags", env, err)
			}
			env.LogO = logfile
		} else {
			env.LogO = os.Stdout
		}
	} else {
		env.LogO = io.Discard
	}
	return env
}

func DumpFlags(config *FlagConfig) {
	values := reflect.ValueOf(config).Elem()
	types := reflect.TypeOf(config).Elem()
	//	fmt.Fprintf(os.Stdout, "Value: %+v\n", values)
	//	fmt.Fprintf(os.Stdout, "Type: %+v\n", types)
	fmt.Fprintf(os.Stdout, "\n/////////// FLAGS ///////////\n")
	for i := 0; i < values.NumField(); i++ {
		fmt.Fprintf(os.Stdout, "%v: %v\n", types.Field(i).Name, values.Field(i).Interface())
	}
	fmt.Fprintf(os.Stdout, "/////////////////////////////\n\n")
}

func main() {
	env := SetEnvAndFlags()
	//	if config.InputFile == "" {
	//		fmt.Printf("You must specify a file to target\n")
	//		os.Exit(1)
	//	}
	//	if env.Config.InputFile != "" {
	//		fmt.Printf("Input file: %s\n", config.InputFile)
	//		ProcessTemplateFile(env.Config.InputFile, env)
	//	}
	//	inputFile, ec, err := OpenFile(config.InputFile, os.O_RDONLY)
	//	if err != nil {
	//		fmt.Printf("Error opening input file:\n")
	//		fmt.Printf("ERROR: %s\n", err)
	//		os.Exit(ec)
	//	}
	//
	//	var outputFile io.Writer
	//	if config.OutputFile != "" {
	//		outputFile, ec, err = OpenFile(config.OutputFile, os.O_WRONLY|os.O_TRUNC|os.O_CREATE)
	//		if err != nil {
	//			fmt.Printf("Error opening output file:\n")
	//			fmt.Printf("ERROR: %s\n", err)
	//			os.Exit(ec)
	//		}
	//	} else {
	//		outputFile = os.Stdout
	//	}
	if env.Config.Recurse {
		TraverseDirectory(env.Config.Directory, env)
	} else {
		ProcessTemplateFile(env.Config.InputFile, env)
	}
	// scanner := bufio.NewScanner(inputFile)
	// lineNumber := 0
	// nIndent := 0
	// tIndent := 0
	//
	//	if config.Directory == "" || config.InputFile == "" {
	//		fmt.Fprintf(os.Stderr, "Error: you must provide a file or directory to operate on")
	//		os.Exit(1005)
	//
	// }
	// meta := config.Meta
	// metaOn := false
	//
	//	if meta {
	//		fmt.Printf("Meta block tagging enabled\n")
	//		for scanner.Scan() {
	//			line := scanner.Bytes()
	//			if meta, tIndent = IsMeta(line); meta {
	//				nIndent = tIndent
	//				metaOn = !metaOn
	//				if metaOn {
	//					nIndent = tIndent
	//				}
	//				continue
	//			}
	//			if metaOn {
	//				WrapLine(lineNumber, scanner.Bytes(), outputFile, logOutput, nIndent, config.Indent)
	//			} else {
	//				outputFile.Write(line)
	//				outputFile.Write([]byte("\n"))
	//			}
	//			lineNumber++
	//		}
	//	} else {
	//
	//		for scanner.Scan() {
	//			WrapIndentedLine(scanner.Bytes(), outputFile)
	//		}
	//	}
}
