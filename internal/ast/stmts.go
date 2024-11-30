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
func (s BaseStmt) IsExpr() bool     { return false }
func (s BaseStmt) IsStmt() bool     { return true }
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

type IfStmt struct {
	BaseStmt
	Expr           Expr
	IfBranch       AstNode
	ElseIfBranches *[]ElseIfStmt
	ElseBranch     *ElseStmt
}

func (s IfStmt) GetExpr() Expr { return s.Expr }
func NewIfStmt(expr Expr, ifBranch AstNode, elseIfBranches *[]ElseIfStmt, elseBranch *ElseStmt, line int) IfStmt {
	return IfStmt{
		Expr:           expr,
		IfBranch:       ifBranch,
		ElseIfBranches: elseIfBranches,
		ElseBranch:     elseBranch,
		BaseStmt: BaseStmt{
			Line: line,
		},
	}
}

type ElseIfStmt struct {
	BaseStmt
	Expr   Expr
	Branch AstNode
}

func (s ElseIfStmt) GetExpr() Expr { return s.Expr }
func NewElseIfStmt(expr Expr, branch AstNode, line int) ElseIfStmt {
	return ElseIfStmt{
		Expr:   expr,
		Branch: branch,
		BaseStmt: BaseStmt{
			Line: line,
		},
	}
}

type ElseStmt struct {
	BaseStmt
	Branch AstNode
}

func (s ElseStmt) GetExpr() Expr { return nil }
func NewElseStmt(branch AstNode, line int) ElseStmt {
	return ElseStmt{
		Branch: branch,
		BaseStmt: BaseStmt{
			Line: line,
		},
	}
}
