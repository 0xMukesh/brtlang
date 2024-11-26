package lexer

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/0xmukesh/interpreter/internal/tokens"
	"github.com/0xmukesh/interpreter/internal/utils"
)

func LexStrLiterals(l *Lexer) (*tokens.Token, *LexerError) {
	strLiteral := ""
	strLiteral += string(l.Char)
	closingQuoteFound := false

	for {
		nextChar := l.PeekChar()

		if nextChar == '\n' || nextChar == 0 {
			break
		}

		l.ReadChar()

		strLiteral += string(l.Char)

		if l.Char == '"' {
			closingQuoteFound = true
			break
		}
	}

	if !closingQuoteFound {
		return nil, NewLexerError(UNTERMINATED_STRING, "Unterminated string.")
	}

	return tokens.NewToken(tokens.STRING, strLiteral, strLiteral[1:len(strLiteral)-1]), nil
}

func LexNumLiterals(l *Lexer) (*tokens.Token, *LexerError) {
	numLiteral := ""
	numLiteral += string(l.Char)
	decimalPointFound := false

	for {
		nextChar := l.PeekChar()

		if !(unicode.IsDigit(rune(nextChar)) || nextChar == '.') {
			break
		}

		if nextChar == '.' {
			if decimalPointFound {
				return nil, NewLexerError(UNTERMINATED_NUMBER, "Unterminated number.")
			}
			decimalPointFound = true
		}

		l.ReadChar()

		numLiteral += string(l.Char)
	}

	if strings.HasSuffix(numLiteral, ".") {
		return nil, NewLexerError(UNTERMINATED_NUMBER, "Unterminated number.")
	}

	var literal string

	if !strings.Contains(numLiteral, ".") {
		literal = fmt.Sprintf("%s.%d", numLiteral, 0)
	} else {
		literal = numLiteral
	}

	return tokens.NewToken(tokens.NUMBER, numLiteral, literal), nil
}

func LexIdentifier(l *Lexer) (*tokens.Token, *LexerError) {
	identLiteral := ""
	identLiteral += string(l.Char)

	for {
		nextChar := l.PeekChar()

		_, existsInTknMap := utils.HasValueMap(tokens.TknLiteralMapping, string(nextChar))

		if utils.IsWhitespace(nextChar) || existsInTknMap {
			break
		}

		l.ReadChar()
		identLiteral += string(l.Char)
	}

	identifierType, isReserved := utils.HasValueMap(tokens.ReservedKeywordsMapping, strings.ToUpper(identLiteral))
	if isReserved {
		return tokens.NewToken(*identifierType, identLiteral, "null"), nil
	}

	return tokens.NewToken(tokens.IDENTIFIER, identLiteral, "null"), nil
}
