package runtime

import (
	"fmt"
	"strings"
)

type RuntimeError struct {
	Message string
	At      string
	Line    int
}

const (
	UNDEFINED_IDENTIFIER      = "damn bruv, this identifier got that invisible drip"
	IDENTIFIER_ALREADY_EXISTS = "nah, the sequel ain't happening for this identifier"

	INVALID_OPERAND_TEMPLATE = "this operand ain't it, chief. got %s"
	INVALID_OPERATOR         = "this operator ain't it, chief"
)

func (e RuntimeError) Error() string {
	return fmt.Sprintf("[line %d] hell naw, im done you caused a runtime error at '%s': %s", e.Line, e.At, e.Message)
}

func NewRuntimeError(msg string, at string, line int) *RuntimeError {
	return &RuntimeError{
		Message: msg,
		At:      at,
		Line:    line,
	}
}

func OperandsMustBeOfErrBuilder(expectedTypes ...string) string {
	if len(expectedTypes) == 1 {
		if expectedTypes[0] == "same" {
			return "the operands out here are wildin'. they need to be of the same type"
		}
		return fmt.Sprintf("the operands out here are wildin'. they need to be of type %s", expectedTypes[0])
	} else {
		return fmt.Sprintf("the operands out here are wildin'. they need to be of type %s or %s", strings.Join(expectedTypes[:len(expectedTypes)-1], ", "), expectedTypes[len(expectedTypes)-1])
	}
}

func ExpectedExprErrBuilder(expectedExprType string) string {
	return fmt.Sprintf("yo, where's the vibe? i was an expecting %s", expectedExprType)
}
