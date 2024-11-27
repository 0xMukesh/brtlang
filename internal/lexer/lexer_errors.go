package lexer

import "fmt"

type LexerErrorType int

const (
	UNTERMINATED_NUMBER LexerErrorType = iota
	UNTERMINATED_STRING
)

type LexerError struct {
	Type    LexerErrorType
	Message string
	Line    int
}

func NewLexerError(errType LexerErrorType, msg string, line int) *LexerError {
	return &LexerError{
		Type:    errType,
		Message: msg,
		Line:    line,
	}
}

func (e LexerError) String() string {
	return fmt.Sprintf("[line %d] Error: %s", e.Line, e.Message)
}
