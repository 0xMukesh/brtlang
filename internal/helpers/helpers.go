package helpers

import (
	"fmt"
	"os"

	"github.com/0xmukesh/interpreter/internal/lexer"
	"github.com/0xmukesh/interpreter/internal/tokens"
	"github.com/0xmukesh/interpreter/internal/utils"
)

func ProcessTokens(l *lexer.Lexer, exitOnError bool) []tokens.Token {
	tkns, err := l.LexAll()
	if err != nil {
		errMsg := fmt.Sprintf("[line %d] Error: %s\n", l.Line, err.Error())
		if exitOnError {
			utils.EPrint(errMsg)
		}
		fmt.Fprint(os.Stderr, errMsg)
		return nil
	}

	var filteredTkns []tokens.Token
	for _, tkn := range tkns {
		switch tkn.Type {
		case tokens.IGNORE:
			continue
		case tokens.ILLEGAL:
			errMsg := fmt.Sprintf("[line %d] Error: Unexpected character: %s\n", l.Line, tkn.Literal)
			if exitOnError {
				utils.EPrint(fmt.Sprintf("[line %d] Error: Unexpected character: %s\n", l.Line, tkn.Literal))
			}
			fmt.Fprint(os.Stderr, errMsg)
			continue
		default:
			filteredTkns = append(filteredTkns, tkn)
		}
	}

	return filteredTkns
}
