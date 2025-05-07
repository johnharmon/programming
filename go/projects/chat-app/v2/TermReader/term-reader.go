package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"
	"time"

	"golang.org/x/term"
)

const (
	MR = "\033[1C"
	ML = "\033[1D"
	MU = "\033[1A"
	MD = "\033[1B"
)

type Env struct {
	Config       *FlagConfig
	DebugWriter  io.Writer
	OutputHeader string
	OutputFooter string
	OutputPrefix string
	OutputSuffix string
}

func (e Env) DWrite(b []byte) {
	fmt.Fprintf(e.DebugWriter, "%b", b)
}

func (e Env) DWriteS(s string) {
	fmt.Fprintf(e.DebugWriter, "%s", s)
}

type InputScanner struct {
	Remaining   []byte
	LastMessage []byte
	LastChunk   []byte
	Input       io.Reader
	Output      io.Writer
	Delimiter   byte
}

type FlagConfig struct {
	Debug      bool
	Verbose    bool
	Terminal   bool
	Verbosity1 bool
	Verbosity2 bool
	Verbosity3 bool
	Verbosity4 bool
	Raw        bool
	Logs       io.Writer
}

func (fc FlagConfig) IsVerbose() bool {
	if fc.Verbosity1 || fc.Verbosity2 || fc.Verbosity3 || fc.Verbosity4 {
		return true
	}
	return false
}

type FormatInfo struct {
	OutputRaw      []byte
	OutputLines    [][]byte // raw bytes representing the output
	WrappingLength int      // how long the output prefix and suffix combined is
	TermWidth      int      // how many characters wide the terminal output is
	TermLength     int      // how many lines long including headers and footers the ouput is
}

type Cell struct {
	formatInfo            *FormatInfo
	RawContent            *bytes.Buffer
	RawInput              *bytes.Buffer
	ContentReader         *bytes.Reader
	Out                   io.Writer
	In                    io.Reader
	DisplayContent        *bytes.Buffer
	CursorPosition        int
	LogicalCursorPosition int
	DisplayCursorPosition int
	CellHistory           []*Cell
}

func (oc *Cell) Display(o io.Writer, env *Env) {
	formattedContent := WrapOutput(env, oc.RawContent.Bytes())
	fmt.Fprint(o, formattedContent)
}

func (oc *Cell) OverWrite(newContent []byte) {
}

type CellHistory struct {
	RawContent     []byte
	DisplayContent []byte
	CursorPosition int
}

type ModificationSequence struct {
	Bytes       []byte
	Name        string
	Raw         string
	ForceRedraw bool
	IsMultiByte bool
}

func (es ModificationSequence) String() string {
	return es.Name
}

func MakeRawTerm() {
	fd := int(os.Stdin.Fd())
	oldState, err := term.MakeRaw(fd)
	if err != nil {
		panic(err)
	}
	defer term.Restore(fd, oldState)
	RawTermInterface()
}

func HandleNormalByte(b byte, out io.Writer) {
	fmt.Fprint(out, b)
}

func makeModificationSequenceDispatcher() (dispatcher map[string]*ModificationSequence) {
	dispatcher = map[string]*ModificationSequence{}
	dispatcher["\x1b[A"] = &ModificationSequence{
		Raw:  "\x1b[A",
		Name: "UpArrow",
	}
	dispatcher["\x1b[B"] = &ModificationSequence{
		Raw:  "\x1b[B",
		Name: "DownArrow",
	}
	dispatcher["\x1b[C"] = &ModificationSequence{
		Raw:  "\x1b[C",
		Name: "RightArrow",
	}
	dispatcher["\x1b[D"] = &ModificationSequence{
		Raw:  "\x1b[D",
		Name: "LeftArrow",
	}
	return dispatcher
}

var ModificationSequenceMap = map[string]*ModificationSequence{
	"\x1b":    {Bytes: []byte("\x1b"), Name: "Escape", Raw: "\x1b", IsMultiByte: true},
	"\x1b[A":  {Bytes: []byte("\x1b[A"), Name: "UpArrow", Raw: "\x1b[A", IsMultiByte: true, ForceRedraw: false},
	"\x1b[B":  {Bytes: []byte("\x1b[B"), Name: "DownArrow", Raw: "\x1b[B", IsMultiByte: true, ForceRedraw: false},
	"\x1b[C":  {Bytes: []byte("\x1b[C"), Name: "RightArrow", Raw: "\x1b[C", IsMultiByte: true, ForceRedraw: false},
	"\x1b[D":  {Bytes: []byte("\x1b[D"), Name: "LeftArrow", Raw: "\x1b[D", IsMultiByte: true, ForceRedraw: false},
	"\x1b[3~": {Bytes: []byte("\x1b[3~"), Name: "Delete", Raw: "\x1b[3~", IsMultiByte: false, ForceRedraw: true},
	"\x7F":    {Bytes: []byte("\x7F"), Name: "Backspace", Raw: "\x7F", IsMultiByte: false, ForceRedraw: true},
}

