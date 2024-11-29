package commands

import (
	"fmt"
	"os"

	"github.com/0xmukesh/interpreter/internal/lexer"
	"github.com/0xmukesh/interpreter/internal/parser"
	"github.com/0xmukesh/interpreter/internal/tokens"
	"github.com/0xmukesh/interpreter/internal/utils"
)

func ParseCmdHandler(src []byte) {
	l := lexer.NewLexer(src)

	tkns, err := l.LexAll()
	if err != nil {
		utils.EPrint(fmt.Sprintf("[line %d] Error: %s\n", l.Line, err.Message))
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

	for range tkns {
		node, err := p.Parse()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		}

		if node != nil {
			expr := node.ExtractExpr()

			if expr != nil {
				fmt.Printf("%+v\n", expr.ParseExpr())
			}
		}
	}

	if hasLexicalErrs == 0 {
		os.Exit(65)
	} else {
		os.Exit(0)
	}
}
