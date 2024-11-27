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
	if !p.IsAtEnd() {
		return p.Tokens[p.Idx]
	} else {
		return p.Curr()
	}
}

func (p *Parser) Advance() {
	if !p.IsAtEnd() {
		p.Idx++
	}
}

func (p *Parser) IsAtEnd() bool {
	return p.Idx >= len(p.Tokens)
}

func (p *Parser) IsSameLine() bool {
	return p.Curr().Line == p.Peek().Line
}

func (p *Parser) Check(expected tokens.TokenType) bool {
	if p.IsAtEnd() {
		return false
	}

	return p.Peek().Type == expected
}

func (p *Parser) MatchAndAdvance(expectedTypes ...tokens.TokenType) bool {
	nextTkn := p.Peek()

	for _, expected := range expectedTypes {
		if nextTkn.Type == expected {
			p.Advance()
			return true
		}
	}

	return false
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
	return p.EqualityRule()
}

func (p *Parser) EqualityRule() (*ast.AstNode, *ParserError) {
	node, err := p.ComparisonRule()
	if err != nil {
		return nil, err
	}

	if node != nil {
		token := p.Peek()

		if p.MatchAndAdvance(tokens.EQUAL_EQUAL, tokens.BANG_EQUAL) {
			if !p.IsSameLine() {
				return nil, NewParserError("expected expression", p.Curr().Lexeme, p.Curr().Line)
			}

			p.Advance()

			right, err := p.EqualityRule()
			if err != nil {
				return nil, err
			}

			if right == nil {
				return nil, NewParserError("expected expression", p.Curr().Lexeme, p.Curr().Line)
			} else {
				if right.Expr.ParseExpr() == "null" {
					return nil, NewParserError("expected expression", p.Curr().Lexeme, p.Curr().Line)
				}

				node = ast.NewAstNode(ast.BINARY, ast.NewBinaryExpr(node.Expr, token.Type, right.Expr, p.Curr().Line))
			}
		}
	}

	return node, nil
}

func (p *Parser) ComparisonRule() (*ast.AstNode, *ParserError) {
	node, err := p.AdditiveRule()
	if err != nil {
		return nil, err
	}

	if node != nil {
		token := p.Peek()

		if p.MatchAndAdvance(tokens.LESS, tokens.LESS_EQUAL, tokens.GREATER, tokens.GREATER_EQUAL) {
			if !p.IsSameLine() {
				return nil, NewParserError("expected expression", p.Curr().Lexeme, p.Curr().Line)
			}

			p.Advance()

			right, err := p.ComparisonRule()
			if err != nil {
				return nil, err
			}

			if right == nil {
				return nil, NewParserError("expected expression", p.Curr().Lexeme, p.Curr().Line)
			} else {
				if right.Expr.ParseExpr() == "null" {
					return nil, NewParserError("expected expression", p.Curr().Lexeme, p.Curr().Line)
				}

				node = ast.NewAstNode(ast.BINARY, ast.NewBinaryExpr(node.Expr, token.Type, right.Expr, p.Curr().Line))
			}
		}
	}

	return node, nil
}

func (p *Parser) AdditiveRule() (*ast.AstNode, *ParserError) {
	node, err := p.MultiplicativeRule()
	if err != nil {
		return nil, err
	}

	if node != nil {
		token := p.Peek()

		if p.MatchAndAdvance(tokens.PLUS, tokens.MINUS) {
			if !p.IsSameLine() {
				return nil, NewParserError("expected expression", p.Curr().Lexeme, p.Curr().Line)
			}

			p.Advance()

			right, err := p.AdditiveRule()
			if err != nil {
				return nil, err
			}

			if right == nil {
				return nil, NewParserError("expected expression", p.Curr().Lexeme, p.Curr().Line)
			} else {
				if right.Expr.ParseExpr() == "null" {
					return nil, NewParserError("expected expression", p.Curr().Lexeme, p.Curr().Line)
				}

				node = ast.NewAstNode(ast.BINARY, ast.NewBinaryExpr(node.Expr, token.Type, right.Expr, p.Curr().Line))
			}
		}
	}

	return node, nil
}

func (p *Parser) MultiplicativeRule() (*ast.AstNode, *ParserError) {
	node, err := p.UnaryRule()
	if err != nil {
		return nil, err
	}

	if node != nil {
		token := p.Peek()

		if p.MatchAndAdvance(tokens.STAR, tokens.SLASH) {
			if !p.IsSameLine() {
				return nil, NewParserError("expected expression", p.Curr().Lexeme, p.Curr().Line)
			}

			p.Advance()

			right, err := p.MultiplicativeRule()
			if err != nil {
				return nil, err
			}

			if right == nil {
				return nil, NewParserError("expected expression", p.Curr().Lexeme, p.Curr().Line)
			} else {
				if right.Expr.ParseExpr() == "null" {
					return nil, NewParserError("expected expression", p.Curr().Lexeme, p.Curr().Line)
				}

				node = ast.NewAstNode(ast.BINARY, ast.NewBinaryExpr(node.Expr, token.Type, right.Expr, p.Curr().Line))
			}

		}
	}

	return node, nil
}

func (p *Parser) UnaryRule() (*ast.AstNode, *ParserError) {
	expectedOperators := []tokens.TokenType{tokens.BANG, tokens.MINUS}

	if utils.HasValueArray(expectedOperators, p.Curr().Type) {
		node, err := p.Parse()
		if err != nil {
			return nil, err
		}

		if node != nil {
			return ast.NewAstNode(ast.UNARY, ast.NewUnaryExpr(p.Curr().Type, node.Expr, p.Curr().Line)), nil
		}

		return nil, NewParserError("expected expression", p.Curr().Lexeme, p.Curr().Line)
	}

	return p.PrimaryRule()
}

func (p *Parser) PrimaryRule() (*ast.AstNode, *ParserError) {
	canIgnore := []tokens.TokenType{tokens.EOF, tokens.ILLEGAL, tokens.IGNORE}

	if utils.HasValueArray(canIgnore, p.Curr().Type) {
		return nil, nil
	}

	if p.Curr().Type == tokens.LEFT_PAREN {
		node, err := p.Parse()
		if err != nil {
			return nil, err
		}

		if node != nil {
			err := NewParserError("expected ')' after this expression", p.Curr().Lexeme, p.Curr().Line)
			p.Consume(tokens.RIGHT_PAREN, *err)
			return ast.NewAstNode(ast.GROUPING, ast.NewGroupingExpr(node.Expr, p.Curr().Line)), nil
		}

		return nil, NewParserError("expected expression", p.Curr().Lexeme, p.Curr().Line)
	}

	return ast.NewAstNode(ast.LITERAL, ast.NewLiteralExpr(p.Curr().Type, p.Curr().Literal, p.Curr().Line)), nil
}
