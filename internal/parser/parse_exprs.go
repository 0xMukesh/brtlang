package parser

import (
	"github.com/0xmukesh/interpreter/internal/ast"
	"github.com/0xmukesh/interpreter/internal/tokens"
)

func (p *Parser) parseGroupingExpr() (*ast.AstNode, *ParserError) {
	node, err := p.Parse()
	if err != nil {
		return nil, err
	}

	if node != nil {
		err := NewParserError(MISSING_RPAREN, p.curr().Lexeme, p.curr().Line)
		p.consume(tokens.RIGHT_PAREN, *err)

		expr, err := p.extractExpr(*node)
		if err != nil {
			return nil, err
		}

		return ast.NewAstNode(ast.EXPR, ast.NewGroupingExpr(expr, p.curr().Line)), nil
	}

	return nil, NewParserError(EXPRESSION_EXPECTED, p.curr().Lexeme, p.curr().Line)
}