func isModificationByte(b byte) (bool, *ModificationSequence) {
	m, ok := ModificationSequenceMap[string(b)]
	return ok, m
}

func makeModificationSequence(sequence []byte) (es *ModificationSequence) {
	ss := string(sequence)
	es = &ModificationSequence{
		Bytes: sequence,
		Name:  "Placeholder",
		Raw:   ss,
	}
	return es
}

func ReadModificationSequence(input io.Reader, timeout time.Duration, esc *ModificationSequence) (*ModificationSequence, error) {
	if esc.IsMultiByte {
		b := make([]byte, 32)
		deadline := time.Now().Add(timeout)
		if f, ok := input.(*os.File); ok {
			f.SetReadDeadline(deadline)
		}
		defer ClearReadDeadline(input)
		n, err := input.Read(b)
		if err != nil {
			return nil, err
		}
		esc, ok := ModificationSequenceMap[string(esc.Bytes[0])+string(b[:n])]
		if !ok {
			return nil, nil
		}
		return esc, nil
	}
	return esc, nil
}

func ClearReadDeadline(input io.Reader) {
	if f, ok := input.(*os.File); ok {
		_ = f.SetReadDeadline(time.Time{})
	}
}

func matchModificationSequence(sequence []byte) {}

func HandleEscabeByte(b byte, out io.Writer) {
}

func RawTermInterface() {
	cell := &Cell{}
	cell.Out = os.Stdout
	cell.In = os.Stdin
	cell.DisplayContent = &bytes.Buffer{}
	cell.DisplayCursorPosition = 0
	cell.RawInput = &bytes.Buffer{}
	var (
		isMod  bool
		modSeq *ModificationSequence
	)
	// typedBytes := cell.RawInput
	buf := make([]byte, 1)
	for {
		nb, err := cell.In.Read(buf)
		if err != nil {
			panic(err)
		}
		if nb > 0 {
			b := buf[0]
			isMod, modSeq = isModificationByte(b)
			if isMod {
				esc, _ := ReadModificationSequence(cell.In, (time.Millisecond * 25), modSeq)
				if esc != nil {
					// HandleModSequence(cell, esc)
					if esc.ForceRedraw {
						newLine, _ := RedrawLine(esc, cell)
						fmt.Fprintf(cell.Out, "\r\x1b[2K%s", newLine)
						DisplayDebugInfo(cell, "EscapeSequenceConditional, ForceRedraw", []string{})
						cell.RawInput.Write(esc.Bytes)
					} else {
						cell.RawInput.Write(esc.Bytes)
						fmt.Fprintf(cell.Out, "%s", esc.Bytes)
						switch esc.Name {
						case "LeftArrow":
							if cell.DisplayCursorPosition > 0 {
								cell.DisplayCursorPosition--
							}
							cell.LogicalCursorPosition += len(esc.Bytes)
							fmt.Fprintf(cell.Out, "%s", esc.Bytes)
						case "RightArrow":
							if cell.DisplayCursorPosition < len(cell.DisplayContent.Bytes()) {
								cell.DisplayCursorPosition++
							}
							cell.LogicalCursorPosition += len(esc.Bytes)
							fmt.Fprintf(cell.Out, "%s", esc.Bytes)
						}
						DisplayDebugInfo(cell, "EscapeSequenceConditional, NoRedraw", []string{})
						cell.RawInput.Write(esc.Bytes)
					}
				}
			} else {
				bslice := buf[0:1]
				cell.RawInput.WriteByte(b)
				if b == 13 {
					fmt.Fprintf(cell.Out, "\n\rYou typed: %s\r\n", cell.RawInput.Bytes())
					cell.RawInput.Reset()
				} else {
					if b == 3 {
						break
					} else {
						cell.WriteDisplayBytes(bslice)
						// cell.DisplayCursorPosition++
						// cell.LogicalCursorPosition++
					}
				}
			}
		}
	}
}

func (cell *Cell) Redraw() {
	fmt.Fprintf(cell.Out, "\r\x1b[2K\r%s", cell.DisplayContent.Bytes())
}

