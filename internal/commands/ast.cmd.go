package commands

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/0xmukesh/interpreter/internal/lexer"
	"github.com/0xmukesh/interpreter/internal/parser"
	"github.com/0xmukesh/interpreter/internal/tokens"
	"github.com/0xmukesh/interpreter/internal/utils"
)

func AstCmdHandler(src []byte) {
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

	b, err := json.MarshalIndent(ast, "", " ")
	if err != nil {
		utils.EPrint(fmt.Sprintf("%s\n", pErr.Error()))
	}

	fmt.Println(string(b))
}
