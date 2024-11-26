package lexer

type LexerErrorType int

const (
	UNTERMINATED_NUMBER LexerErrorType = iota
	UNTERMINATED_STRING
)

type LexerError struct {
	Type    LexerErrorType
	Message string
}

func NewLexerError(errType LexerErrorType, msg string) *LexerError {
	return &LexerError{
		Type:    errType,
		Message: msg,
	}
}
