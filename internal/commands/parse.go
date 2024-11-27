package commands

import (
	"fmt"
	"os"

	"github.com/0xmukesh/interpreter/internal/lexer"
	"github.com/0xmukesh/interpreter/internal/parser"
	"github.com/0xmukesh/interpreter/internal/tokens"
)

func ParseCmdHandler(src []byte) {
	lexer := lexer.NewLexer(src)
	var tkns []tokens.Token

	for {
		lexer.ReadChar()

		tkn, err := lexer.Lex()
		if err != nil {
			fmt.Fprintf(os.Stderr, "[line %d] Error: %s\n", lexer.Line, err.Message)
		}

		if tkn != nil {
			tkns = append(tkns, *tkn)
		}

		if lexer.Char == 0 {
			break
		}
	}

	parser := parser.NewParser(tkns)

	for range tkns {
		node := parser.Parse()
		if node != nil {
			fmt.Printf("%+v\n", node.Expr.ParseExpr())
		}
	}
}
