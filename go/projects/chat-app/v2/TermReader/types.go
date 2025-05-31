package main

import (
	"bytes"
	"io"
	"os"
)

type EphemeralLogger interface {
	Log(string, ...any)
	Logln(string, ...any)
}

type ConcreteLogger struct {
	Out   io.Writer
	LogCh chan string
	RunCh chan interface{}
}

type KeyAction struct {
	Children   map[byte]*KeyAction
	Value      []byte
	IsTerminal bool
	Action     string
	PrintRaw   bool
	FromPool   bool
}

type State struct {
	Windows      []*Window
	ActiveWindow *Window
}

type MainConfig struct {
	LogConfig *LogConfig
	In        io.Reader
	Out       io.Writer
	State     State
}

type Window struct { // Represents a sliding into its backing buffer of Window.Buf as well as the space it takes up in the terminal window
	TermTopLine       int
	BufTopLine        int
	StartIndex        int
	StartLine         int
	Buf               *DisplayBuffer
	Height            int
	Width             int
	StartCol          int
	CursorLine        int
	CursorCol         int
	EndIndex          int
	RawStartIndex     int
	RawEndIndex       int
	Out               io.Writer
	EventChan         chan *KeyAction
	RawEventChan      chan []byte
	KeyActionReturner chan *KeyAction
	Logger            EphemeralLogger
}

type LogConfig struct {
	File *os.File
	Link *os.File
}

type PooledKeyAction struct {
	KA       *KeyAction
	FromPool bool
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

type Env struct {
	Config       *FlagConfig
	DebugWriter  io.Writer
	OutputHeader string
	OutputFooter string
	OutputPrefix string
	OutputSuffix string
	Logger       EphemeralLogger
}
type DisplayBuffer struct { // This represents the full backing buffer to any window view
	RawBuf         []byte
	Lines          [][]byte
	DisplayedLines [][]byte
	Size           int
	ActiveLine     int
	TopLine        int
	Height         int
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
	Window     bool
	Raw        bool
	Logs       io.Writer
}

type FormatInfo struct {
	OutputRaw      []byte
	OutputLines    [][]byte // raw bytes representing the output
	WrappingLength int      // how long the output prefix and suffix combined is
	TermWidth      int      // how many characters wide the terminal output is
	TermLength     int      // how many lines long including headers and footers the ouput is
}

type VirtualBuffer struct {
	Buf [][]byte
}

type Cell struct {
	formatInfo            *FormatInfo
	RawContent            *bytes.Buffer
	RawInput              *bytes.Buffer
	ContentReader         *bytes.Reader
	Out                   io.Writer
	VBuf                  [][]byte
	In                    io.Reader
	DisplayContent        *bytes.Buffer
	CursorPosition        int
	LogicalCursorPosition int
	DisplayCursorPosition int
	CursorLine            int
	CursorColumn          int
	DisplayBuffer         *DisplayBuffer
	Window                *Window
	VirtualBuffer         [][]byte
	DebugInfo             []string
	CellHistory           []*Cell
	ActiveLineIdx         int
	ActiveLineLength      int
	LogicalRowIdx         int
	EffecitveRowIdx       int
	LogCh                 chan string
	BufferLen             int
	Logger                io.Writer
	LogFile               *os.File
	LogLink               string
}
