package lexer

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/0xmukesh/interpreter/internal/tokens"
	"github.com/0xmukesh/interpreter/internal/utils"
)

func (l *Lexer) LexStrLiterals() (*tokens.Token, *LexerError) {
	strLiteral := ""
	strLiteral += string(l.Char)
	closingQuoteFound := false

	for {
		nextChar := l.Peek()

		if nextChar == '\n' || nextChar == 0 {
			break
		}

		l.Read()

		strLiteral += string(l.Char)

		if l.Char == '"' {
			closingQuoteFound = true
			break
		}
	}

	if !closingQuoteFound {
		return nil, NewLexerError(UNTERMINATED_STRING, "Unterminated string.", l.Line)
	}

	return tokens.NewToken(tokens.STRING, strLiteral, strLiteral[1:len(strLiteral)-1], l.Line), nil
}

func (l *Lexer) LexNumLiterals() (*tokens.Token, *LexerError) {
	numLiteral := ""
	numLiteral += string(l.Char)
	decimalPointFound := false

	for {
		nextChar := l.Peek()

		if !(unicode.IsDigit(rune(nextChar)) || nextChar == '.') {
			break
		}

		if nextChar == '.' {
			if decimalPointFound {
				return nil, NewLexerError(UNTERMINATED_NUMBER, "Unterminated number.", l.Line)
			}
			decimalPointFound = true
		}

		l.Read()

		numLiteral += string(l.Char)
	}

	if strings.HasSuffix(numLiteral, ".") {
		return nil, NewLexerError(UNTERMINATED_NUMBER, "Unterminated number.", l.Line)
	}

	var literal string

	if !strings.Contains(numLiteral, ".") {
		literal = fmt.Sprintf("%s.%d", numLiteral, 0)
	} else {
		literal = numLiteral
	}

	return tokens.NewToken(tokens.NUMBER, numLiteral, literal, l.Line), nil
}

func (l *Lexer) LexIdentifier() (*tokens.Token, *LexerError) {
	identLiteral := ""
	identLiteral += string(l.Char)

	for {
		nextChar := l.Peek()

		_, existsInTknMap := utils.HasValueMap(tokens.TknLiteralMapping, string(nextChar))

		if utils.IsWhitespace(nextChar) || existsInTknMap {
			break
		}

		l.Read()
		identLiteral += string(l.Char)
	}

	identifierType, isReserved := utils.HasValueMap(tokens.ReservedKeywordsMapping, strings.ToUpper(identLiteral))
	if isReserved {
		return tokens.NewToken(*identifierType, identLiteral, "null", l.Line), nil
	}

	return tokens.NewToken(tokens.IDENTIFIER, identLiteral, "null", l.Line), nil
}
