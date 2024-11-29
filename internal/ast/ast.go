package ast

type AstNodeType int

const (
	EXPR AstNodeType = iota
	STMT
)

type AstNodeValue interface {
	isAstValue() bool
}

type AstNode struct {
	Type  AstNodeType
	Value AstNodeValue
}

func (n AstNode) ExtractExpr() Expr {
	value, ok := n.Value.(Expr)
	if !ok {
		return nil
	}

	return value
}

func (n AstNode) ExtractStmtExpr() *Expr {
	value, ok := n.Value.(Stmt)
	if !ok {
		return nil
	}

	expr := value.GetExpr()

	return &expr
}

type Ast []AstNode

func NewAstNode(nodeType AstNodeType, value AstNodeValue) *AstNode {
	return &AstNode{
		Type:  nodeType,
		Value: value,
	}
}
