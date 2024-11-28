package evaluator

import (
	"fmt"
	"strconv"

	"github.com/0xmukesh/interpreter/internal/ast"
	"github.com/0xmukesh/interpreter/internal/tokens"
)

type Evaluator struct {
	Ast ast.Ast
	Idx int
}

func NewEvaluator(ast ast.Ast) *Evaluator {
	return &Evaluator{
		Ast: ast,
		Idx: 0,
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

func (e *Evaluator) Evaluate() (*RuntimeValue, *RuntimeError) {
	curr := e.curr()

	if e.isAtEnd() {
		return nil, nil
	}

	val, err := e.evaluteExpr(curr.Expr)
	if err != nil {
		return nil, err
	}

	e.advance()
	return val, nil
}

func (e *Evaluator) evaluteExpr(expr ast.Expr) (*RuntimeValue, *RuntimeError) {
	switch v := expr.(type) {
	case ast.LiteralExpr:
		return e.evaluteLiteralExpr(v)
	case ast.GroupingExpr:
		return e.evaluteGroupingExpr(v)
	case ast.UnaryExpr:
		return e.evaluteUnaryExpr(v)
	case ast.BinaryExpr:
		return e.evaluateBinaryExpr(v)
	default:
		return nil, nil
	}
}

func (e *Evaluator) evaluteLiteralExpr(literalExpr ast.LiteralExpr) (*RuntimeValue, *RuntimeError) {
	switch literalExpr.TokenType {
	case tokens.STRING:
		return NewRuntimeValue(literalExpr.Value), nil
	case tokens.NUMBER:
		num, err := strconv.ParseFloat(literalExpr.Value, 64)
		if err != nil {
			return nil, NewRuntimeError(err.Error(), literalExpr.Value, literalExpr.Line)
		}

		return NewRuntimeValue(num), nil
	case tokens.TRUE:
		return NewRuntimeValue(true), nil
	case tokens.FALSE:
		return NewRuntimeValue(false), nil
	case tokens.NIL:
		return NewRuntimeValue(nil), nil
	default:
		return nil, nil
	}
}

func (e *Evaluator) evaluteGroupingExpr(groupingExpr ast.GroupingExpr) (*RuntimeValue, *RuntimeError) {
	return e.evaluteExpr(groupingExpr.Expr)
}

func (e *Evaluator) evaluteUnaryExpr(unaryExpr ast.UnaryExpr) (*RuntimeValue, *RuntimeError) {
	val, err := e.evaluteExpr(unaryExpr.Expr)
	if err != nil {
		return nil, err
	}

	if unaryExpr.Operator == tokens.MINUS {
		val = NewRuntimeValue(fmt.Sprintf("%s%v", unaryExpr.Operator.Literal(), val))
	} else {
		if val.Value != true && val.Value != false {
			val = NewRuntimeValue(false)
		} else {
			val = NewRuntimeValue(!val.Value.(bool))
		}
	}

	return val, nil
}

func (e *Evaluator) evaluateBinaryExpr(binaryExpr ast.BinaryExpr) (*RuntimeValue, *RuntimeError) {
	left, err := e.evaluteExpr(binaryExpr.Left)
	if err != nil {
		return nil, err
	}

	right, err := e.evaluteExpr(binaryExpr.Right)
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

		if isLeftStr && isRightStr {
			return NewRuntimeValue(leftStr + rightStr), nil
		} else if isLeftNum && isRightNum {
			return NewRuntimeValue(leftNum + rightNum), nil
		} else {
			return nil, NewRuntimeError("either both of the operands must be strings or else both of them must be numbers", binaryExpr.Operator.Literal(), binaryExpr.Line)
		}
	default:
		return nil, NewRuntimeError("invalid binary expression operator", binaryExpr.Operator.Literal(), binaryExpr.Line)
	}
}
