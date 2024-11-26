package lexer

import "github.com/0xmukesh/interpreter/internal/tokens"

// handles scanning "=" and "==" tokens
func LexEqualChar(l *Lexer) (*tokens.Token, *LexerError) {
	nextChar := l.PeekChar()

	if nextChar == '=' {
		l.ReadChar()
		return tokens.NewToken(tokens.EQUAL_EQUAL, tokens.EQUAL_EQUAL.Literal(), "null"), nil
	}

	return tokens.NewToken(tokens.EQUAL, tokens.EQUAL.Literal(), "null"), nil
}

// handles scanning "!" and "!=" tokens
func LexBangChar(l *Lexer) (*tokens.Token, *LexerError) {
	nextChar := l.PeekChar()

	if nextChar == '=' {
		l.ReadChar()
		return tokens.NewToken(tokens.BANG_EQUAL, tokens.BANG_EQUAL.Literal(), "null"), nil
	}

	return tokens.NewToken(tokens.BANG, tokens.BANG.Literal(), "null"), nil
}

// handles scanning "<" and "<=" tokens
func LexLessChar(l *Lexer) (*tokens.Token, *LexerError) {
	nextChar := l.PeekChar()

	if nextChar == '=' {
		l.ReadChar()
		return tokens.NewToken(tokens.LESS_EQUAL, tokens.LESS_EQUAL.Literal(), "null"), nil
	}

	return tokens.NewToken(tokens.LESS, tokens.LESS.Literal(), "null"), nil
}

// handles scanning ">" and ">=" tokens
func LexGreaterChar(l *Lexer) (*tokens.Token, *LexerError) {
	nextChar := l.PeekChar()

	if nextChar == '=' {
		l.ReadChar()
		return tokens.NewToken(tokens.GREATER_EQUAL, tokens.GREATER_EQUAL.Literal(), "null"), nil
	}

	return tokens.NewToken(tokens.GREATER, tokens.GREATER.Literal(), "null"), nil
}

// handles scanning "/" token and "//" (aka comment)
func LexSlashChar(l *Lexer) (*tokens.Token, *LexerError) {
	nextChar := l.PeekChar()

	if nextChar == '/' {
		for {
			l.ReadChar()
			// read until end of the line/file
			if l.Char == '\n' || l.Char == 0 {
				break
			}
		}

		return tokens.NewToken(tokens.IGNORE, "", "null"), nil
	}

	return tokens.NewToken(tokens.SLASH, tokens.SLASH.Literal(), "null"), nil
}
