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
	vars := make(runtime.RuntimeVarMapping)
	funcs := make(runtime.RuntimeFuncMapping)
	globaEnv := runtime.NewEnvironment(vars, funcs, nil)
	runtime := runtime.NewRuntime(&[]runtime.Environment{*globaEnv})

	l := lexer.NewLexer(src)
	tkns := helpers.ProcessTokens(l, true)
	p := parser.NewParser(tkns, runtime)

	programAst, err := p.BuildAst()
	if err != nil {
		utils.EPrint(fmt.Sprintf("%s\n", err.Error()))
	}

	e := evaluator.NewEvaluator(programAst, runtime)
	r := runner.NewRunner(programAst, runtime, e)

	for !r.IsAtEnd() {
		r.Run()
	}
}
