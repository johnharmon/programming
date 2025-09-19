package set

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"
)

func (cl *ConcreteLogger) Logln(message string, vars ...any) {
	// fmt.Fprintf(os.Stderr, "Logln called\n")
	rawLogArgs := cl.RawLogArgPool.Get().(*RawLogArgs)
	rawLogArgs.FormatMessage = message
	rawLogArgs.FormatArgs = vars
	rawLogArgs.Timestamp = time.Now().Format(time.StampMicro)
	cl.RawLogCh <- rawLogArgs
}

func (cl *ConcreteLogger) RawLogHandler() {
	for logArgs := range cl.RawLogCh {
		// fmt.Fprintf(os.Stderr, "RawLogHandler received log: %s\n", logArgs.FormatMessage)
		entry := cl.LogEntryPool.Get().(*LogEntry)
		entry.Message = fmt.Sprintf(logArgs.FormatMessage, logArgs.FormatArgs...)
		entry.Timestamp = logArgs.Timestamp
		cl.LogEntryCh <- entry
		cl.RawLogArgPool.Put(logArgs)
	}
	close(cl.LogEntryCh)
}

func (cl *ConcreteLogger) JsonMarshaler() {
	encodeBuffer := bytes.NewBuffer(make([]byte, 0, 2048))
	encoder := json.NewEncoder(encodeBuffer)
	for rawLog := range cl.LogEntryCh {
		encodeSlice := cl.MessageBufferPool.New().([]byte)
		encodeSlice = encodeSlice[:0]
		// fmt.Fprintf(os.Stderr, "JsonMarshaler received Log entry: +%v\n", rawLog)
		encodeBuffer.Reset()
		encoder.Encode(rawLog)
		encodeSlice = append(encodeSlice, encodeBuffer.Bytes()...)
		cl.LogOutput <- encodeSlice
		// os.Stderr.Write(encodeBuffer.Bytes())
		cl.LogEntryPool.Put(rawLog)
	}
	close(cl.LogOutput)
}

func (cl *ConcreteLogger) JsonWriter() {
	// cl.Mu.Lock()
	var flushToken *FlushToken
	bufFlushSize := 1024
	ticker := time.NewTicker(time.Millisecond * 100)
LogWriteLoop:
	for {
		select {
		case msg, ok := <-cl.LogOutput:
			// fmt.Fprintf(os.Stderr, "JsonWriter received Log: %s\n", msg)
			if !ok {
				if cl.ActiveBuffer.Len() > 0 {
					flushToken = <-cl.FlushReceiver
					flushToken.HandledBy = "JsonWriter(): case msg, ok := <- cl.Logch; if !ok {<this>}"
					flushToken.SentBy = "JsonWriter(): case msg, ok := <- cl.Logch; if !ok {<this>}"
					cl.Out.Write(cl.ActiveBuffer.Bytes())
				}
				cl.Mu.Unlock()
				return
			}
			cl.SwapMu.Lock()
			cl.ActiveBuffer.Write(msg)
			os.Stderr.Write(msg)
			cl.SwapMu.Unlock()
			cl.MessageBufferPool.Put(msg)
			if cl.ActiveBuffer.Len() >= bufFlushSize {
				select {
				case flushToken = <-cl.FlushReceiver:
					flushToken.HandledBy = "JsonWriter(): case msg, ok := <- cl.Logch"
					flushToken.SentBy = "JsonWriter(): case msg, ok := <- cl.Logch"
					cl.FlushAndSwapActiveBuffer()
					cl.FlushSender <- flushToken
				default:
					continue LogWriteLoop
				}
			}
		case <-ticker.C:
			if cl.ActiveBuffer.Len() > 0 {
				select {
				case flushToken = <-cl.FlushReceiver:
					flushToken.HandledBy = "JsonWriter(): case <- ticker.C:"
					flushToken.SentBy = "JsonWriter(): case <- ticker.C:"
					cl.FlushAndSwapActiveBuffer()
					cl.FlushSender <- flushToken
				default:
					continue LogWriteLoop
				}
			}
		}
	}
}

