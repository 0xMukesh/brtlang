package parser

import (
	"github.com/0xmukesh/interpreter/internal/ast"
	"github.com/0xmukesh/interpreter/internal/tokens"
)

func (p *Parser) parseGroupingExpr() (*ast.AstNode, *ParserError) {
	nodePtr, err := p.Parse()
	if err != nil {
		return nil, err
	}

	if nodePtr != nil {
		err := NewParserError(MISSING_RPAREN, p.curr().Lexeme, p.curr().Line)
		p.consume(tokens.RIGHT_PAREN, *err)

		node := *nodePtr
		return ast.NewAstNode(ast.EXPR, ast.NewGroupingExpr(node, p.curr().Line)), nil
	}

	return nil, NewParserError(EXPRESSION_EXPECTED, p.curr().Lexeme, p.curr().Line)
}
