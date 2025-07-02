package main

import (
	"io"
	"os"
	"sync"
	"time"

	"golang.org/x/term"
)

func ReadInput(input io.Reader, out chan byte) {
	b := make([]byte, 1)
	for {
		input.Read(b)
		out <- b[0]
	}
}

func ReadWithTimeout(buf []byte, in chan byte, timeout time.Duration) (n int) {
	// GlobalLogger.Logln("Reading with timeout")
	timeChan := time.After(timeout)
	var b byte = 0
	bytesRead := 0
listenLoop:
	for {
		select {
		case b = <-in:
			// GlobalLogger.Logln("Caught multi byte sequence")
			if bytesRead < len(buf) {
				buf[bytesRead] = b
				bytesRead++
			}
		case <-timeChan:
			break listenLoop
		}
	}
	return bytesRead
}

func InputParser(in chan (byte), out chan []byte, bufPool *sync.Pool, timeout time.Duration) {
	bufPool.New = func() any {
		return make([]byte, 1, 8)
	}
	for b := range in {
		buf := bufPool.Get().([]byte)[0:1]
		buf[0] = b
		if buf[0] == '\x1b' {
			// clear(buf[1:8])
			n := ReadWithTimeout(buf[1:8], in, timeout)
			// GlobalLogger.Logln("Read %d bytes from timeout", n)
			buf = buf[0 : 1+n]
		}
		// GlobalLogger.Logln("Sending %b to key action generator", buf)
		out <- buf
	}
}

func KeyActionGenerator(in chan []byte, out chan *KeyAction, closer chan struct{}, kaPool *sync.Pool, bufPool *sync.Pool) {
	var ka *KeyAction
	for b := range in {
		if len(b) == 1 {
			if b[0] == 3 {
				closer <- struct{}{}
				return
			} else if b[0] >= 0x20 && b[0] <= 0x7E {
				GlobalLogger.Logln("Getting *KeyAction from pool")
				ka = kaPool.Get().(*KeyAction)
				ka.Value = ka.Value[0:len(b)]
				copy(ka.Value, b)
			} else {
				ka = ValidateSequence(b)
			}
		} else {
			ka = ValidateSequence(b)
		}
		if ka == nil {
			GlobalLogger.Logln("Getting *KeyAction from pool")
			ka = kaPool.Get().(*KeyAction)
			ka.Value = ka.Value[0:len(b)]
			copy(ka.Value, b)
			ka.Action = "Unknown"
		}
		out <- ka
		bufPool.Put(b[0:cap(b)])
	}
}

func KeyActionGenerator2(in chan []byte, out chan *KeyAction, closer chan struct{}, kaPool *sync.Pool, bufPool *sync.Pool) {
	var ka *KeyAction
	for b := range in {
		if len(b) == 1 {
			if b[0] == 3 {
				closer <- struct{}{}
				return
			} else {
				ka = ValidateSequence(b)
				if ka == nil {
					GlobalLogger.Logln("Getting *KeyAction from pool")
					ka = kaPool.Get().(*KeyAction)
					ka.Value = ka.Value[0:len(b)]
					copy(ka.Value, b)
					ka.Action = "Unknown"
				}
				out <- ka
				bufPool.Put(b[0:cap(b)])
			}
		}
	}
}

func StartupListeners(input io.Reader, closer chan struct{}, terminalInput chan byte, inputSequences chan []byte, keyActions chan *KeyAction, inputParserPool *sync.Pool, kaPool *sync.Pool) {
	go ReadInput(input, terminalInput)
	go InputParser(terminalInput, inputSequences, inputParserPool, time.Millisecond*25)
	go KeyActionGenerator(inputSequences, keyActions, closer, kaPool, inputParserPool)
}

func MainEventHandler2(mc *MainConfig) {
	GlobalLogger.Logln("Making closer channel")
	closer := make(chan struct{})
	GlobalLogger.Logln("Making directTerminalInput channel")
	directTerminalInput := make(chan byte, 1)
	GlobalLogger.Logln("Making terminalInputSequences channel")
	terminalInputSequences := make(chan []byte, 1)
	GlobalLogger.Logln("Making keyActions channel")
	keyActions := make(chan *KeyAction, 1)
	GlobalLogger.Logln("Making inputParserPool pool")
	inputParserPool := &sync.Pool{}
	GlobalLogger.Logln("Making kaPool pool")
	kaPool := MakeKeyActionPool()
	StartupListeners(mc.In, closer, directTerminalInput, terminalInputSequences, keyActions, inputParserPool, kaPool)
	gl := GlobalLogger.(*ConcreteLogger)
	RegisterCleanupTask(gl.RunCh, gl.Cleanup, LOGGER_CLEANUP_UNIQUE_KEY, true)
	var ka *KeyAction
	fd := int(os.Stdin.Fd())
	gl.Logln("Setting terminal to raw mode")
	oldState, err := term.MakeRaw(fd)
	if err != nil {
		panic(err)
	}
	gl.Logln("Making closer channel")
	gl.Logln("Making *KeyAction pool")
	sp := MakeKeyActionPool() // Create a pool of *SequenceNode references for dispatching normal ascii printable characters on the event channel for the window (trying to avoid as much re-allocation as possible
	go mc.State.ActiveWindow.RunKeyActionReturner(sp)
	gl.Logln("Making *KeyAction return channel for the pool")
	keyActionReturner := make(chan *KeyAction, 1000)
	gl.Logln("Spinning off goroutine for returning *KeyActions to the pool")
	go ReturnKeyActionsToPool(sp, keyActionReturner)
	gl.Logln("Spinning off cleanup goroutine")
	go Cleanup(closer, fd, oldState, CleanupTaskMap)
	gl.Logln("Creating byte handler from closure")
	gl.Logln("Entering main event loop")
	for ka = range keyActions {
		mc.State.ActiveWindow.EventChan <- ka
	}
}
