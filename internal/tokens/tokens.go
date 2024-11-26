package tokens

type Token struct {
	Type    TokenType
	Lexeme  string
	Literal string
	Line    int
}

func NewToken(tokenType TokenType, lexeme string, literal string) *Token {
	return &Token{
		Type:    tokenType,
		Lexeme:  lexeme,
		Literal: literal,
	}
}
