package evaluator

import (
	"fmt"
)

type RuntimeValue struct {
	Value interface{}
}

func NewRuntimeValue(value interface{}) *RuntimeValue {
	return &RuntimeValue{
		Value: value,
	}
}

func (e RuntimeValue) String() string {
	return fmt.Sprintf("%v", e.Value)
}

type RuntimeError struct {
	Message string
	At      string
	Line    int
}

func NewRuntimeError(msg string, at string, line int) *RuntimeError {
	return &RuntimeError{
		Message: msg,
		At:      at,
		Line:    line,
	}
}
func (e RuntimeError) Error() string {
	return fmt.Sprintf("[line %d] Error at '%s': %s", e.Line, e.At, e.Message)
}
