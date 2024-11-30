package evaluator

import (
	"reflect"
	"strconv"

	"github.com/0xmukesh/interpreter/internal/ast"
	"github.com/0xmukesh/interpreter/internal/runtime"
	"github.com/0xmukesh/interpreter/internal/tokens"
)

type Evaluator struct {
	Ast     ast.Ast
	Runtime *runtime.Runtime
	Idx     int
}

func NewEvaluator(ast ast.Ast, runtime *runtime.Runtime) *Evaluator {
	return &Evaluator{
		Ast:     ast,
		Runtime: runtime,
		Idx:     0,
	}
}

func (e *Evaluator) isAtEnd() bool {
	return e.Idx >= len(e.Ast)
}

func (e *Evaluator) advance() {
	if !e.isAtEnd() {
		e.Idx++
	}
}

func (e *Evaluator) curr() ast.AstNode {
	if !e.isAtEnd() {
		return e.Ast[e.Idx]
	} else {
		return e.Ast[len(e.Ast)-1]
	}
}

func (e *Evaluator) Evaluate() (*runtime.RuntimeValue, *runtime.RuntimeError) {
	curr := e.curr()

	if e.isAtEnd() {
		return nil, nil
	}

	expr, isExpr := curr.Value.(ast.Expr)
	stmt, _ := curr.Value.(ast.Stmt)

	if isExpr {
		val, err := e.EvaluateExpr(expr)
		if err != nil {
			return nil, err
		}

		e.advance()
		return val, nil
	} else {
		stmtExpr := stmt.GetExpr()
		val, err := e.EvaluateExpr(stmtExpr)
		if err != nil {
			return nil, err
		}

		return val, nil
	}
}

func (e *Evaluator) EvaluateExpr(expr ast.Expr) (*runtime.RuntimeValue, *runtime.RuntimeError) {
	switch v := expr.(type) {
	case ast.LiteralExpr:
		return e.evaluteLiteralExpr(v)
	case ast.GroupingExpr:
		return e.evaluteGroupingExpr(v)
	case ast.LogicalExpr:
		return e.evaluteLogicalExpr(v)
	case ast.UnaryExpr:
		return e.evaluteUnaryExpr(v)
	case ast.BinaryExpr:
		return e.evaluateBinaryExpr(v)
	default:
		return nil, nil
	}
}

func (e *Evaluator) evaluteLiteralExpr(literalExpr ast.LiteralExpr) (*runtime.RuntimeValue, *runtime.RuntimeError) {
	switch literalExpr.TokenType {
	case tokens.STRING:
		return runtime.NewRuntimeValue(literalExpr.Value), nil
	case tokens.NUMBER:
		num, err := strconv.ParseFloat(literalExpr.Value, 64)
		if err != nil {
			return nil, runtime.NewRuntimeError(err.Error(), literalExpr.Value, literalExpr.Line)
		}

		return runtime.NewRuntimeValue(num), nil
	case tokens.TRUE:
		return runtime.NewRuntimeValue(true), nil
	case tokens.FALSE:
		return runtime.NewRuntimeValue(false), nil
	case tokens.NIL:
		return runtime.NewRuntimeValue(nil), nil
	case tokens.IDENTIFIER:
		currEnv := e.Runtime.CurrEnv()
		val, _ := currEnv.GetVar(literalExpr.Value)

		if val != nil {
			return runtime.NewRuntimeValue(val.Value), nil
		}

		return nil, runtime.NewRuntimeError(runtime.UNDEFINED_IDENTIFIER, literalExpr.Value, literalExpr.Line)
	default:
		return nil, nil
	}
}

func (e *Evaluator) evaluteGroupingExpr(groupingExpr ast.GroupingExpr) (*runtime.RuntimeValue, *runtime.RuntimeError) {
	return e.EvaluateExpr(groupingExpr.Expr)
}

