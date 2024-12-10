package parser

import "fmt"

const (
	EXPRESSION_EXPECTED                  = "expression expected"
	EXPRESSION_AFTER_ASSIGNMENT_EXPECTED = "expression after assignment expected"
	VARIABLE_NAME_EXPECTED               = "variable name expected"

	MISSING_SEMICOLON = "missing ';'"
	MISSING_LPAREN    = "missing '('"
	MISSING_RPAREN    = "missing ')'"
	MISSING_RBRACE    = "missing '}'"
	MISSING_IF_BRANCH = "missing 'if' branch"

	INVALID_VARIABLE_NAME = "invalid variable name"
	INVALID_EXPRESSION    = "invalid expression"

	IDENTIFIER_ALREADY_EXISTS = "identifier already exists"
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
	return fmt.Sprintf("An error occured in parser: [line %d] Error at '%s': %s", e.Line, e.At, e.Message)
}
