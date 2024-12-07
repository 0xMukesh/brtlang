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
		err := runtime.NewRuntimeError(runtime.ExpectedExprErrBuilder("boolean"), expr.ParseExpr(), expr.GetLine())
		utils.EPrint(err.Error())
	}

	if evaledCondition != nil {
		conditionVal := (*evaledCondition).Value
		if conditionVal != true && conditionVal != false {
			err := runtime.NewRuntimeError(runtime.ExpectedExprErrBuilder("boolean"), evaledCondition.String(), expr.GetLine())
			utils.EPrint(err.Error())
		}

		if conditionVal == true {
			r.RunNode(node, nil, nil)
			return true
		}
	}

	return false
}

func (r *Runner) RunNode(node ast.AstNode, envVars *runtime.RuntimeVarMapping, envFuncs *runtime.RuntimeFuncMapping) {
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
			localEnvVars := make(runtime.RuntimeVarMapping)
			localEnvFuncs := make(runtime.RuntimeFuncMapping)

			if envVars != nil {
				localEnvVars = *envVars
			}
			if envFuncs != nil {
				localEnvFuncs = *envFuncs
			}

			localEnv := runtime.NewEnvironment(localEnvVars, localEnvFuncs, r.Runtime.CurrEnv())
			r.Runtime.AddNewEnv(*localEnv)

			for _, node := range value.Nodes {
				r.RunNode(node, &localEnvVars, &localEnvFuncs)
			}
		case ast.CloseBlockStmt:
			r.Runtime.RemoveLastEnv()
		case ast.VarAssignStmt:
			val, err := r.Evaluator.EvaluateExpr(value.Expr)
			if err != nil {
				utils.EPrint(fmt.Sprintf("%s\n", err.Error()))
			}

			if val != nil {
				currEnv := r.Runtime.CurrEnv()
				varVal, _ := currEnv.GetVar(value.Name)
				if varVal != nil {
					err := runtime.NewRuntimeError(runtime.IDENTIFIER_ALREADY_EXISTS, value.Name, value.Line)
					utils.EPrint(err.Error())
				}

				currEnv.SetVar(value.Name, *runtime.NewRuntimeValue(val.Value))
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
		case ast.WhileStmt:
			for {
				val, err := r.Evaluator.EvaluateExpr(value.Expr)
				if err != nil {
					utils.EPrint(fmt.Sprintf("%s\n", err.Error()))
				}

				if val != nil {
					conditionVal, isConditionBool := val.Value.(bool)

					if !isConditionBool {
						err := runtime.NewRuntimeError(runtime.ExpectedExprErrBuilder("bool"), value.Expr.ParseExpr(), value.Expr.GetLine())
						utils.EPrint(err.Error())
					}

					if !conditionVal {
						break
					}

					r.RunNode(value.Node, nil, nil)
				}
			}
		case ast.FuncDeclarationStmt:
			currEnv := r.Runtime.CurrEnv()
			funcNode, _ := currEnv.GetFunc(value.Name)
			if funcNode != nil {
				err := runtime.NewRuntimeError(runtime.IDENTIFIER_ALREADY_EXISTS, value.Name, value.Line)
				utils.EPrint(err.Error())
			}

			currEnv.SetFunc(value.Name, value.Node, value.Args)
		case ast.FuncCallStmt:
			currEnv := r.Runtime.CurrEnv()
			funcMappingPtr, _ := currEnv.GetFunc(value.Name)

			if funcMappingPtr == nil {
				err := runtime.NewRuntimeError(runtime.UNDEFINED_IDENTIFIER, value.Name, value.Line)
				utils.EPrint(err.Error())
			}

			funcMapping := *funcMappingPtr

			if len(value.Args) != len(funcMapping.Args) {
				err := runtime.NewRuntimeError(fmt.Sprintf("invalid number of arguments. expected %d arguments but got %d arguments", len(funcMapping.Args), len(value.Args)), value.Name, value.Line)
				utils.EPrint(err.Error())
			}

			argsMapping := make(runtime.RuntimeVarMapping)

			for i, arg := range value.Args {
				argName := funcMapping.Args[i].Value.(ast.LiteralExpr).Value
				argValue, _ := r.Evaluator.EvaluateExpr(arg.Value.(ast.Expr))
				argsMapping[argName] = *argValue
			}

			r.RunNode(funcMapping.Node, &argsMapping, nil)
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
		r.RunNode(curr, nil, nil)
		r.advance()
	}
}
