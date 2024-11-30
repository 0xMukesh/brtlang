package runner

import (
	"fmt"

	"github.com/0xmukesh/interpreter/internal/ast"
	"github.com/0xmukesh/interpreter/internal/evaluator"
	"github.com/0xmukesh/interpreter/internal/runtime"
	"github.com/0xmukesh/interpreter/internal/tokens"
	"github.com/0xmukesh/interpreter/internal/utils"
)

type Runner struct {
	Ast       ast.Ast
	Runtime   *runtime.Runtime
	Evaluator *evaluator.Evaluator
	Idx       int
}

func NewRunner(ast ast.Ast, runtime *runtime.Runtime, evaluator *evaluator.Evaluator) *Runner {
	return &Runner{
		Ast:       ast,
		Runtime:   runtime,
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

func (r *Runner) EvalAndRunNode(expr ast.Expr, node ast.AstNode) bool {
	evaledCondition, err := r.Evaluator.EvaluateExpr(expr)
	if err != nil {
		err := runtime.NewRuntimeError(runtime.ExpectedExprErrBuilder("boolean"), "smth", expr.GetLine())
		utils.EPrint(err.Error())
	}

	if evaledCondition != nil {
		conditionVal := (*evaledCondition).Value
		if conditionVal != true && conditionVal != false {
			err := runtime.NewRuntimeError(runtime.ExpectedExprErrBuilder("boolean"), evaledCondition.String(), expr.GetLine())
			utils.EPrint(err.Error())
		}

		if conditionVal == true {
			r.RunNode(node)
			return true
		}
	}

	return false
}

func (r *Runner) RunNode(node ast.AstNode) {
	expr, isExpr := node.Value.(ast.Expr)

	if !isExpr {
		switch value := node.Value.(type) {
		case ast.PrintStmt:
			val, err := r.Evaluator.EvaluateExpr(value.GetExpr())
			if err != nil {
				utils.EPrint(fmt.Sprintf("%s\n", err.Error()))
			}

			fmt.Println(val)
		case ast.CreateBlockStmt:
			localEnvVars := make(map[string]runtime.RuntimeValue)
			localEnv := runtime.NewEnvironment(localEnvVars, r.Runtime.CurrEnv())
			r.Runtime.AddNewEnv(*localEnv)

			for _, node := range value.Nodes {
				r.RunNode(node)
			}
		case ast.CloseBlockStmt:
			r.Runtime.RemoveLastEnv()
		case ast.VarAssignStmt:
			val, err := r.Evaluator.EvaluateExpr(value.Expr)
			if err != nil {
				utils.EPrint(fmt.Sprintf("%s\n", err.Error()))
			}

			if val != nil {
				env := r.Runtime.CurrEnv()
				env.SetVar(value.Name, *runtime.NewRuntimeValue(val.Value))
			}
		case ast.VarReassignStmt:
			currEnv := r.Runtime.CurrEnv()
			val, env := currEnv.GetVar(value.Name)

			if val == nil {
				err := runtime.NewRuntimeError(runtime.UNDEFINED_IDENTIFIER, value.Name, expr.GetLine())
				utils.EPrint(err.Error())
			}

			exprVal, err := r.Evaluator.EvaluateExpr(value.Expr)
			if err != nil {
				utils.EPrint(fmt.Sprintf("%s\n", err.Error()))
			}

			if exprVal != nil {
				env.SetVar(value.Name, *exprVal)
			}
		case ast.IfStmt:
			var res bool
			res = r.EvalAndRunNode(value.Expr, value.IfBranch)

			if !res {
				elseIfBranches := value.ElseIfBranches

				if elseIfBranches != nil {
					for _, elseIfBranch := range *elseIfBranches {
						res = r.EvalAndRunNode(elseIfBranch.Expr, elseIfBranch.Branch)
						if res {
							break
						}
					}
				}
			}

			if !res {
				elseBranch := value.ElseBranch

				if elseBranch != nil {
					r.EvalAndRunNode(ast.NewLiteralExpr(tokens.TRUE, "", 0), elseBranch.Branch)
				}
			}
		}
	} else {
		_, err := r.Evaluator.EvaluateExpr(expr)
		if err != nil {
			utils.EPrint(fmt.Sprintf("%s\n", err.Error()))
		}
	}
}

func (r *Runner) Run() {
	curr := r.curr()

	if !r.IsAtEnd() {
		r.RunNode(curr)
		r.advance()
	}
}
