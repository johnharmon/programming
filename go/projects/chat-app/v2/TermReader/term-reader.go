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

type DisplayBuffer struct {
	Buffer []byte
	Lines  [][]byte
	Size   int
}

func (db *DisplayBuffer) AllocateLines(size int) {
	for i := range db.Lines {
		db.Lines[i] = make([]byte, size)
	}
}

func (db *DisplayBuffer) Write(p []byte) (int, error) {
	db.Buffer = append(db.Buffer, p...)
	return len(p), nil
}

func (db *DisplayBuffer) Reset() {
	if cap(db.Buffer) > 4096 {
		db.Buffer = make([]byte, 4096)
	} else {
		clear(db.Buffer)
		db.Buffer = db.Buffer[:0]
	}
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
	CursorLine            int
	CursorColumn          int
	DisplayBuffer         *DisplayBuffer
	DebugInfo             []string
	CellHistory           []*Cell
	ActiveLineIdx         int
	LogicalRowIdx         int
	EffecitveRowIdx       int
}

func (oc *Cell) Display(o io.Writer, env *Env) {
	formattedContent := WrapOutput(env, oc.RawContent.Bytes())
	fmt.Fprint(o, formattedContent)
}

func (oc *Cell) OverWrite(newContent []byte) {
}

func (cell *Cell) IncrementActiveLine(incr int) {
	nextLine := cell.ActiveLineIdx + incr
	if nextLine > 0 || nextLine < cell.DisplayBuffer.Size-1 {
		cell.ActiveLineIdx = nextLine
		return
	} else {
		if nextLine < 0 {
			cell.ActiveLineIdx = cell.DisplayBuffer.Size - 1
			return
		} else if nextLine >= cell.DisplayBuffer.Size {
			cell.ActiveLineIdx = 0
			return
		}
	}
}

func (cell *Cell) ScrollUp(newLine []byte) {
	numLines := len(cell.DisplayBuffer.Lines)
	tmpLine := cell.DisplayBuffer.Lines[numLines-1][:0]
	copy(cell.DisplayBuffer.Lines[1:], cell.DisplayBuffer.Lines[:numLines-1])
	tmpLine = append(tmpLine, newLine...)
	cell.DisplayBuffer.Lines[0] = tmpLine
	cell.IncrementActiveLine(1)
}

func (cell *Cell) ScrollDown(newLine []byte) {
	numLines := len(cell.DisplayBuffer.Lines)
	tmpLine := cell.DisplayBuffer.Lines[0][:0]
	copy(cell.DisplayBuffer.Lines[:numLines-1], cell.DisplayBuffer.Lines[1:])
	tmpLine = append(tmpLine, newLine...)
	cell.DisplayBuffer.Lines[numLines-1] = tmpLine
	cell.IncrementActiveLine(-1)
}

//func (cell *Cell) IncrementActiveLine(incr int) {
//	newLine := cell.ActiveLineIdx + incr
//	if newLine > cell.DisplayBuffer.Size -1
//	cell.IncrementActiveLine(1)
//}

