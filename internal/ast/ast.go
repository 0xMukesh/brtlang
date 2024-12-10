package ast

type AstNodeType int

const (
	EXPR AstNodeType = iota
	STMT
)

type AstNodeValue interface {
	IsExpr() bool
	IsStmt() bool
	isAstValue() bool
}

type AstNode struct {
	Type  AstNodeType
	Value AstNodeValue
}

func (n AstNode) ExtractExpr() Expr {
	var expr Expr

	switch v := n.Value.(type) {
	case Expr:
		value, ok := n.Value.(Expr)

		if !ok {
			value = nil
		}

		expr = value
	case Stmt:
		expr = v.GetExpr()
	}

	return expr
}

type Ast []AstNode

func NewAstNode(nodeType AstNodeType, value AstNodeValue) *AstNode {
	return &AstNode{
		Type:  nodeType,
		Value: value,
	}
}
