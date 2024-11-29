package main

import (
	"os"

	"github.com/0xmukesh/interpreter/internal/commands"
	"github.com/0xmukesh/interpreter/internal/utils"
)

func main() {
	args := os.Args

	if len(args) < 3 {
		utils.EPrint("invalid usage")
	}

	command := args[1]
	filename := args[2]

	src, err := os.ReadFile(filename)
	if err != nil {
		utils.EPrint(err.Error())
	}

	if command == "tokenize" {
		commands.TokenizeCmdHandler(src)
	} else if command == "parse" {
		commands.ParseCmdHandler(src)
	} else if command == "evalute" || command == "eval" {
		commands.EvaluteCmdHandler(src)
	} else if command == "ast" {
		commands.AstCmdHandler(src)
	} else if command == "run" {
		commands.RunCmdHandler(src)
	} else {
		utils.EPrint("invalid command\n")
	}
}
