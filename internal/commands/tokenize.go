package commands

import (
	"fmt"
	"os"

	"github.com/0xmukesh/interpreter/internal/lexer"
	"github.com/0xmukesh/interpreter/internal/tokens"
)

func TokenizeCmdHandler(src []byte) {
	l := lexer.NewLexer(src)

	tkns, err := l.LexAll()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
	}

	hasLexicalErrs := 1

	for _, tkn := range tkns {
		if tkn.Type == tokens.ILLEGAL {
			fmt.Fprintf(os.Stderr, "[line %d] Error: Unexpected character: %s\n", l.Line, tkn.Literal)
			hasLexicalErrs *= 0
		} else if tkn.Type == tokens.IGNORE {
			continue
		} else {
			msg := fmt.Sprintf("%s %s %s", tkn.Type.String(), tkn.Lexeme, tkn.Literal)
			fmt.Println(msg)
		}
	}

	if hasLexicalErrs == 0 {
		os.Exit(65)
	} else {
		os.Exit(0)
	}
}
