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
		p.Idx++
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

	if p.curr().Type == tokens.LEFT_PAREN {
		node, err := p.Parse()
		if err != nil {
			return nil, err
		}

		if node != nil {
			err := NewParserError("expected ')' after this expression", p.curr().Lexeme, p.curr().Line)
			p.consume(tokens.RIGHT_PAREN, *err)

			expr, err := p.extractExpr(*node)
			if err != nil {
				return nil, err
			}

			return ast.NewAstNode(ast.EXPR, ast.NewGroupingExpr(expr, p.curr().Line)), nil
		}

		return nil, NewParserError("expected expression", p.curr().Lexeme, p.curr().Line)
	}

	if p.curr().Type.IsReserved() {
		if p.curr().Type == tokens.PRINT {
			node, err := p.Parse()
			if err != nil {
				return nil, err
			}

			if node != nil {
				err := NewParserError("expected ';' after print statement", p.curr().Lexeme, p.curr().Line)
				p.consume(tokens.SEMICOLON, *err)

				expr, err := p.extractExpr(*node)
				if err != nil {
					return nil, err
				}

				return ast.NewAstNode(ast.STMT, ast.NewPrintStmt(expr, p.curr().Line)), nil
			} else {
				return nil, NewParserError("expected expression after print statement", p.curr().Lexeme, p.curr().Line)
			}
		} else if p.curr().Type == tokens.VAR {
			varNameNode, err := p.Parse()
			if err != nil {
				return nil, err
			}
			if varNameNode == nil {
				return nil, NewParserError("missing variable name after `var`", p.curr().Lexeme, p.curr().Line)
			}

			varNameExpr := varNameNode.ExtractExpr()
			if varNameExpr == nil {
				return nil, NewParserError("invalid variable name identifier", p.curr().Lexeme, p.curr().Line)
			}

			varNameLiteralExpr, isLiteralExpr := varNameExpr.(ast.LiteralExpr)
			if !isLiteralExpr {
				return nil, NewParserError("invalid variable name identifier", p.curr().Lexeme, p.curr().Line)
			}

			varName := varNameLiteralExpr.Value

			err = NewParserError("expected '=' after variable name identifier", p.curr().Lexeme, p.curr().Line)
			p.consume(tokens.EQUAL, *err)

			varValueNode, err := p.Parse()
			if err != nil {
				return nil, err
			}

			if varValueNode == nil {
				return nil, NewParserError("missing assignment value to the variable", p.curr().Lexeme, p.curr().Line)
			}

			varValueExpr := varValueNode.ExtractExpr()

			if varValueExpr == nil {
				return nil, NewParserError("invalid variable value assignment", p.curr().Lexeme, p.curr().Line)
			}

			err = NewParserError("missing ';' after a statement", p.curr().Lexeme, p.curr().Line)
			p.consume(tokens.SEMICOLON, *err)

			return ast.NewAstNode(ast.STMT, ast.NewVarAssignStmt(varName, varValueExpr, p.curr().Line)), nil
		}
	}

	literal := strings.TrimSuffix(strings.TrimPrefix(p.curr().Lexeme, `"`), `"`)

	return ast.NewAstNode(ast.EXPR, ast.NewLiteralExpr(p.curr().Type, literal, p.curr().Line)), nil
}
