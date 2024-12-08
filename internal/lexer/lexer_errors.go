package lexer

import (
	"fmt"
)

type LexerError struct {
	Message string
	Line    int
}

func NewLexerError(msg string, line int) *LexerError {
	return &LexerError{
		Message: msg,
		Line:    line,
	}
}

func (e LexerError) Error() string {
	return fmt.Sprintf("An error occured in lexer: [line %d] Error: %s", e.Line, e.Message)
}
