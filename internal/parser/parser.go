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

func (p *Parser) Advance() {
	if !p.IsAtEnd() {
		p.Idx++
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
	expr, err := p.factor()
	if err != nil {
		return nil, err
	}

	if expr != nil {
		// "peeks" and checks if the next token is the expected operator. if yes, then skip onto the part after the operator
		token := p.Peek()
		expectedOperators := []tokens.TokenType{tokens.PLUS, tokens.MINUS}

		if utils.HasValueArray(expectedOperators, token.Type) {
			p.Advance()
			p.Advance()

			right, err := p.factor()
			if err != nil {
				return nil, err
			}

			if right != nil {
				expr = ast.NewAstNode(ast.BINARY, ast.NewBinaryExpr(expr.Expr, token.Type, right.Expr))
			}
		}
	}

	return expr, nil
}

func (p *Parser) factor() (*ast.AstNode, *ParserError) {
	expr, err := p.unary()
	if err != nil {
		return nil, err
	}

	if expr != nil {
		// "peeks" and checks if the next token is the expected operator. if yes, then skip onto the part after the operator
		token := p.Peek()
		expectedOperators := []tokens.TokenType{tokens.STAR, tokens.SLASH}

		if utils.HasValueArray(expectedOperators, token.Type) {
			p.Advance()
			p.Advance()

			right, err := p.unary()
			if err != nil {
				return nil, err
			}

			if right != nil {
				expr = ast.NewAstNode(ast.BINARY, ast.NewBinaryExpr(expr.Expr, token.Type, right.Expr))
			}
		}
	}

	return expr, nil
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
