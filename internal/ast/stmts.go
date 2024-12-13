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
	Node AstNode
	Name string
}

func (s VarAssignStmt) GetExpr() Expr { return s.Node.ExtractExpr() }
func NewVarAssignStmt(name string, node AstNode, line int) VarAssignStmt {
	return VarAssignStmt{
		Name: name,
		Node: node,
		BaseStmt: BaseStmt{
			Line: line,
		},
	}
}

type VarReassignStmt struct {
	BaseStmt
	Node AstNode
	Name string
}

func (s VarReassignStmt) GetExpr() Expr { return s.Node.ExtractExpr() }
func NewVarReassignStmt(name string, node AstNode, line int) VarReassignStmt {
	return VarReassignStmt{
		Name: name,
		Node: node,
		BaseStmt: BaseStmt{
			Line: line,
		},
	}
}

type PrintStmt struct {
	BaseStmt
	Node AstNode
}

func (s PrintStmt) GetExpr() Expr { return s.Node.ExtractExpr() }
func NewPrintStmt(node AstNode, line int) PrintStmt {
	return PrintStmt{
		Node: node,
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
	Node           AstNode
	IfBranch       AstNode
	ElseIfBranches *[]ElseIfStmt
	ElseBranch     *ElseStmt
}

func (s IfStmt) GetExpr() Expr { return s.Node.ExtractExpr() }
func NewIfStmt(node AstNode, ifBranch AstNode, elseIfBranches *[]ElseIfStmt, elseBranch *ElseStmt, line int) IfStmt {
	return IfStmt{
		Node:           node,
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
	Node   AstNode
	Branch AstNode
}

func (s ElseIfStmt) GetExpr() Expr { return s.Node.ExtractExpr() }
func NewElseIfStmt(node AstNode, branch AstNode, line int) ElseIfStmt {
	return ElseIfStmt{
		Node:   node,
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

type WhileStmt struct {
	BaseStmt
	Node   AstNode
	Branch AstNode
}

func (s WhileStmt) GetExpr() Expr { return s.Node.ExtractExpr() }
func NewWhileStmt(node AstNode, branch AstNode, line int) WhileStmt {
	return WhileStmt{
		Node:   node,
		Branch: branch,
		BaseStmt: BaseStmt{
			Line: line,
		},
	}
}

type FuncDeclarationStmt struct {
	BaseStmt
	Name string
	Node AstNode
	Args []AstNode
}

func (s FuncDeclarationStmt) GetExpr() Expr { return nil }
func NewFuncDeclarationStmt(name string, args []AstNode, node AstNode, line int) FuncDeclarationStmt {
	return FuncDeclarationStmt{
		Name: name,
		Node: node,
		Args: args,
		BaseStmt: BaseStmt{
			Line: line,
		},
	}
}

type FuncCallStmt struct {
	BaseStmt
	Name       string
	Args       []AstNode
	ReturnNode AstNode
}

func (s FuncCallStmt) GetExpr() Expr { return s.ReturnNode.ExtractExpr() }
func NewFuncCallStmt(name string, args []AstNode, returnNode AstNode, line int) FuncCallStmt {
	return FuncCallStmt{
		Name:       name,
		Args:       args,
		ReturnNode: returnNode,
		BaseStmt: BaseStmt{
			Line: line,
		},
	}
}

type ReturnStmt struct {
	BaseStmt
	Node AstNode
}

func (s ReturnStmt) GetExpr() Expr {
	return s.Node.ExtractExpr()
}
func NewReturnStmt(node AstNode, line int) ReturnStmt {
	return ReturnStmt{
		Node: node,
		BaseStmt: BaseStmt{
			Line: line,
		},
	}
}