func (cell *Cell) WriteDisplayBytes(b []byte) {
	extra := []string{}
	blen := len(b)
	if cell.DisplayCursorPosition == 0 {
		content := append([]byte{}, cell.DisplayContent.Bytes()...)
		cell.DisplayContent.Reset()
		cell.DisplayContent.Write(b)
		extra = append(extra, string(b))
		cell.DisplayContent.Write(content)
		extra = append(extra, string(content))
		cell.DisplayCursorPosition += blen
		cell.LogicalCursorPosition += blen
		cell.Redraw()
		DisplayDebugInfo(cell, "DisplayCursorPosition = 0", extra)
	} else if cell.DisplayCursorPosition == len(cell.DisplayContent.Bytes()) {
		cell.DisplayContent.Write(b)
		cell.DisplayCursorPosition += blen
		cell.LogicalCursorPosition += blen
		fmt.Fprintf(cell.Out, "%s", b)
		DisplayDebugInfo(cell, "DisplayCursorPosition = end of display content", extra)
	} else {
		temp := append([]byte{}, cell.DisplayContent.Bytes()...)
		cell.DisplayContent.Reset()
		before := temp[0:cell.DisplayCursorPosition]
		after := temp[cell.DisplayCursorPosition:]
		cell.DisplayContent.Write(before)
		extra = append(extra, string(before))
		cell.DisplayContent.Write(b)
		extra = append(extra, string(b))
		cell.DisplayContent.Write(after)
		extra = append(extra, string(after))
		cell.DisplayCursorPosition += blen
		cell.LogicalCursorPosition += blen
		cell.Redraw()
		DisplayDebugInfo(cell, "DisplayCursorPostition = inside display content", extra)
	}
}

func DisplayDebugInfo(cell *Cell, callingInfo string, extras []string) {
	var cursorRight string
	fmt.Fprintf(cell.Out, "\x1b[B\r\x1b[2K")
	fmt.Fprintf(cell.Out, "DisplayCursorPosition: %d | LogicalCursorPosition: %d | CalledBy: %s", cell.DisplayCursorPosition, cell.LogicalCursorPosition, callingInfo)
	cursorUp := "\x1b[A\r"
	cursorRight = fmt.Sprintf("\x1b[%dC", cell.DisplayCursorPosition)
	if cell.DisplayCursorPosition > 0 {
		cursorRight = fmt.Sprintf("\x1b[%dC", cell.DisplayCursorPosition)
	}
	cursorUp = fmt.Sprintf("\x1b[%dA\r", len(extras)+1)
	//	if cell.DisplayCursorPosition > 0 {
	//		cursorRight = fmt.Sprintf("\x1b[%dC", cell.DisplayCursorPosition)
	//	}
	if len(extras) > 0 {
		for _, extra := range extras {
			fmt.Fprintf(cell.Out, "\r\n%s\r", extra)
		}
	}
	fmt.Fprintf(cell.Out, "%s%s", cursorUp, cursorRight)
}

func HandleModSequence(cell *Cell, esc *ModificationSequence) {
}

func RedrawLine(mod *ModificationSequence, cell *Cell) (newLine []byte, err error) {
	if mod.Name == "Backspace" || mod.Name == "Delete" {
		if cell.DisplayCursorPosition > 0 {
			cell.DisplayCursorPosition--
			cell.LogicalCursorPosition--
			newLine = append(cell.DisplayContent.Bytes()[:cell.DisplayCursorPosition], cell.DisplayContent.Bytes()[cell.DisplayCursorPosition+1:]...)
		}
	} else {
		newLine = cell.DisplayContent.Bytes()
	}
	cell.DisplayContent.Reset()
	cell.DisplayContent.Write(newLine)
	return newLine, nil
}

func (w *FormatInfo) Debug(env *Env) {
	env.DWriteS("Entered function \"WrapOutput\"\n")
	env.DWriteS(fmt.Sprintf("OutputRaw: %s\n", string(w.OutputRaw)))
	env.DWriteS(fmt.Sprintf("termLength: %d\n", w.TermLength))
	env.DWriteS(fmt.Sprintf("termWidth: %d\n", w.TermWidth))
	env.DWriteS(fmt.Sprintf("wrappingLength: %d\n", w.WrappingLength))
	for idx, ln := range w.OutputLines {
		env.DWriteS(fmt.Sprintf("%d: %s\n", idx, ln))
	}
}

type PrintFunc func(*Env, []byte, io.Writer)

func WrapOutput(env *Env, output []byte) (wrappedOutput string) {
	wrapInfo := GetOutputDimensions(env, output)
	wrappedOutputLines := []string{}
	wrapInfo.Debug(env)
	for i := 0; i < wrapInfo.TermLength; i++ {
		wrappedOutputLines = append(wrappedOutputLines, wrapOutputLine(env, wrapInfo, i))
	}
	return strings.Join(wrappedOutputLines, "\n") + "\n"
}

