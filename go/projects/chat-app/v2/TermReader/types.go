package main

import (
	"bytes"
	"io"
	"os"
	"regexp"
	"sync"
)

type EphemeralLogger interface {
	Logln(string, ...any)
	Cleanup()
}

type LogEntry struct {
	Message   string `json:"message"`
	Timestamp any    `json:"timestamp"`
}

type RawLogArgs struct {
	FormatMessage string
	FormatArgs    []any
	Timestamp     string
}

type FlushToken struct {
	Iteration  int
	HandledBy  string
	SentBy     string
	ReceivedBy string
	Values     map[string]any
}

type ConcreteLogger struct {
	ActiveBuffer      *bytes.Buffer
	FlushBuffer       *bytes.Buffer
	Out               io.Writer
	Mu                *sync.Mutex
	FlushMu           *sync.Mutex
	SwapMu            *sync.Mutex
	FlushSender       chan *FlushToken
	FlushReceiver     chan *FlushToken
	LogOutput         chan []byte
	LogEntryCh        chan *LogEntry
	RawLogCh          chan *RawLogArgs
	RunCh             chan *sync.WaitGroup
	Done              chan struct{}
	LogFileName       string
	LogEntryPool      *sync.Pool
	RawLogArgPool     *sync.Pool
	MessageBufferPool *sync.Pool
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

type CleanupTask struct {
	Closer chan *sync.WaitGroup
	Name   string
	Func   func()
	Start  bool
	// A CleanupTask represents a task to run that will wait to receive a waitgroup from the closer channel
	// All cleanup tasks are run via a simple func() call, often requiring some sort of wrapper/closure for the real logic
	// They are also backgrounded as long running tasks and notified via a channel carrying the main cleanup orchestrator WaitGroup which it only needs to call wg.Done()
}

type Window struct { // Represents a sliding into its backing buffer of Window.Buf as well as the space it takes up in the terminal window
	TermTopLine       int    // Actual terminal line that represents the top of this window
	BufTopLine        int    // unused
	StartIndex        int    // unused
	FileName          string // name/path of the file that this buffer represents
	NeedRedraw        bool   // boolean indicating if the window needs to be redrawn
	StartLine         int
	Height            int // how many rows high this window is in the terminal
	Width             int // how many columns/characters wide this terminal window is
	StartCol          int
	CursorLine        int             // the line **in the backing buffer** that the terminal is on
	CursorCol         int             // the column in the terminal that the cursor should be on for normal buffer operations
	DesiredCursorCol  int             // the column that the cursor *wants* to be on, really only for allowing the cursor to move back to a desired position after scrolling across lines that are shorter than a pervous one
	Mode              int             // insert, normal, visual, cmd, used to decide key sequence routing
	PrevCursorCol     int             // essentially a storage field when the cursor column only needs to be updated temporarily
	CmdCursorCol      int             // column that the cursor needs to be in the command line, separately tracked for ease of rendering and modification
	EndIndex          int             // unused
	RawStartIndex     int             // unused
	RawEndIndex       int             // unused
	CmdBuf            []byte          // the command buffer
	ExpectingInputBuf []byte          // the buffer to store key input in while waiting for an actionable sequence to be read
	Buf               *DisplayBuffer  // pointer to the backing display buffer for this window
	Out               io.Writer       // simple writer for where things should write to if say saving this buffer
	EventChan         chan *KeyAction // Channel that the main even handler routes resolved keyactions to
	RawEventChan      chan []byte
	KeyActionReturner chan *KeyAction // channel on which to return pool-generated KeyActions to
	Logger            EphemeralLogger // Mostly Unused, reference to any struct that implements the interface
	ActionCount       int
	MotionCount       int
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

type Cmd struct {
	ExecFunc func(*Window, ...string) error
}

type Action struct {
	ExecFunc    func(*Window, *ActionContext)
	MustBeFirst bool
	Context     *ActionContext
}

type ActionContext struct {
	Count     int    // Number of characters that make up the full action
	FullInput []byte // Full action input including the actionable key itself
	Prefix    []byte // arguments/modifiers that came before the action key
	Suffix    []byte // arguments/modifiers that came after the action key
}

type ActionWithContext struct {
	Action  *Action
	Context *ActionContext
}

type CommandEntry struct {
	AcceptsCount  bool
	AcceptsSuffix bool
	MustBeFirst   bool
	ExecFunc      func(*Window, *ActionContext)
}

type Motion struct {
	Expression   regexp.Regexp
	AcceptsCount bool
}
