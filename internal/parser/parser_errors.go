package parser

import "fmt"

type ParserErrorType int

const (
	EXPRESSION_EXPECTED ParserErrorType = iota
	TOKEN_EXPECTED
)

type ParserError struct {
	Type    ParserErrorType
	Message string
	At      string
	Line    int
}

func NewParserError(errType ParserErrorType, msg string, at string, line int) *ParserError {
	return &ParserError{
		Type:    errType,
		Message: msg,
		At:      at,
		Line:    line,
	}
}

func (e ParserError) String() string {
	return fmt.Sprintf("[line %d] Error at '%s': %s", e.Line, e.At, e.Message)
}