func wrapOutputLine(env *Env, wrapInfo *FormatInfo, lineNumber int) (newLine string) {
	env.DWriteS(fmt.Sprintf("Processing line %d for terminal output...\n", lineNumber))
	switch lineNumber {
	case 0:
		newLine = fmt.Sprintf("%s", strings.Repeat(env.OutputHeader, wrapInfo.TermWidth))
	case wrapInfo.TermLength - 1:
		newLine = fmt.Sprintf("%s", strings.Repeat(env.OutputFooter, wrapInfo.TermWidth))
	default:
		padding := strings.Repeat(" ", wrapInfo.TermWidth-len(wrapInfo.OutputLines[lineNumber-1])-wrapInfo.WrappingLength)
		newLine = fmt.Sprintf("%s%s%s%s", env.OutputPrefix, wrapInfo.OutputLines[lineNumber-1], padding, env.OutputSuffix)
	}
	env.DWriteS(fmt.Sprintf("Processed line #%d: %s\n", lineNumber, newLine))
	return newLine
}

func GetOutputDimensions(env *Env, output []byte) (wrapInfo *FormatInfo) {
	wrapInfo = &FormatInfo{}
	wrapInfo.OutputRaw = output
	wrapInfo.OutputLines = ExpandBytesLinewise(env, output)
	wrapInfo.WrappingLength = (len(env.OutputPrefix) + len(env.OutputSuffix))
	wrapInfo.TermWidth = LongestByteSlice(wrapInfo.OutputLines) + wrapInfo.WrappingLength
	wrapInfo.TermLength = len(wrapInfo.OutputLines) + 2
	return wrapInfo
}

func wrapOutputDebugHelper(env *Env, wrapInfo *FormatInfo) {
	env.DWriteS("Entered function \"WrapOutput\"\n")
	env.DWriteS(fmt.Sprintf("OutputRaw: %s\n", string(wrapInfo.OutputRaw)))
	env.DWriteS(fmt.Sprintf("termLength: %d\n", wrapInfo.TermLength))
	env.DWriteS(fmt.Sprintf("termWidth: %d\n", wrapInfo.TermWidth))
	env.DWriteS(fmt.Sprintf("wrappingLength: %d\n", wrapInfo.WrappingLength))
	for idx, ln := range wrapInfo.OutputLines {
		env.DWriteS(fmt.Sprintf("%d: %s\n", idx, ln))
	}
}

func ExpandBytesLinewise(env *Env, iBytes []byte) (byteLines [][]byte) {
	env.DWriteS("Entered Function: \"ExpandBytesLinewise\"\n")
	for {
		splitIndex := bytes.IndexByte(iBytes, '\n')
		env.DWriteS(fmt.Sprintf("Encountered newLine at %d\n", splitIndex))
		if splitIndex == -1 {
			byteLines = append(byteLines, iBytes)
			break
		} else {
			byteLines = append(byteLines, iBytes[0:splitIndex])
			iBytes = iBytes[splitIndex+1:]
		}
	}
	return byteLines
}

func LongestByteSlice(slices [][]byte) (longest int) {
	longest = -1
	for _, s := range slices {
		if len(s) > longest {
			longest = len(s)
		}
	}
	return longest
}

func WrapOutputLines(env *Env, output []byte) (wrappedOutput []string) {
	outputLines := ExpandBytesLinewise(env, output)
	termLength := len(outputLines) + 2
	termWidth := LongestByteSlice(outputLines)
	var termLine string
	for i := 0; i < termLength; i++ {
		if i == 0 || i == termLength-1 {
			termLine = fmt.Sprintf("%s\n", strings.Repeat("-", termWidth))
			wrappedOutput = append(wrappedOutput, termLine)
		} else {
			padding := strings.Repeat(" ", termWidth-len(outputLines[i-1]))
			termLine = fmt.Sprintf("%s%s%s%s\n", "|", outputLines[i-1], padding, "|")
		}
	}
	return wrappedOutput
}

func BasicPrinter(env *Env, input []byte, output io.Writer) {
	fmt.Fprint(output, input)
}

func BasicStringPrinter(env *Env, input []byte, output io.Writer) {
	fmt.Fprint(output, string(input))
}

func NewLineStringPrinter(env *Env, input []byte, output io.Writer) {
	fmt.Fprint(output, string(input), "\n")
}

func WrapOutputPrinter(env *Env, input []byte, output io.Writer) {
	env.DWriteS("Wrap output printer called\n")
	outputString := WrapOutput(env, input)
	output.Write([]byte(outputString))
}

