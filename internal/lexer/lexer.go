package lexer

import (
	"unicode"
	"unicode/utf8"

	"github.com/0xmukesh/interpreter/internal/tokens"
	"github.com/0xmukesh/interpreter/internal/utils"
)

type Lexer struct {
	Src  []byte
	Line int
	Idx  int
	Char byte
}

func NewLexer(src []byte) *Lexer {
	return &Lexer{
		Src:  src,
		Idx:  0,
		Char: 0,
	}
}

func (l *Lexer) isAtEnd() bool {
	return l.Idx >= len(l.Src)
}

func (l *Lexer) read() {
	if l.isAtEnd() {
		l.Char = 0
		return
	}

	l.Char = l.Src[l.Idx]
	l.Line = utils.FindLineNumber(l.Idx, l.Src)
	l.Idx++
}

func (l *Lexer) peek() byte {
	if l.isAtEnd() {
		return 0
	}

	return l.Src[l.Idx]
}

func (l *Lexer) LexAll() ([]tokens.Token, *LexerError) {
	var tkns []tokens.Token

	for !l.isAtEnd() {
		l.read()
		tkn, err := l.Lex()
		if err != nil {
			return nil, err
		}

		if tkn != nil {
			tkns = append(tkns, *tkn)
		}
	}

	return tkns, nil
}

func (l *Lexer) Lex() (*tokens.Token, *LexerError) {
	if l.Char == 0 {
		return tokens.NewToken(tokens.EOF, "", "null", l.Line), nil
	}

	if utils.IsWhitespace(l.Char) {
		return tokens.NewToken(tokens.IGNORE, "", "null", l.Line), nil
	}

	if utf8.Valid([]byte{l.Char}) {
		tknType, doesExists := utils.HasValueMap(tokens.TknLiteralMapping, string(l.Char))

		if doesExists {
			var tkn *tokens.Token
			var err *LexerError

			switch *tknType {
			case tokens.EQUAL:
				tkn, err = l.LexEqualChar()
			case tokens.BANG:
				tkn, err = l.LexBangChar()
			case tokens.LESS:
				tkn, err = l.LexLessChar()
			case tokens.GREATER:
				tkn, err = l.LexGreaterChar()
			case tokens.SLASH:
				tkn, err = l.LexSlashChar()
			default:
				tkn = tokens.NewToken(*tknType, string(l.Char), "null", l.Line)
			}

			return tkn, err
		} else {
			if l.Char == '&' {
				tkn, err := l.LexAmpersandChar()
				return tkn, err
			} else if l.Char == '|' {
				tkn, err := l.LexPipeChar()
				return tkn, err
			} else if l.Char == '"' {
				tkn, err := l.LexStrLiterals()
				return tkn, err
			} else if unicode.IsDigit(rune(l.Char)) {
				tkn, err := l.LexNumLiterals()
				return tkn, err
			} else {
				tkn, err := l.LexIdentifier()
				return tkn, err
			}
		}
	}

	return tokens.NewToken(tokens.ILLEGAL, "", string(l.Char), l.Line), nil
}
