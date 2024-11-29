package commands

import (
	"fmt"

	"github.com/0xmukesh/interpreter/internal/evaluator"
	"github.com/0xmukesh/interpreter/internal/helpers"
	"github.com/0xmukesh/interpreter/internal/lexer"
	"github.com/0xmukesh/interpreter/internal/parser"
	"github.com/0xmukesh/interpreter/internal/runner"
	"github.com/0xmukesh/interpreter/internal/runtime"
	"github.com/0xmukesh/interpreter/internal/utils"
)

func RunCmdHandler(src []byte) {
	l := lexer.NewLexer(src)
	tkns := helpers.ProcessTokens(l, true)
	p := parser.NewParser(tkns)

	ast, err := p.BuildAst()
	if err != nil {
		utils.EPrint(fmt.Sprintf("%s\n", err.Error()))
	}

	vars := runtime.RuntimeVarMapping{}
	globaEnv := runtime.NewEnvironment(vars)
	runtime := runtime.NewRuntime(&[]runtime.Environment{*globaEnv})
	e := evaluator.NewEvaluator(ast, runtime)
	r := runner.NewRunner(ast, runtime, e)

	for !r.IsAtEnd() {
		r.Run()
	}
}
