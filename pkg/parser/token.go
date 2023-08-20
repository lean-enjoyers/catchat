package parser

import (
	"unicode"
)

type Token int

// Tokens
const (
	ILLEGAL Token = iota
	EOF
	EOL
	WS
	IDENT      // Command name, parameter
	STR        // String encased in quotes
	DASH       // Indication of a short command
	DDASH      // Indication of a long command
	ASSIGNMENT // =
)

const (
	eof   = rune(0)
	dash  = rune('-')
	quote = rune('"')
	equal = rune('=')
)

func isWhiteSpace(ch rune) bool {
	return ch == ' ' || ch == '\t'
}

func isLetter(ch rune) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
}

func isQuote(ch rune) bool {
	return ch == quote
}

func isString(ch rune) bool {
	return isLetter(ch) || isQuote(ch)
}

func isDigit(ch rune) bool {
	return unicode.IsDigit(ch)
}

func isEOF(ch rune) bool {
	return ch == eof
}
