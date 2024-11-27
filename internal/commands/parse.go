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
	lexer := lexer.NewLexer(src)
	var tkns []tokens.Token

	for {
		lexer.Read()

		tkn, err := lexer.Lex()
		if err != nil {
			// exit if lexer throws an error
			utils.EPrint(fmt.Sprintf("[line %d] Error: %s\n", lexer.Line, err.Message))
		}

		if tkn != nil {
			tkns = append(tkns, *tkn)
		}

		if lexer.Char == 0 {
			break
		}
	}

	for i, tkn := range tkns {
		if tkn.Type == tokens.IGNORE {
			tkns = append(tkns[:i], tkns[i+1:]...)
		}
	}

	parser := parser.NewParser(tkns)

	for range tkns {
		node, err := parser.Parse()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err.String())
		}

		if node != nil {
			fmt.Printf("%+v\n", node.Expr.ParseExpr())
		}
	}
}
