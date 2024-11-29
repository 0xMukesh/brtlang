package commands

import (
	"fmt"
	"os"

	"github.com/0xmukesh/interpreter/internal/helpers"
	"github.com/0xmukesh/interpreter/internal/lexer"
	"github.com/0xmukesh/interpreter/internal/parser"
)

func ParseCmdHandler(src []byte) {
	l := lexer.NewLexer(src)
	tkns := helpers.ProcessTokens(l, true)
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
}
