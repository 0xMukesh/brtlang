package ast

type Ast []AstNode

type AstNode struct {
	Type ExprType
	Expr Expr
}

func NewAstNode(exprType ExprType, expr Expr) *AstNode {
	return &AstNode{
		Type: exprType,
		Expr: expr,
	}
}
