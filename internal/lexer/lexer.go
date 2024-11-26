package lexer

import (
	"unicode"
	"unicode/utf8"

	"github.com/0xmukesh/interpreter/internal/tokens"
	"github.com/0xmukesh/interpreter/internal/utils"
)

type Lexer struct {
	Content []byte
	Line    int
	Idx     int
	Char    byte
}

func NewLexer(content []byte) *Lexer {
	return &Lexer{
		Content: content,
		Idx:     0,
		Char:    0,
	}
}

func (l *Lexer) ReadChar() {
	if l.Idx >= len(l.Content) {
		l.Char = 0
		return
	}

	l.Char = l.Content[l.Idx]
	l.Line = utils.FindLineNumber(l.Idx, l.Content)
	l.Idx++
}

func (l *Lexer) PeekChar() byte {
	if l.Idx >= len(l.Content) {
		return 0
	}

	return l.Content[l.Idx]
}

func (l *Lexer) LexTokenWithPeek() {

}

func (l *Lexer) ParseChar() (*tokens.Token, *LexerError) {
	if l.Char == 0 {
		return tokens.NewToken(tokens.EOF, "", "null"), nil
	}

	if utils.IsWhitespace(l.Char) {
		return tokens.NewToken(tokens.IGNORE, "", "null"), nil
	}

	if utf8.Valid([]byte{l.Char}) {
		tknType, doesExists := utils.HasValueMap(tokens.TknLiteralMapping, string(l.Char))

		if doesExists {
			if *tknType == tokens.EQUAL {
				tkn, err := LexEqualChar(l)
				return tkn, err
			}

			if *tknType == tokens.BANG {
				tkn, err := LexBangChar(l)
				return tkn, err
			}

			if *tknType == tokens.LESS {
				tkn, err := LexLessChar(l)
				return tkn, err
			}

			if *tknType == tokens.GREATER {
				tkn, err := LexGreaterChar(l)
				return tkn, err
			}

			if *tknType == tokens.SLASH {
				tkn, err := LexSlashChar(l)
				return tkn, err
			}

			return tokens.NewToken(*tknType, string(l.Char), "null"), nil
		} else {
			if l.Char == '"' {
				tkn, err := LexStrLiterals(l)
				return tkn, err
			} else if unicode.IsDigit(rune(l.Char)) {
				tkn, err := LexNumLiterals(l)
				return tkn, err
			} else {
				tkn, err := LexIdentifier(l)
				return tkn, err
			}
		}
	}

	return tokens.NewToken(tokens.ILLEGAL, "", string(l.Char)), nil
}
