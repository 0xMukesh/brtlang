package ast

type Scope int

const (
	GLOBAL Scope = iota
	LOCAL
)

type Stmt interface {
	ParseStmt()
	GetExpr() Expr
	GetLine() int
}

type BaseStmt struct {
	Line int
}

func (s BaseStmt) ParseStmt()       {}
func (s BaseStmt) GetLine() int     { return s.Line }
func (s BaseStmt) isAstValue() bool { return true }

type VarAssignStmt struct {
	BaseStmt
	Expr Expr
	Name string
}

func (s VarAssignStmt) GetExpr() Expr { return s.Expr }
func NewVarAssignStmt(name string, expr Expr, line int) VarAssignStmt {
	return VarAssignStmt{
		Name: name,
		Expr: expr,
		BaseStmt: BaseStmt{
			Line: line,
		},
	}
}

type PrintStmt struct {
	BaseStmt
	Expr Expr
}

func (s PrintStmt) GetExpr() Expr { return s.Expr }
func NewPrintStmt(expr Expr, line int) PrintStmt {
	return PrintStmt{
		Expr: expr,
		BaseStmt: BaseStmt{
			Line: line,
		},
	}
}

type CreateBlockStmt struct {
	BaseStmt
	Nodes []AstNode
}

func (s CreateBlockStmt) GetExpr() Expr { return nil }
func NewCreateBlockStmt(nodes []AstNode, line int) CreateBlockStmt {
	return CreateBlockStmt{
		Nodes: nodes,
		BaseStmt: BaseStmt{
			Line: line,
		},
	}
}

type CloseBlockStmt struct {
	BaseStmt
}

func (s CloseBlockStmt) GetExpr() Expr { return nil }
func NewCloseBlockStmt(line int) CloseBlockStmt {
	return CloseBlockStmt{
		BaseStmt: BaseStmt{
			Line: line,
		},
	}
}
