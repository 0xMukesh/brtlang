package utils

import (
	"fmt"
	"os"
	"slices"
	"strconv"
	"unicode"
	"unicode/utf8"

	"github.com/0xmukesh/interpreter/internal/runtime"
	"github.com/0xmukesh/interpreter/internal/tokens"
)

func EPrint(err string) {
	fmt.Fprintf(os.Stderr, "%s\n", err)
	os.Exit(1)
}

func IsString(char byte) bool {
	return utf8.Valid([]byte{char})
}

// check if there exists any key-value pair where V == e and if it exists, it returns the key
func HasValueMap[K comparable, V comparable](dict map[K]V, e V) (*K, bool) {
	for k, v := range dict {
		if e == v {
			return &k, true
		}
	}

	return nil, false
}

func HasValueArray[E comparable](arr []E, e E) bool {
	for _, v := range arr {
		if v == e {
			return true
		}
	}

	return false
}

// checks if `b` is a subset of `a`
func IsSubset[T ~[]E, E comparable](a, b T) bool {
	for _, v := range a {
		if !slices.Contains(b, v) {
			return false
		}
	}

	return true
}

func IsAlphaOnly(s string) bool {
	for _, v := range s {
		if !unicode.IsLetter(v) {
			return false
		}
	}

	return true
}

func IsReservedKeyword(keyword string) bool {
	for _, v := range tokens.TknLiteralMapping {
		if v == keyword {
			return true
		}
	}
	return false
}

func IsNativeFunc(funcName string) bool {
	for k := range runtime.NativeFnsReturnExprMapping {
		if k == funcName {
			return true
		}
	}
	return false
}

func IsNumber(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

func FindLineNumber(charIdx int, src []byte) int {
	lineNumber := 1

	for i := range src {
		if i == charIdx {
			return lineNumber
		}

		if src[i] == '\n' {
			lineNumber++
		}

	}

	return lineNumber
}

func IsWhitespace(char byte) bool {
	return char == '\n' || char == ' ' || char == '\t' || char == 0
}
