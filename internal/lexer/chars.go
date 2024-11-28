package lexer

import "github.com/0xmukesh/interpreter/internal/tokens"

// handles scanning "=" and "==" tokens
func (l *Lexer) LexEqualChar() (*tokens.Token, *LexerError) {
	nextChar := l.peek()

	if nextChar == '=' {
		l.read()
		return tokens.NewToken(tokens.EQUAL_EQUAL, tokens.EQUAL_EQUAL.Literal(), "null", l.Line), nil
	}

	return tokens.NewToken(tokens.EQUAL, tokens.EQUAL.Literal(), "null", l.Line), nil
}

// handles scanning "!" and "!=" tokens
func (l *Lexer) LexBangChar() (*tokens.Token, *LexerError) {
	nextChar := l.peek()

	if nextChar == '=' {
		l.read()
		return tokens.NewToken(tokens.BANG_EQUAL, tokens.BANG_EQUAL.Literal(), "null", l.Line), nil
	}

	return tokens.NewToken(tokens.BANG, tokens.BANG.Literal(), "null", l.Line), nil
}

// handles scanning "<" and "<=" tokens
func (l *Lexer) LexLessChar() (*tokens.Token, *LexerError) {
	nextChar := l.peek()

	if nextChar == '=' {
		l.read()
		return tokens.NewToken(tokens.LESS_EQUAL, tokens.LESS_EQUAL.Literal(), "null", l.Line), nil
	}

	return tokens.NewToken(tokens.LESS, tokens.LESS.Literal(), "null", l.Line), nil
}

// handles scanning ">" and ">=" tokens
func (l *Lexer) LexGreaterChar() (*tokens.Token, *LexerError) {
	nextChar := l.peek()

	if nextChar == '=' {
		l.read()
		return tokens.NewToken(tokens.GREATER_EQUAL, tokens.GREATER_EQUAL.Literal(), "null", l.Line), nil
	}

	return tokens.NewToken(tokens.GREATER, tokens.GREATER.Literal(), "null", l.Line), nil
}

// handles scanning "/" token and "//" (aka comment)
func (l *Lexer) LexSlashChar() (*tokens.Token, *LexerError) {
	nextChar := l.peek()

	if nextChar == '/' {
		for {
			l.read()
			// read until end of the line/file
			if l.Char == '\n' || l.Char == 0 {
				break
			}
		}

		return tokens.NewToken(tokens.IGNORE, "", "null", l.Line), nil
	}

	return tokens.NewToken(tokens.SLASH, tokens.SLASH.Literal(), "null", l.Line), nil
}
