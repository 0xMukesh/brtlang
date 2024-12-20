package parser

import "fmt"

const (
	EXPRESSION_EXPECTED          = "yo, where's the vibe? i was expecting an expression over here"
	EXPRESSION_EXPECTED_TEMPLATE = "yo, where's the vibe? i was expecting an %s expression over here"

	STATEMENT_EXPECTED = "yo, where's the vibe? i was expecting a statement over here"

	EXPRESSION_AFTER_ASSIGNMENT_EXPECTED = "vibe check failed, blud misunderstood the assignment. i was expecting an expression after assignment"
	VARIABLE_NAME_EXPECTED               = "yo, where's the vibe? i was expecting a variable name over here"

	MISSING_SEMICOLON = "nahh, you left me hanging. where's ';' at?"
	MISSING_LPAREN    = "bruh, where's the '('? you can't just skip it like that"
	MISSING_RPAREN    = "nahh, you left me hanging. where's ')' at?"
	MISSING_RBRACE    = "nahh, you left me hanging. where's '}' at?"
	MISSING_IF_BRANCH = "bruh, where's the 'if' branch? you can't just skip it like that"

	INVALID_TOKEN_TYPE_TEMPLATE = "ay, that token isn't allowed. expected %s token"

	INVALID_VARIABLE_NAME      = "who tf even allowed you to name this variable?"
	INVALID_EXPRESSION         = "invalid expression? wow, didn't know we were coding in clown mode today"
	INVALID_STATEMENT          = "invalid statement? wow, didn't know we were coding in clown mode today"
	INVALID_STATEMENT_TEMPLATE = "invalid statement? wow, didn't know we were coding in clown mode today. expected %s statement"
	INVALID_OPERATOR_TEMPLATE  = "this operator ain't it, chief. expected a comparision operator but got %s"

	IDENTIFIER_ALREADY_EXISTS = "nah, the sequel ain't happening for this identifier"
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
	return fmt.Sprintf("[line %d] hell naw, im done with you. you caused a parser error at '%s': %s", e.Line, e.At, e.Message)
}
