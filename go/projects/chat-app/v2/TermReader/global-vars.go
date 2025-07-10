package main

import (
	"log"
	"os"
)

var (
	COMMANDS                  map[string]Cmd = make(map[string]Cmd)
	COMMAND_ALIASES           map[string]Cmd = make(map[string]Cmd)
	TermHeight, TermWidth                    = GetTermSize()
	GlobalLogger              EphemeralLogger
	CleanupTaskMap            = map[string]*CleanupTask{}
	GenCleanupKey             = CreateCleanupKeyGenerator(CleanupTaskMap)
	InsertCleanupKey          = CreateCleanupKeyInserter(CleanupTaskMap)
	StartCleanupTasks         = CreateCleanupTaskStarter(CleanupTaskMap)
	RegisterCleanupTask       = CreateCleanupTaskRegistrar(CleanupTaskMap)
	LOGGER_CLEANUP_UNIQUE_KEY = "LOGGER_CLEANUP"
	KeyActionTree             map[byte]*KeyAction
	TERM_CLEAR_LINE           = []byte{0x1b, '[', '2', 'K'}
	TERM_CLEAR_SCREEN         = []byte{0x1b, '[', '2', 'J'}
	InsertModeDispatchMap     = map[string]Action{}
	NormalModeDispatchMap     = map[int]CommandEntry{}
	MotionDispatchMap         = map[string]Motion{}
	Commands                  = map[string]Cmd{}
)

func SetDefaultLogOutput() {
	log.Default().SetOutput(os.Stderr)
}

func InitCoreCommands() {
	COMMANDS["edit"] = Cmd{ExecFunc: Edit}
	COMMANDS["write"] = Cmd{ExecFunc: Write}
}

func InitCommandShortcuts() {
	COMMANDS["e"] = Cmd{ExecFunc: Edit}
	COMMANDS["w"] = Cmd{ExecFunc: Write}
}

func InitGlobalVars() {
	GlobalLogger = NewConcreteLogger()
	SetDefaultLogOutput()
	InitCoreCommands()
	InitCommandShortcuts()
	InitKeyActionTree()
}

func InitKeyActionTree() {
	KeyActionTree = make(map[byte]*KeyAction)
	KeyActionTree[0x1b] = NewKeyAction(false, "Escape", false, 0x1b)
	err := InitializeArrowKeys()
	if err != nil {
		log.Fatalf(err.Error())
	}
	err = InitializeControlCodes()
	if err != nil {
		log.Fatal(err.Error())
	}
	InitAsciiPrintableCharactersWithPrintAction()
}

func InitNormalModeCommands() {
	NormalModeDispatchMap['h'] = CommandEntry{
		SuffixRequired: false,
		AcceptsMotion:  false,
		MustBeFirst:    false,
		ExecFunc:       NormalHandleLeftMoveCmd,
	}

	NormalModeDispatchMap['j'] = CommandEntry{
		SuffixRequired: false,
		AcceptsMotion:  false,
		MustBeFirst:    false,
		ExecFunc:       NormalHandleDownMoveCmd,
	}

	NormalModeDispatchMap['k'] = CommandEntry{
		SuffixRequired: false,
		AcceptsMotion:  false,
		MustBeFirst:    false,
		ExecFunc:       NormalHandleUpMoveCmd,
	}

	NormalModeDispatchMap['l'] = CommandEntry{
		SuffixRequired: false,
		AcceptsMotion:  false,
		MustBeFirst:    false,
		ExecFunc:       NormalHandleRightMoveCmd,
	}
}
