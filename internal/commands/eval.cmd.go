package commands

import (
	"fmt"
	"os"

	"github.com/0xmukesh/interpreter/internal/evaluator"
	"github.com/0xmukesh/interpreter/internal/helpers"
	"github.com/0xmukesh/interpreter/internal/lexer"
	"github.com/0xmukesh/interpreter/internal/parser"
	"github.com/0xmukesh/interpreter/internal/runtime"
	"github.com/0xmukesh/interpreter/internal/utils"
)

func EvaluteCmdHandler(src []byte) {
	l := lexer.NewLexer(src)
	tkns := helpers.ProcessTokens(l, true)
	p := parser.NewParser(tkns)

	ast, pErr := p.BuildAst()
	if pErr != nil {
		utils.EPrint(fmt.Sprintf("%s\n", pErr.Error()))
	}

	vars := map[string]runtime.RuntimeValue{}
	globaEnv := runtime.NewEnvironment(vars)
	runtime := runtime.NewRuntime([]*runtime.Environment{globaEnv})
	e := evaluator.NewEvaluator(ast, runtime)

	for {
		val, eErr := e.Evaluate()
		if eErr != nil {
			fmt.Fprintf(os.Stderr, "%s\n", eErr.Error())
		}

		if val == nil {
			break
		}

		fmt.Println(val.String())
	}
}
