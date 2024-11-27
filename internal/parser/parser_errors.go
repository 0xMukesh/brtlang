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
}

func NewParserError(errType ParserErrorType, msg string) *ParserError {
	return &ParserError{
		Type:    errType,
		Message: msg,
	}
}

func (e ParserError) String() string {
	return fmt.Sprintf("Error: %s", e.Message)
}