func (i *InputScanner) Scan(env *Env, pf PrintFunc) {
	buf := make([]byte, 1024)
	n, _ := i.Input.Read(buf)
	message, read := i.ScanInput(buf[0:n])
	if read {
		pf(env, message, i.Output)
	}
}

func (i *InputScanner) ScanInput(input []byte) (message []byte, messageRead bool) {
	remaining := []byte{}

	delimIndex := bytes.IndexByte(input, i.Delimiter)
	if delimIndex != -1 {
		message = append(i.Remaining, input[:delimIndex]...)
		remaining = input[delimIndex+1:]
		messageRead = true
	} else {
		message = append(i.Remaining, input...)
		remaining = input[:0]
		messageRead = false
	}
	i.Remaining = remaining
	if messageRead {
		i.LastMessage = message
	}
	i.LastChunk = input
	return message, messageRead
}

func ReadUntil(input io.Reader, delim byte) (message []byte, err error) {
	return message, err
}

func ScanInput(input []byte, delim byte) (message []byte, remaining []byte, isTerminated, err error) {
	delimIndex := bytes.IndexByte(input, delim)
	if delimIndex != -1 {
		message = input[:delimIndex]
		remaining = input[delimIndex+1:]
	} else {
		message = input
		remaining = input[:0]
	}
	return message, remaining, isTerminated, err
}

func NewInputScanner(input io.Reader) *InputScanner {
	is := InputScanner{}
	is.Remaining = []byte{}
	is.LastMessage = []byte{}
	is.LastChunk = []byte{}
	is.Input = input
	is.Delimiter = '\n'
	return &is
}

func SelectPrinter(env *Env) (pf PrintFunc) {
	switch {
	case env.Config.Verbosity1:
		pf = WrapOutputPrinter
	case env.Config.Raw:
		pf = BasicPrinter
	case !env.Config.Raw:
		pf = NewLineStringPrinter
	default:
		pf = BasicStringPrinter
	}
	return pf
}

func RunScanner(env *Env) {
	scanner := NewInputScanner(os.Stdin)
	if env.Config.Debug {
		scanner.Output = os.Stdout
	}
	pf := SelectPrinter(env)
	for {
		scanner.Scan(env, pf)
	}
}

func DumpFlags(config *FlagConfig) {
	values := reflect.ValueOf(config).Elem()
	types := reflect.TypeOf(config).Elem()
	fmt.Fprintf(os.Stdout, "\n/////////// FLAGS ///////////\n")
	for i := 0; i < values.NumField(); i++ {
		fmt.Fprintf(os.Stdout, "%v: %v\n", types.Field(i).Name, values.Field(i).Interface())
	}
	fmt.Fprintf(os.Stdout, "/////////////////////////////\n\n")
}

func NewEnv(config *FlagConfig) (env *Env) {
	env = &Env{}
	env.Config = config
	if env.Config.IsVerbose() {
		env.DebugWriter = os.Stdout
	}
	return env
}

func NewDefaultEnv(config *FlagConfig) (env *Env) {
	env = NewEnv(config)
	env.OutputFooter = "-"
	env.OutputHeader = "-"
	env.OutputPrefix = "| "
	env.OutputSuffix = " |"
	return env
}

func ParseFlags() (config *FlagConfig) {
	config = &FlagConfig{}
	flag.BoolVar(&config.Debug, "d", false, "use debug mode (boolean toggle)")
	flag.BoolVar(&config.Terminal, "t", false, "use terminal mode (boolean toggle)")
	flag.BoolVar(&config.Verbose, "verbose", false, "use verbose mode (boolean toggle)")
	flag.BoolVar(&config.Verbosity1, "v", false, "use verbose mode (boolean toggle)")
	flag.BoolVar(&config.Verbosity2, "vv", false, "use verbose mode (boolean toggle)")
	flag.BoolVar(&config.Verbosity3, "vvv", false, "use verbose mode (boolean toggle)")
	flag.BoolVar(&config.Verbosity4, "vvvv", false, "use verbose mode (boolean toggle)")
	flag.BoolVar(&config.Raw, "raw", false, "use raw output mode (boolean toggle)")

	flag.Parse()
	if config.Verbose {
		config.Logs = os.Stdout
		DumpFlags(config)
	} else {
		config.Logs = io.Discard
	}
	return config
}

func main() {
	config := ParseFlags()
	env := NewDefaultEnv(config)
	if config.Terminal {
		MakeRawTerm()
	} else if config.Debug {
		RunScanner(env)
	}
}
