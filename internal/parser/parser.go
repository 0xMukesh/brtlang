package parser

import (
	"fmt"

	"github.com/0xmukesh/interpreter/internal/ast"
	"github.com/0xmukesh/interpreter/internal/tokens"
	"github.com/0xmukesh/interpreter/internal/utils"
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

func (p *Parser) IsAtEnd() bool {
	return p.Idx >= len(p.Tokens)
}

func (p *Parser) Check(expected tokens.TokenType) bool {
	if p.IsAtEnd() {
		return false
	}

	return p.Peek().Type == expected
}

func (p *Parser) Consume(expected tokens.TokenType, err string) {
	if p.Check(expected) {
		p.Idx++
		return
	}

	utils.EPrint(fmt.Sprintf("%s\n", err))
}

// rules in order of precedence (lowest to highest)
// 1. equality
// 2. comparison
// 3. term
// 4. factor
// 5. unary
// 6. primary

func (p *Parser) Parse() *ast.AstNode {
	if p.IsAtEnd() {
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

	if curr.Type == tokens.LEFT_PAREN {
		node := p.Parse()
		if node != nil {
			p.Consume(tokens.RIGHT_PAREN, "expected ')' after this expression")
			return ast.NewAstNode(ast.GROUPING, ast.NewGroupingExpr(node.Expr))
		}
	}

	return ast.NewAstNode(ast.LITERAL, ast.NewLiteralExpr(curr.Type, curr.Lexeme))
}
