package parser

import (
	"fmt"
	"strings"

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
		p.advance()
		return
	}

	utils.EPrint(fmt.Sprintf("%s\n", err.Error()))
}

func (p *Parser) extractExpr(node ast.AstNode) (ast.Expr, *ParserError) {
	expr := node.ExtractExpr()
	if expr == nil {
		return nil, NewParserError("expected expression", p.curr().Lexeme, p.curr().Line)
	}

	return expr, nil
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

func (p *Parser) binaryRuleBuilder(selfRule func() (*ast.AstNode, *ParserError), nextRule func() (*ast.AstNode, *ParserError), expectedTokens ...tokens.TokenType) (*ast.AstNode, *ParserError) {
	leftNode, err := nextRule()
	if err != nil {
		return nil, err
	}

	if leftNode != nil {
		token := p.peek()

		if p.matchAndAdvance(expectedTokens...) {
			if !p.isSameLine() {
				return nil, NewParserError("expected expression", p.curr().Lexeme, p.curr().Line)
			}

			p.advance()

			rightNode, err := selfRule()
			if err != nil {
				return nil, err
			}

			if rightNode == nil {
				return nil, NewParserError("expected expression", p.curr().Lexeme, p.curr().Line)
			} else {
				leftExpr, err := p.extractExpr(*leftNode)
				if err != nil {
					return nil, err
				}

				rightExpr, err := p.extractExpr(*rightNode)
				if err != nil {
					return nil, err
				}

				if rightExpr.ParseExpr() == "null" {
					return nil, NewParserError("expected expression", p.curr().Lexeme, p.curr().Line)
				}

				leftNode = ast.NewAstNode(ast.EXPR, ast.NewBinaryExpr(leftExpr, token.Type, rightExpr, p.curr().Line))
			}
		}
	}

	return leftNode, nil
}

func (p *Parser) equalityRule() (*ast.AstNode, *ParserError) {
	return p.binaryRuleBuilder(p.equalityRule, p.comparisonRule, tokens.EQUAL_EQUAL, tokens.BANG_EQUAL)
}

func (p *Parser) comparisonRule() (*ast.AstNode, *ParserError) {
	return p.binaryRuleBuilder(p.comparisonRule, p.subtractionRule, tokens.LESS, tokens.LESS_EQUAL, tokens.GREATER, tokens.GREATER_EQUAL)
}

func (p *Parser) subtractionRule() (*ast.AstNode, *ParserError) {
	return p.binaryRuleBuilder(p.subtractionRule, p.additionRule, tokens.MINUS)
}

func (p *Parser) additionRule() (*ast.AstNode, *ParserError) {
	return p.binaryRuleBuilder(p.additionRule, p.multiplicationRule, tokens.PLUS)
}

func (p *Parser) multiplicationRule() (*ast.AstNode, *ParserError) {
	return p.binaryRuleBuilder(p.multiplicationRule, p.divisionRule, tokens.STAR)
}

func (p *Parser) divisionRule() (*ast.AstNode, *ParserError) {
	return p.binaryRuleBuilder(p.divisionRule, p.unaryRule, tokens.SLASH)
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
			expr, err := p.extractExpr(*node)
			if err != nil {
				return nil, err
			}

			return ast.NewAstNode(ast.EXPR, ast.NewUnaryExpr(operator.Type, expr, p.curr().Line)), nil
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

	switch p.curr().Type {
	case tokens.LEFT_PAREN:
		return p.parseGroupingExpr()
	case tokens.LEFT_BRACE:
		return p.parseCreateBlockStmt()
	case tokens.PRINT:
		return p.parsePrintStmt()
	case tokens.VAR:
		return p.parseVarAssignStmt()
	}

	literal := p.curr().Lexeme

	if p.curr().Type == tokens.STRING {
		literal = strings.TrimSuffix(strings.TrimPrefix(p.curr().Lexeme, `"`), `"`)
	}

	return ast.NewAstNode(ast.EXPR, ast.NewLiteralExpr(p.curr().Type, literal, p.curr().Line)), nil
}
