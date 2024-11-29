package ast

type Stmt interface {
	ParseStmt()
	GetExpr() Expr
	GetLine() int
}

type BaseStmt struct {
	Line int
	Expr Expr
}

func (s BaseStmt) ParseStmt() {}

func (s BaseStmt) GetExpr() Expr    { return s.Expr }
func (s BaseStmt) GetLine() int     { return s.Line }
func (s BaseStmt) isAstValue() bool { return true }

type VarAssignStmt struct {
	BaseStmt
	Name string
}

func NewVarAssignStmt(name string, expr Expr, line int) VarAssignStmt {
	return VarAssignStmt{
		Name: name,
		BaseStmt: BaseStmt{
			Expr: expr,
			Line: line,
		},
	}
}

type PrintStmt struct {
	BaseStmt
}

func NewPrintStmt(expr Expr, line int) PrintStmt {
	return PrintStmt{
		BaseStmt: BaseStmt{
			Expr: expr,
			Line: line,
		},
	}
}