func (cell *Cell) EnterDisplayLoop() {
	var (
		isMod  bool
		modSeq *ModificationSequence
		out    io.Writer
		buf    = make([]byte, 1)
	)
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
						fmt.Fprintf(out, "\r\x1b[2K%s", newLine)
						DisplayDebugInfo(cell, "EscapeSequenceConditional, ForceRedraw", []string{})
						cell.RawInput.Write(esc.Bytes)
					} else {
						cell.RawInput.Write(esc.Bytes)
						fmt.Fprintf(out, "%s", esc.Bytes)
						switch esc.Name {
						case "LeftArrow":
							if cell.DisplayCursorPosition > 0 {
								cell.DisplayCursorPosition--
							}
							cell.LogicalCursorPosition += len(esc.Bytes)
							fmt.Fprintf(out, "%s", esc.Bytes)
						case "RightArrow":
							if cell.DisplayCursorPosition < len(cell.DisplayContent.Bytes()) {
								cell.DisplayCursorPosition++
							}
							cell.LogicalCursorPosition += len(esc.Bytes)
							fmt.Fprintf(out, "%s", esc.Bytes)
						}
						DisplayDebugInfo(cell, "EscapeSequenceConditional, NoRedraw", []string{})
						cell.RawInput.Write(esc.Bytes)
					}
				}
			} else {
				bslice := buf[0:1]
				cell.RawInput.WriteByte(b)
				if b == 13 {
					cell.Out.Write([]byte("\n\r"))
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

type DisplayWrapper struct {
	TopPattern    string
	BottomPattern string
	LinePrefix    string
	LineSuffix    string
}

func (es ModificationSequence) String() string {
	return es.Name
}

func MakeRawTerm(config *FlagConfig) {
	fd := int(os.Stdin.Fd())
	oldState, err := term.MakeRaw(fd)
	if err != nil {
		panic(err)
	}
	defer term.Restore(fd, oldState)
	if config.Raw {
		RawTermInterface(true)
	} else {
		RawTermInterface(false)
	}
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

func CloneBuffer(b *bytes.Buffer) []byte {
	temp := b.Bytes()
	clone := make([]byte, len(temp))
	copy(clone, temp)
	return clone
}

func matchModificationSequence(sequence []byte) {}

func HandleEscabeByte(b byte, out io.Writer) {
}

func Display(cell *Cell) {
	// dw := &DisplayWrapper{TopPattern: "-", BottomPattern: "-", LinePrefix: "| ", LineSuffix: " |"}
	// lines := bytes.Split(cell.DisplayBuffer.Buffer, []byte("\n"))
	// numLines := len(lines)
}

func (cell Cell) FindCursorCoordFromPos() (row int, col int) {
	lines := bytes.Split(cell.DisplayBuffer.Buffer, []byte("\n"))
	var cumulativeBytes int = 0
	for idx, line := range lines {
		ll := len(line)
		if ll+cumulativeBytes > cell.CursorPosition {
			row = idx
			col = cell.CursorPosition - cumulativeBytes
			return row, col
		} else {
			cumulativeBytes += ll + 1 // + 1 on the end for the \n we lost during the split, it still exists in the original display buffer
		}
	}
	return row, col
}

func WrapOutput2(dw *DisplayWrapper, cell *Cell) {
	displayLines := bytes.Split(cell.DisplayBuffer.Buffer, []byte("\n"))
	numLines := 2 + len(cell.DebugInfo) + len(displayLines)
	jumpUp := numLines - 1 - len(displayLines) // For now, this *should* be the number of lines we need to move the cursor up after we finish printing everything, assuming we are at the end of the interactive buffer
	row, col := cell.FindCursorCoordFromPos()
}

func NewDefaultCell() (cell *Cell) {
	cell = &Cell{}
	cell.Out = os.Stdout
	cell.In = os.Stdin
	cell.DisplayContent = &bytes.Buffer{}
	cell.DisplayCursorPosition = 0
	cell.RawInput = &bytes.Buffer{}
	cell.DisplayBuffer = &DisplayBuffer{Buffer: make([]byte, 4096), Lines: make([][]byte, 100)}
	cell.DisplayBuffer.AllocateLines(4096)
	return cell
}

func RawTermInterface(toStdout bool) {
	var (
		isMod  bool
		modSeq *ModificationSequence
		out    io.Writer
	)
	cell := NewDefaultCell()

	if toStdout {
		out = cell.Out
	} else {
		out = cell.DisplayBuffer
	}
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
						fmt.Fprintf(out, "\r\x1b[2K%s", newLine)
						DisplayDebugInfo(cell, "EscapeSequenceConditional, ForceRedraw", []string{})
						cell.RawInput.Write(esc.Bytes)
					} else {
						cell.RawInput.Write(esc.Bytes)
						fmt.Fprintf(out, "%s", esc.Bytes)
						switch esc.Name {
						case "LeftArrow":
							if cell.DisplayCursorPosition > 0 {
								cell.DisplayCursorPosition--
							}
							cell.LogicalCursorPosition += len(esc.Bytes)
							fmt.Fprintf(out, "%s", esc.Bytes)
						case "RightArrow":
							if cell.DisplayCursorPosition < len(cell.DisplayContent.Bytes()) {
								cell.DisplayCursorPosition++
							}
							cell.LogicalCursorPosition += len(esc.Bytes)
							fmt.Fprintf(out, "%s", esc.Bytes)
						}
						DisplayDebugInfo(cell, "EscapeSequenceConditional, NoRedraw", []string{})
						cell.RawInput.Write(esc.Bytes)
					}
				}
			} else {
				bslice := buf[0:1]
				cell.RawInput.WriteByte(b)
				if b == 13 {
					fmt.Fprintf(out, "\n\rYou typed: %s\r\n", cell.RawInput.Bytes())
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
	// fmt.Fprintf(cell.Out, "\r\x1b[2K\r%s", cell.DisplayContent.Bytes())
}

func (cell *Cell) WriteDisplayBytes(b []byte) {
	extra := []string{}
	blen := len(b)
	if cell.DisplayCursorPosition == 0 {
		temp := cell.DisplayContent.Bytes()
		content := make([]byte, len(temp))
		copy(content, temp)
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
		// temp := append([]byte{}, cell.DisplayContent.Bytes()...)
		temp := cell.DisplayContent.Bytes()
		clone := make([]byte, len(temp))
		copy(clone, temp)
		cell.DisplayContent.Reset()
		before := clone[0:cell.DisplayCursorPosition]
		after := clone[cell.DisplayCursorPosition:]
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

func (cell *Cell) WriteDisplayBytesByBuffer(b []byte) {
	extra := []string{}
	blen := len(b)
	if cell.DisplayCursorPosition == 0 {
		activeLine := cell.DisplayBuffer.Lines[cell.ActiveLineIdx]
		activeLine = activeLine[0 : len(activeLine)+blen]
		copy(activeLine[len(b):], activeLine[0:len(activeLine)-len(b)])
		copy(activeLine[:len(b)], b)
		extra = append(extra, string(b))
		extra = append(extra, string(activeLine))
		cell.DisplayCursorPosition += blen
		cell.Redraw()
		DisplayDebugInfo(cell, "DisplayCursorPosition = 0", extra)
	} else if cell.DisplayCursorPosition == len(cell.DisplayContent.Bytes()) {
		cell.DisplayContent.Write(b)
		cell.DisplayCursorPosition += blen
		cell.LogicalCursorPosition += blen
		fmt.Fprintf(cell.Out, "%s", b)
		DisplayDebugInfo(cell, "DisplayCursorPosition = end of display content", extra)
	} else {
		// temp := append([]byte{}, cell.DisplayContent.Bytes()...)
		activeLine := cell.DisplayBuffer.Lines[cell.ActiveLineIdx]
		clone := make([]byte, len(activeLine))
		copy(clone, activeLine)
		before := clone[0:cell.DisplayCursorPosition]
		after := clone[cell.DisplayCursorPosition:]
		activeLine = append(activeLine[:0], before, b, after)
		cell.DisplayContent.Write(before)
		extra = append(extra, string(before))
		cell.DisplayContent.Write(b)
		extra = append(extra, string(b))
		cell.DisplayContent.Write(after)
		extra = append(extra, string(after))
		cell.DisplayCursorPosition += blen
		cell.Redraw()
		DisplayDebugInfo(cell, "DisplayCursorPostition = inside display content", extra)
	}
}

func InsertAt(a []byte, b []byte, startIdx int) []byte {
	al := len(a)
	bl := len(b)
	if cap(a) >= len(a)+len(b) {
		if bl == 1 {
			a = a[0 : al+1]
			for i := al - 2; i >= startIdx; i-- {
				a[i+1] = a[i]
				if i == startIdx {
					a[i] = b[0]
				}
			}
		} else {
			a = a[0 : al+bl]
			for i := al - 1 - bl; i >= startIdx; i-- {
				a[i+bl] = a[i]
			}
		}
		copy(a[startIdx:startIdx+bl], b)
		return a
	}
	tmp := make([]byte, 0, (al+bl)*2)
	tmp = append(tmp, a[0:startIdx]...)
	tmp = append(tmp, b...)
	tmp = append(tmp, a[startIdx:]...)
	return tmp
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
		MakeRawTerm(config)
	} else if config.Debug {
		RunScanner(env)
	}
}
