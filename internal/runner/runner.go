package runner

import (
	"fmt"

	"github.com/0xmukesh/interpreter/internal/ast"
	"github.com/0xmukesh/interpreter/internal/evaluator"
	"github.com/0xmukesh/interpreter/internal/runtime"
	"github.com/0xmukesh/interpreter/internal/utils"
)

type Runner struct {
	Ast       ast.Ast
	Env       *runtime.Environment
	Evaluator *evaluator.Evaluator
	Idx       int
}

func NewRunner(ast ast.Ast, env *runtime.Environment, evaluator *evaluator.Evaluator) *Runner {
	return &Runner{
		Ast:       ast,
		Env:       env,
		Evaluator: evaluator,
		Idx:       0,
	}
}

func (r *Runner) IsAtEnd() bool {
	return r.Idx >= len(r.Ast)
}

func (r *Runner) advance() {
	if !r.IsAtEnd() {
		r.Idx++
	}
}

func (r *Runner) curr() ast.AstNode {
	if !r.IsAtEnd() {
		return r.Ast[r.Idx]
	} else {
		return r.Ast[len(r.Ast)-1]
	}
}

func (r *Runner) Run() {
	curr := r.curr()

	if !r.IsAtEnd() {
		expr, isExpr := curr.Value.(ast.Expr)
		stmt, _ := curr.Value.(ast.Stmt)

		if !isExpr {
			_, isPrintStmt := curr.Value.(ast.PrintStmt)
			varAssignStmt, isVarAssignStmt := curr.Value.(ast.VarAssignStmt)

			if isPrintStmt {
				val, err := r.Evaluator.EvaluateExpr(stmt.GetExpr())
				if err != nil {
					utils.EPrint(fmt.Sprintf("%s\n", err.Error()))
				}

				fmt.Println(val)
			} else if isVarAssignStmt {
				val, err := r.Evaluator.EvaluateExpr(varAssignStmt.Expr)
				if err != nil {
					utils.EPrint(fmt.Sprintf("%s\n", err.Error()))
				}

				if val != nil {
					r.Env.Vars[varAssignStmt.Name] = *runtime.NewRuntimeValue(val.Value)
				}
			}
		} else {
			_, err := r.Evaluator.EvaluateExpr(expr)
			if err != nil {
				utils.EPrint(fmt.Sprintf("%s\n", err.Error()))
			}
		}

		r.advance()
	}
}