func (e *Evaluator) evaluteLogicalExpr(logicalExpr ast.LogicalExpr) (*runtime.RuntimeValue, *runtime.RuntimeError) {
	left, err := e.EvaluateExpr(logicalExpr.Left)
	if err != nil {
		return nil, err
	}

	right, err := e.EvaluateExpr(logicalExpr.Right)
	if err != nil {
		return nil, err
	}

	operator := logicalExpr.Operator

	switch operator {
	case tokens.AND:
		leftBool, isLeftBool := left.Value.(bool)
		rightBool, isRightBool := right.Value.(bool)

		if !(isLeftBool && isRightBool) {
			return nil, runtime.NewRuntimeError(runtime.OperandsMustBeOfErrBuilder("bool"), logicalExpr.Operator.Literal(), logicalExpr.Line)
		}

		return runtime.NewRuntimeValue(leftBool && rightBool), nil
	case tokens.OR:
		leftBool, isLeftBool := left.Value.(bool)
		rightBool, isRightBool := right.Value.(bool)

		if !(isLeftBool && isRightBool) {
			return nil, runtime.NewRuntimeError(runtime.OperandsMustBeOfErrBuilder("bool"), logicalExpr.Operator.Literal(), logicalExpr.Line)
		}

		return runtime.NewRuntimeValue(leftBool || rightBool), nil
	default:
		return nil, runtime.NewRuntimeError(runtime.INVALID_OPERATOR, logicalExpr.Operator.Literal(), logicalExpr.Line)
	}
}

func (e *Evaluator) evaluteUnaryExpr(unaryExpr ast.UnaryExpr) (*runtime.RuntimeValue, *runtime.RuntimeError) {
	val, err := e.EvaluateExpr(unaryExpr.Expr)
	if err != nil {
		return nil, err
	}

	if unaryExpr.Operator == tokens.MINUS {
		valNum, isNum := val.Value.(float64)
		if !isNum {
			return nil, runtime.NewRuntimeError(runtime.OperandsMustBeOfErrBuilder("number"), unaryExpr.Operator.Literal(), unaryExpr.Line)
		}
		val = runtime.NewRuntimeValue(-1 * valNum)
	} else {
		if val.Value != true && val.Value != false {
			val = runtime.NewRuntimeValue(false)
		} else {
			val = runtime.NewRuntimeValue(!val.Value.(bool))
		}
	}

	return val, nil
}

