package ast

import (
	"fmt"

	"github.com/0xmukesh/interpreter/internal/tokens"
)

type ExprType int

const (
	LITERAL ExprType = iota
	GROUPING
	UNARY
	BINARY
)

type Expr interface {
	ParseExpr() string
	GetLine() int
}

type BaseExpr struct {
	Line int
}

func (e BaseExpr) GetLine() int {
	return e.Line
}

type LiteralExpr struct {
	BaseExpr
	TokenType tokens.TokenType
	Value     string
}

func (e LiteralExpr) ParseExpr() string {
	return e.Value
}
func NewLiteralExpr(tokenType tokens.TokenType, value string, line int) LiteralExpr {
	return LiteralExpr{
		TokenType: tokenType,
		Value:     value,
		BaseExpr: BaseExpr{
			Line: line,
		},
	}
}

type GroupingExpr struct {
	BaseExpr
	Expr Expr
}

func (e GroupingExpr) ParseExpr() string {
	return fmt.Sprintf("(group %s)", e.Expr.ParseExpr())
}
func NewGroupingExpr(expr Expr, line int) GroupingExpr {
	return GroupingExpr{
		Expr: expr,
		BaseExpr: BaseExpr{
			Line: line,
		},
	}
}

type UnaryExpr struct {
	BaseExpr
	Operator tokens.TokenType
	Expr     Expr
}

func (e UnaryExpr) ParseExpr() string {
	return fmt.Sprintf("(%s %s)", e.Operator.Literal(), e.Expr.ParseExpr())
}
func NewUnaryExpr(operator tokens.TokenType, expr Expr, line int) UnaryExpr {
	return UnaryExpr{
		Operator: operator,
		Expr:     expr,
		BaseExpr: BaseExpr{
			Line: line,
		},
	}
}

type BinaryExpr struct {
	BaseExpr
	Left     Expr
	Operator tokens.TokenType
	Right    Expr
}

func (e BinaryExpr) ParseExpr() string {
	return fmt.Sprintf("(%s %s %s)", e.Operator.Literal(), e.Left.ParseExpr(), e.Right.ParseExpr())
}
func NewBinaryExpr(left Expr, operator tokens.TokenType, right Expr, line int) BinaryExpr {
	return BinaryExpr{
		Left:     left,
		Operator: operator,
		Right:    right,
		BaseExpr: BaseExpr{
			Line: line,
		},
	}
}
