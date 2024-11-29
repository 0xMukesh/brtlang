package commands

import (
	"encoding/json"
	"fmt"

	"github.com/0xmukesh/interpreter/internal/helpers"
	"github.com/0xmukesh/interpreter/internal/lexer"
	"github.com/0xmukesh/interpreter/internal/parser"
	"github.com/0xmukesh/interpreter/internal/utils"
)

func AstCmdHandler(src []byte) {
	l := lexer.NewLexer(src)
	tkns := helpers.ProcessTokens(l, true)
	p := parser.NewParser(tkns)

	ast, pErr := p.BuildAst()
	if pErr != nil {
		utils.EPrint(fmt.Sprintf("%s\n", pErr.Error()))
	}

	b, err := json.MarshalIndent(ast, "", " ")
	if err != nil {
		utils.EPrint(fmt.Sprintf("%s\n", err.Error()))
	}

	fmt.Println(string(b))
}