func (e *Evaluator) evaluateBinaryExpr(binaryExpr ast.BinaryExpr) (*runtime.RuntimeValue, *runtime.RuntimeError) {
	left, err := e.EvaluateExpr(binaryExpr.Left)
	if err != nil {
		return nil, err
	}

	right, err := e.EvaluateExpr(binaryExpr.Right)
	if err != nil {
		return nil, err
	}

	operator := binaryExpr.Operator

	switch operator {
	case tokens.PLUS:
		leftStr, isLeftStr := left.Value.(string)
		rightStr, isRightStr := right.Value.(string)
		leftNum, isLeftNum := left.Value.(float64)
		rightNum, isRightNum := right.Value.(float64)

		if isLeftStr {
			if !isRightStr {
				return nil, runtime.NewRuntimeError(runtime.OperandsMustBeOfErrBuilder("string"), binaryExpr.Operator.Literal(), binaryExpr.Line)
			}

			return runtime.NewRuntimeValue(leftStr + rightStr), nil
		} else if isLeftNum {
			if !isRightNum {
				return nil, runtime.NewRuntimeError(runtime.OperandsMustBeOfErrBuilder("number"), binaryExpr.Operator.Literal(), binaryExpr.Line)
			}

			return runtime.NewRuntimeValue(leftNum + rightNum), nil
		} else {
			return nil, runtime.NewRuntimeError(runtime.OperandsMustBeOfErrBuilder("string", "number"), binaryExpr.Operator.Literal(), binaryExpr.Line)
		}
	case tokens.MINUS:
		leftNum, isLeftNum := left.Value.(float64)
		rightNum, isRightNum := right.Value.(float64)

		if !(isLeftNum && isRightNum) {
			return nil, runtime.NewRuntimeError(runtime.OperandsMustBeOfErrBuilder("number"), binaryExpr.Operator.Literal(), binaryExpr.Line)
		}

		return runtime.NewRuntimeValue(leftNum - rightNum), nil
	case tokens.STAR:
		leftNum, isLeftNum := left.Value.(float64)
		rightNum, isRightNum := right.Value.(float64)

		if !(isLeftNum && isRightNum) {
			return nil, runtime.NewRuntimeError(runtime.OperandsMustBeOfErrBuilder("number"), binaryExpr.Operator.Literal(), binaryExpr.Line)
		}

		return runtime.NewRuntimeValue(leftNum * rightNum), nil
	case tokens.SLASH:
		leftNum, isLeftNum := left.Value.(float64)
		rightNum, isRightNum := right.Value.(float64)

		if !(isLeftNum && isRightNum) {
			return nil, runtime.NewRuntimeError(runtime.OperandsMustBeOfErrBuilder("number"), binaryExpr.Operator.Literal(), binaryExpr.Line)
		}

		return runtime.NewRuntimeValue(leftNum / rightNum), nil
	case tokens.LESS:
		leftNum, isLeftNum := left.Value.(float64)
		rightNum, isRightNum := right.Value.(float64)

		if !(isLeftNum && isRightNum) {
			return nil, runtime.NewRuntimeError(runtime.OperandsMustBeOfErrBuilder("number"), binaryExpr.Operator.Literal(), binaryExpr.Line)
		}

		return runtime.NewRuntimeValue(leftNum < rightNum), nil
	case tokens.LESS_EQUAL:
		leftNum, isLeftNum := left.Value.(float64)
		rightNum, isRightNum := right.Value.(float64)

		if !(isLeftNum && isRightNum) {
			return nil, runtime.NewRuntimeError(runtime.OperandsMustBeOfErrBuilder("number"), binaryExpr.Operator.Literal(), binaryExpr.Line)
		}

		return runtime.NewRuntimeValue(leftNum <= rightNum), nil
	case tokens.GREATER:
		leftNum, isLeftNum := left.Value.(float64)
		rightNum, isRightNum := right.Value.(float64)

		if !(isLeftNum && isRightNum) {
			return nil, runtime.NewRuntimeError(runtime.OperandsMustBeOfErrBuilder("number"), binaryExpr.Operator.Literal(), binaryExpr.Line)
		}

		return runtime.NewRuntimeValue(leftNum > rightNum), nil
	case tokens.GREATER_EQUAL:
		leftNum, isLeftNum := left.Value.(float64)
		rightNum, isRightNum := right.Value.(float64)

		if !(isLeftNum && isRightNum) {
			return nil, runtime.NewRuntimeError(runtime.OperandsMustBeOfErrBuilder("number"), binaryExpr.Operator.Literal(), binaryExpr.Line)
		}

		return runtime.NewRuntimeValue(leftNum >= rightNum), nil
	case tokens.EQUAL_EQUAL:
		if reflect.TypeOf(left.Value) != reflect.TypeOf(right.Value) {
			return nil, runtime.NewRuntimeError(runtime.OperandsMustBeOfErrBuilder("same"), binaryExpr.Operator.Literal(), binaryExpr.Line)
		}

		return runtime.NewRuntimeValue(left.Value == right.Value), nil
	case tokens.BANG_EQUAL:
		if reflect.TypeOf(left.Value) != reflect.TypeOf(right.Value) {
			return nil, runtime.NewRuntimeError(runtime.OperandsMustBeOfErrBuilder("same"), binaryExpr.Operator.Literal(), binaryExpr.Line)
		}

		return runtime.NewRuntimeValue(left.Value != right.Value), nil
	default:
		return nil, runtime.NewRuntimeError(runtime.INVALID_OPERATOR, binaryExpr.Operator.Literal(), binaryExpr.Line)
	}
}
