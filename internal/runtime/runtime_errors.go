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
	UNDEFINED_IDENTIFIER      = "undefined identifier"
	IDENTIFIER_ALREADY_EXISTS = "identifier already exists"

	INVALID_OPERAND  = "invalid operand"
	INVALID_OPERATOR = "invalid operator"

	EXPECTED_EXPRESSION = "expected expression"

	OPERANDS_MUST_BE_OF = "operands must be of"
)

func (e RuntimeError) Error() string {
	return fmt.Sprintf("[line %d] Error at '%s': %s", e.Line, e.At, e.Message)
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
			return fmt.Sprintf("%s same type", OPERANDS_MUST_BE_OF)
		}
		return fmt.Sprintf("%s type %s", OPERANDS_MUST_BE_OF, expectedTypes[0])
	} else {
		return fmt.Sprintf("%s type %s or %s", OPERANDS_MUST_BE_OF, strings.Join(expectedTypes[:len(expectedTypes)-1], ", "), expectedTypes[len(expectedTypes)-1])
	}
}

func ExpectedExprErrBuilder(expectedExprType string) string {
	parts := strings.Split(EXPECTED_EXPRESSION, " ")
	return fmt.Sprintf("%s %s %s", parts[0], expectedExprType, parts[1])
}
