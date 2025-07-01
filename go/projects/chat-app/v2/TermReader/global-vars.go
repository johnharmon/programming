package main

var COMMANDS map[string]Cmd = make(map[string]Cmd)

func InitCommands() {
	COMMANDS["edit"] = Cmd{ExecFunc: Edit}
}
