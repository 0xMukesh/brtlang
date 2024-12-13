package lexer

import (
	"github.com/0xmukesh/interpreter/internal/tokens"
)

func (l *Lexer) LexDoubleCharBuilder(expectedNextChar byte, doubleCharTknType, singleCharTknType tokens.TokenType) (*tokens.Token, *LexerError) {
	nextChar := l.peek()

	if nextChar == expectedNextChar {
		l.read()
		return tokens.NewToken(doubleCharTknType, doubleCharTknType.Literal(), "null", l.Line), nil
	}

	return tokens.NewToken(singleCharTknType, singleCharTknType.Literal(), "null", l.Line), nil
}

// scans "=" and "==" tokens
func (l *Lexer) LexEqualChar() (*tokens.Token, *LexerError) {
	return l.LexDoubleCharBuilder('=', tokens.EQUAL_EQUAL, tokens.EQUAL)
}

// scans "!" and "!=" tokens
func (l *Lexer) LexBangChar() (*tokens.Token, *LexerError) {
	return l.LexDoubleCharBuilder('=', tokens.BANG_EQUAL, tokens.BANG)
}

// scans "<" and "<=" tokens
func (l *Lexer) LexLessChar() (*tokens.Token, *LexerError) {
	return l.LexDoubleCharBuilder('=', tokens.LESS_EQUAL, tokens.LESS)
}

// scans ">" and ">=" tokens
func (l *Lexer) LexGreaterChar() (*tokens.Token, *LexerError) {
	return l.LexDoubleCharBuilder('=', tokens.GREATER_EQUAL, tokens.GREATER)
}

// scans "+" and "++" tokens
func (l *Lexer) LexPlusChar() (*tokens.Token, *LexerError) {
	return l.LexDoubleCharBuilder('+', tokens.PLUS_PLUS, tokens.PLUS)
}

// scans "-" and "--" tokens
func (l *Lexer) LexMinusChar() (*tokens.Token, *LexerError) {
	return l.LexDoubleCharBuilder('-', tokens.MINUS_MINUS, tokens.MINUS)
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
