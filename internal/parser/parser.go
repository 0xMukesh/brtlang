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

func (p *Parser) Consume(expected tokens.TokenType, err ParserError) {
	if p.Check(expected) {
		p.Idx++
		return
	}

	utils.EPrint(fmt.Sprintf("%s\n", err.String()))
}

func (p *Parser) Parse() (*ast.AstNode, *ParserError) {
	if p.IsAtEnd() {
		return nil, nil
	}

	p.Idx++
	return p.equality()
}

func (p *Parser) equality() (*ast.AstNode, *ParserError) {
	return p.comparison()
}

func (p *Parser) comparison() (*ast.AstNode, *ParserError) {
	return p.term()
}

func (p *Parser) term() (*ast.AstNode, *ParserError) {
	return p.factor()
}

func (p *Parser) factor() (*ast.AstNode, *ParserError) {
	return p.unary()
}

func (p *Parser) unary() (*ast.AstNode, *ParserError) {
	operator := p.Curr()
	expectedOperators := []tokens.TokenType{tokens.BANG, tokens.MINUS}

	if utils.HasValueArray(expectedOperators, operator.Type) {
		node, err := p.Parse()
		if err != nil {
			return nil, err
		}

		if node != nil {
			return ast.NewAstNode(ast.UNARY, ast.NewUnaryExpr(operator.Type, node.Expr)), nil
		}

		return nil, NewParserError(EXPRESSION_EXPECTED, "expected expression")
	}

	return p.primary()
}

func (p *Parser) primary() (*ast.AstNode, *ParserError) {
	curr := p.Curr()
	canIgnore := []tokens.TokenType{tokens.EOF, tokens.ILLEGAL, tokens.IGNORE}

	if utils.HasValueArray(canIgnore, curr.Type) {
		return nil, nil
	}

	if curr.Type == tokens.LEFT_PAREN {
		node, err := p.Parse()
		if err != nil {
			return nil, err
		}

		if node != nil {
			err := NewParserError(TOKEN_EXPECTED, "expected ')' after this expression")
			p.Consume(tokens.RIGHT_PAREN, *err)
			return ast.NewAstNode(ast.GROUPING, ast.NewGroupingExpr(node.Expr)), nil
		}

		return nil, NewParserError(EXPRESSION_EXPECTED, "expected expression")
	}

	return ast.NewAstNode(ast.LITERAL, ast.NewLiteralExpr(curr.Type, curr.Lexeme)), nil
}
