package parser

import "fmt"

const (
	EXPRESSION_EXPECTED = "expression expected"
	TOKEN_EXPECTED      = ""
)

type ParserError struct {
	Message string
	At      string
	Line    int
}

func NewParserError(msg string, at string, line int) *ParserError {
	return &ParserError{
		Message: msg,
		At:      at,
		Line:    line,
	}
}

func (e ParserError) Error() string {
	return fmt.Sprintf("[line %d] Error at '%s': %s", e.Line, e.At, e.Message)
}
