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

func (p *Parser) curr() tokens.Token {
	return p.Tokens[p.Idx-1]
}

func (p *Parser) peek() tokens.Token {
	if !p.isAtEnd() {
		return p.Tokens[p.Idx]
	} else {
		return p.curr()
	}
}

func (p *Parser) advance() {
	if !p.isAtEnd() {
		p.Idx++
	}
}

func (p *Parser) isAtEnd() bool {
	return p.Idx >= len(p.Tokens)
}

func (p *Parser) isSameLine() bool {
	return p.curr().Line == p.peek().Line
}

func (p *Parser) check(expected tokens.TokenType) bool {
	if p.isAtEnd() {
		return false
	}

	return p.peek().Type == expected
}

func (p *Parser) matchAndAdvance(expectedTypes ...tokens.TokenType) bool {
	nextTkn := p.peek()

	for _, expected := range expectedTypes {
		if nextTkn.Type == expected {
			p.advance()
			return true
		}
	}

	return false
}

func (p *Parser) consume(expected tokens.TokenType, err ParserError) {
	if p.check(expected) {
		p.Idx++
		return
	}

	utils.EPrint(fmt.Sprintf("%s\n", err.Error()))
}

func (p *Parser) Parse() (*ast.AstNode, *ParserError) {
	if p.isAtEnd() {
		return nil, nil
	}

	p.Idx++
	return p.equalityRule()
}

func (p *Parser) BuildAst() (ast.Ast, *ParserError) {
	var ast ast.Ast

	for {
		node, err := p.Parse()
		if err != nil {
			return nil, err
		}

		if node != nil {
			ast = append(ast, *node)
		}

		if p.isAtEnd() {
			return ast, nil
		}
	}
}

func (p *Parser) equalityRule() (*ast.AstNode, *ParserError) {
	node, err := p.comparisonRule()
	if err != nil {
		return nil, err
	}

	if node != nil {
		token := p.peek()

		if p.matchAndAdvance(tokens.EQUAL_EQUAL, tokens.BANG_EQUAL) {
			if !p.isSameLine() {
				return nil, NewParserError("expected expression", p.curr().Lexeme, p.curr().Line)
			}

			p.advance()

			right, err := p.equalityRule()
			if err != nil {
				return nil, err
			}

			if right == nil {
				return nil, NewParserError("expected expression", p.curr().Lexeme, p.curr().Line)
			} else {
				if right.Expr.ParseExpr() == "null" {
					return nil, NewParserError("expected expression", p.curr().Lexeme, p.curr().Line)
				}

				node = ast.NewAstNode(ast.BINARY, ast.NewBinaryExpr(node.Expr, token.Type, right.Expr, p.curr().Line))
			}
		}
	}

	return node, nil
}

func (p *Parser) comparisonRule() (*ast.AstNode, *ParserError) {
	node, err := p.additiveRule()
	if err != nil {
		return nil, err
	}

	if node != nil {
		token := p.peek()

		if p.matchAndAdvance(tokens.LESS, tokens.LESS_EQUAL, tokens.GREATER, tokens.GREATER_EQUAL) {
			if !p.isSameLine() {
				return nil, NewParserError("expected expression", p.curr().Lexeme, p.curr().Line)
			}

			p.advance()

			right, err := p.comparisonRule()
			if err != nil {
				return nil, err
			}

			if right == nil {
				return nil, NewParserError("expected expression", p.curr().Lexeme, p.curr().Line)
			} else {
				if right.Expr.ParseExpr() == "null" {
					return nil, NewParserError("expected expression", p.curr().Lexeme, p.curr().Line)
				}

				node = ast.NewAstNode(ast.BINARY, ast.NewBinaryExpr(node.Expr, token.Type, right.Expr, p.curr().Line))
			}
		}
	}

	return node, nil
}

func (p *Parser) additiveRule() (*ast.AstNode, *ParserError) {
	node, err := p.multiplicativeRule()
	if err != nil {
		return nil, err
	}

	if node != nil {
		token := p.peek()

		if p.matchAndAdvance(tokens.PLUS, tokens.MINUS) {
			if !p.isSameLine() {
				return nil, NewParserError("expected expression", p.curr().Lexeme, p.curr().Line)
			}

			p.advance()

			right, err := p.additiveRule()
			if err != nil {
				return nil, err
			}

			if right == nil {
				return nil, NewParserError("expected expression", p.curr().Lexeme, p.curr().Line)
			} else {
				if right.Expr.ParseExpr() == "null" {
					return nil, NewParserError("expected expression", p.curr().Lexeme, p.curr().Line)
				}

				node = ast.NewAstNode(ast.BINARY, ast.NewBinaryExpr(node.Expr, token.Type, right.Expr, p.curr().Line))
			}
		}
	}

	return node, nil
}

func (p *Parser) multiplicativeRule() (*ast.AstNode, *ParserError) {
	node, err := p.unaryRule()
	if err != nil {
		return nil, err
	}

	if node != nil {
		token := p.peek()

		if p.matchAndAdvance(tokens.STAR, tokens.SLASH) {
			if !p.isSameLine() {
				return nil, NewParserError("expected expression", p.curr().Lexeme, p.curr().Line)
			}

			p.advance()

			right, err := p.multiplicativeRule()
			if err != nil {
				return nil, err
			}

			if right == nil {
				return nil, NewParserError("expected expression", p.curr().Lexeme, p.curr().Line)
			} else {
				if right.Expr.ParseExpr() == "null" {
					return nil, NewParserError("expected expression", p.curr().Lexeme, p.curr().Line)
				}

				node = ast.NewAstNode(ast.BINARY, ast.NewBinaryExpr(node.Expr, token.Type, right.Expr, p.curr().Line))
			}

		}
	}

	return node, nil
}

func (p *Parser) unaryRule() (*ast.AstNode, *ParserError) {
	expectedOperators := []tokens.TokenType{tokens.BANG, tokens.MINUS}
	operator := p.curr()

	if utils.HasValueArray(expectedOperators, operator.Type) {
		node, err := p.Parse()
		if err != nil {
			return nil, err
		}

		if node != nil {
			return ast.NewAstNode(ast.UNARY, ast.NewUnaryExpr(operator.Type, node.Expr, p.curr().Line)), nil
		}

		return nil, NewParserError("expected expression", p.curr().Lexeme, p.curr().Line)
	}

	return p.primaryRule()
}

func (p *Parser) primaryRule() (*ast.AstNode, *ParserError) {
	canIgnore := []tokens.TokenType{tokens.EOF, tokens.ILLEGAL, tokens.IGNORE}

	if utils.HasValueArray(canIgnore, p.curr().Type) {
		return nil, nil
	}

	if p.curr().Type == tokens.LEFT_PAREN {
		node, err := p.Parse()
		if err != nil {
			return nil, err
		}

		if node != nil {
			err := NewParserError("expected ')' after this expression", p.curr().Lexeme, p.curr().Line)
			p.consume(tokens.RIGHT_PAREN, *err)
			return ast.NewAstNode(ast.GROUPING, ast.NewGroupingExpr(node.Expr, p.curr().Line)), nil
		}

		return nil, NewParserError("expected expression", p.curr().Lexeme, p.curr().Line)
	}

	return ast.NewAstNode(ast.LITERAL, ast.NewLiteralExpr(p.curr().Type, p.curr().Literal, p.curr().Line)), nil
}
