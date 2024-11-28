package utils

import (
	"fmt"
	"os"
	"unicode/utf8"
)

func EPrint(err string) {
	fmt.Fprint(os.Stderr, err)
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
