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
}

type LiteralExpr struct {
	TokenType tokens.TokenType
	Value     string
}

func (e LiteralExpr) ParseExpr() string {
	return e.Value
}
func NewLiteralExpr(tokenType tokens.TokenType, value string) LiteralExpr {
	return LiteralExpr{
		TokenType: tokenType,
		Value:     value,
	}
}

type GroupingExpr struct {
	Expr Expr
}

func (e GroupingExpr) ParseExpr() string {
	return fmt.Sprintf("(group %s)", e.Expr.ParseExpr())
}
func NewGroupingExpr(expr Expr) GroupingExpr {
	return GroupingExpr{
		Expr: expr,
	}
}

type UnaryExpr struct {
	Operator tokens.TokenType
	Expr     Expr
}

func (e UnaryExpr) ParseExpr() string {
	return fmt.Sprintf("(%s %s)", e.Operator.Literal(), e.Expr.ParseExpr())
}
func NewUnaryExpr(operator tokens.TokenType, expr Expr) UnaryExpr {
	return UnaryExpr{
		Operator: operator,
		Expr:     expr,
	}
}
