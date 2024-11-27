package parser

import (
	"github.com/0xmukesh/interpreter/internal/ast"
	"github.com/0xmukesh/interpreter/internal/tokens"
)

type Parser struct {
	Tokens []tokens.Token
	Idx    int
}

func NewParser(tokens []tokens.Token) *Parser {
	return &Parser{
		Tokens: tokens,
		Idx:    0,
	}
}

func (p *Parser) Curr() tokens.Token {
	return p.Tokens[p.Idx-1]
}

func (p *Parser) Peek() tokens.Token {
	return p.Tokens[p.Idx]
}

// rules in order of precedence (lowest to highest)
// 1. equality
// 2. comparison
// 3. term
// 4. factor
// 5. unary
// 6. primary

func (p *Parser) Parse() *ast.AstNode {
	if p.Idx >= len(p.Tokens) {
		return nil
	}

	p.Idx++
	return p.equality()
}

func (p *Parser) equality() *ast.AstNode {
	return p.comparison()
}

func (p *Parser) comparison() *ast.AstNode {
	return p.term()
}

func (p *Parser) term() *ast.AstNode {
	return p.factor()
}

func (p *Parser) factor() *ast.AstNode {
	return p.unary()
}

func (p *Parser) unary() *ast.AstNode {
	return p.primary()
}

func (p *Parser) primary() *ast.AstNode {
	curr := p.Curr()

	if curr.Type == tokens.EOF || curr.Type == tokens.ILLEGAL || curr.Type == tokens.IGNORE {
		return nil
	}

	return ast.NewAstNode(ast.LITERAL, ast.NewLiteralExpr(curr.Type, curr.Lexeme))
}
