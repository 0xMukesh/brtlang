package commands

import (
	"fmt"
	"os"

	"github.com/0xmukesh/interpreter/internal/evaluator"
	"github.com/0xmukesh/interpreter/internal/lexer"
	"github.com/0xmukesh/interpreter/internal/parser"
	"github.com/0xmukesh/interpreter/internal/runtime"
	"github.com/0xmukesh/interpreter/internal/tokens"
	"github.com/0xmukesh/interpreter/internal/utils"
)

func EvaluteCmdHandler(src []byte) {
	l := lexer.NewLexer(src)

	tkns, lErr := l.LexAll()
	if lErr != nil {
		utils.EPrint(fmt.Sprintf("%s\n", lErr.Error()))
	}

	hasLexicalErrs := 1

	for i, tkn := range tkns {
		if tkn.Type == tokens.IGNORE {
			tkns = append(tkns[:i], tkns[i+1:]...)
		} else if tkn.Type == tokens.ILLEGAL {
			fmt.Fprintf(os.Stderr, "[line %d] Error: Unexpected character: %s\n", l.Line, tkn.Literal)
			hasLexicalErrs *= 0
		}
	}

	p := parser.NewParser(tkns)

	ast, pErr := p.BuildAst()
	if pErr != nil {
		utils.EPrint(fmt.Sprintf("%s\n", pErr.Error()))
	}

	vars := map[string]runtime.RuntimeValue{}
	env := runtime.NewEnvironment(vars)
	e := evaluator.NewEvaluator(ast, env)

	for {
		val, eErr := e.Evaluate()
		if eErr != nil {
			fmt.Fprintf(os.Stderr, "%s\n", eErr.Error())
		}

		if val == nil {
			break
		} else {
			fmt.Println(val.String())
		}
	}

	if hasLexicalErrs == 0 {
		os.Exit(65)
	} else {
		os.Exit(0)
	}
}
