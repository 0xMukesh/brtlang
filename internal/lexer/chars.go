package lexer

import (
	"github.com/0xmukesh/interpreter/internal/tokens"
)

// scans "=" and "==" tokens
func (l *Lexer) LexEqualChar() (*tokens.Token, *LexerError) {
	nextChar := l.peek()

	if nextChar == '=' {
		l.read()
		return tokens.NewToken(tokens.EQUAL_EQUAL, tokens.EQUAL_EQUAL.Literal(), "null", l.Line), nil
	}

	return tokens.NewToken(tokens.EQUAL, tokens.EQUAL.Literal(), "null", l.Line), nil
}

// scans "!" and "!=" tokens
func (l *Lexer) LexBangChar() (*tokens.Token, *LexerError) {
	nextChar := l.peek()

	if nextChar == '=' {
		l.read()
		return tokens.NewToken(tokens.BANG_EQUAL, tokens.BANG_EQUAL.Literal(), "null", l.Line), nil
	}

	return tokens.NewToken(tokens.BANG, tokens.BANG.Literal(), "null", l.Line), nil
}

// scans "<" and "<=" tokens
func (l *Lexer) LexLessChar() (*tokens.Token, *LexerError) {
	nextChar := l.peek()

	if nextChar == '=' {
		l.read()
		return tokens.NewToken(tokens.LESS_EQUAL, tokens.LESS_EQUAL.Literal(), "null", l.Line), nil
	}

	return tokens.NewToken(tokens.LESS, tokens.LESS.Literal(), "null", l.Line), nil
}

// scans ">" and ">=" tokens
func (l *Lexer) LexGreaterChar() (*tokens.Token, *LexerError) {
	nextChar := l.peek()

	if nextChar == '=' {
		l.read()
		return tokens.NewToken(tokens.GREATER_EQUAL, tokens.GREATER_EQUAL.Literal(), "null", l.Line), nil
	}

	return tokens.NewToken(tokens.GREATER, tokens.GREATER.Literal(), "null", l.Line), nil
}

// scans "/" token and "//" (comment)
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

// scans "&&" token
func (l *Lexer) LexAmpersandChar() (*tokens.Token, *LexerError) {
	nextChar := l.peek()

	if nextChar == '&' {
		l.read()
		return tokens.NewToken(tokens.AND, tokens.AND.Literal(), "null", l.Line), nil
	}

	return nil, NewLexerError(`missing "&" character`, l.Line)
}

// scans "||" token
func (l *Lexer) LexPipeChar() (*tokens.Token, *LexerError) {
	nextChar := l.peek()

	if nextChar == '|' {
		l.read()
		return tokens.NewToken(tokens.OR, tokens.OR.Literal(), "null", l.Line), nil
	}

	return nil, NewLexerError(`missing "|" character`, l.Line)
}

// scans "%" token
func (l *Lexer) LexModuloChar() (*tokens.Token, *LexerError) {
	return tokens.NewToken(tokens.MODULO, tokens.MODULO.Literal(), "null", l.Line), nil
}