func (cl *ConcreteLogger) FlushAndSwapActiveBuffer() {
	go func() {
		cl.SwapMu.Lock()
		cl.ActiveBuffer, cl.FlushBuffer = cl.FlushBuffer, cl.ActiveBuffer
		cl.SwapMu.Unlock()
		flushToken := <-cl.FlushSender
		flushToken.Iteration++
		flushToken.HandledBy = "FlushAndSwapActiveBuffer()"
		flushToken.ReceivedBy = "FlushAndSwapActiveBuffer()"
		cl.Out.Write(cl.FlushBuffer.Bytes())
		//if cl.Out != os.Stderr {
		//	fmt.Fprintf(os.Stderr, cl.FlushBuffer.String())
		//}
		cl.FlushBuffer.Reset()
		cl.FlushReceiver <- flushToken
	}()
}

func (cl *ConcreteLogger) StartAsync() {
	cl.Mu.Lock()
	flushToken := &FlushToken{Iteration: 0}
	cl.FlushReceiver <- flushToken
	cl.RawLogHandler()
	cl.JsonMarshaler()
	cl.JsonWriter()
	for {
		_, ok := <-cl.Done
		if !ok {
			// fmt.Fprintf(os.Stderr, "Logging close received, exiting")
			close(cl.RawLogCh)
		}
	}
}

func (cl *ConcreteLogger) Start() {
	cl.Mu.Lock()
	flushToken := &FlushToken{Iteration: 0}
	go cl.RawLogHandler()
	go cl.JsonMarshaler()
	go cl.JsonWriter()
	cl.FlushReceiver <- flushToken
	for {
		// fmt.Fprintf(os.Stderr, "Waiting on cl.Done to be closed\n")
		_, ok := <-cl.Done
		if !ok {
			// fmt.Fprintf(os.Stderr, "cl.Done closed, closing cl.RawLogCh\n")
			cl.RawLogCh <- &RawLogArgs{FormatMessage: "FlushToken{Iterations: %d HandledBy: %s, SentBy: %s", FormatArgs: []any{flushToken.Iteration, flushToken.HandledBy, flushToken.SentBy}}
			close(cl.RawLogCh)
			break
		}
	}
}

func (cl *ConcreteLogger) Init() {
	// cl.LogCh = make(chan string)
	cl.Mu = &sync.Mutex{}
	cl.FlushMu = &sync.Mutex{}
	cl.SwapMu = &sync.Mutex{}
	cl.ActiveBuffer = &bytes.Buffer{}
	cl.FlushBuffer = &bytes.Buffer{}
	cl.FlushReceiver = make(chan *FlushToken)
	cl.FlushSender = make(chan *FlushToken)
	cl.LogOutput = make(chan []byte, 1000)
	cl.RawLogCh = make(chan *RawLogArgs, 1000)
	cl.LogEntryCh = make(chan *LogEntry, 1000)
	cl.RunCh = make(chan *sync.WaitGroup)
	cl.Done = make(chan struct{})
	cl.LogEntryPool = &sync.Pool{}
	cl.LogEntryPool.New = func() any {
		return &LogEntry{}
	}
	cl.RawLogArgPool = &sync.Pool{}
	cl.RawLogArgPool.New = func() any {
		return &RawLogArgs{}
	}
	cl.MessageBufferPool = &sync.Pool{}
	cl.MessageBufferPool.New = func() any {
		return make([]byte, 2048)
	}

	f, err := os.CreateTemp("./", ".term-reader-logger.json.")
	if err != nil {
		fmt.Printf("Error opening tmp file: %s\n", err)
		cl.Out = os.Stderr
	} else {
		os.Remove("term-reader-logger.json")
		cl.LogFileName = f.Name()
		err := os.Symlink(f.Name(), "term-reader-logger.json")
		if err != nil {
			fmt.Printf("Error creating logger symlink: %s\n", err)
		}
		cl.Out = f
		go cl.Start()
		cl.Logln("Opened New logger at %s", f.Name())
	}
}

func NewConcreteLogger() (cl *ConcreteLogger) {
	cl = &ConcreteLogger{}
	cl.Init()
	return cl
}
