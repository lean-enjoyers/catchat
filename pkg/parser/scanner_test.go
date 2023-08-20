package parser

import (
	"strings"
	"testing"
)

func TestIdentifers(t *testing.T) {
	input := "hello world"
	reader := strings.NewReader(input)
	sc := NewScanner(reader)

	identifiers := []string{"hello", "world"}
	ident_c := 0
	expected := []Token{IDENT, IDENT, EOF}

	for i, v := range expected {
		token := sc.Advance()
		tok := token.Token
		lit := token.string
		switch v {
		case IDENT:
			if tok != IDENT || lit != identifiers[ident_c] {
				t.Fatalf("'%s' not identified as IDENTIFIER", lit)
			}
			ident_c += 1
		case WS:
			if tok != WS {
				t.Fatalf("'%s' not identified as WHITESPACE", lit)
			}
		case EOF:
			if tok != EOF {
				t.Fatalf("Failed to identify EOF")
			}
			if i != len(expected)-1 {
				t.Fatalf("Failed to read all")
			}
		}
	}
}

func TestIdentifersWeirdSpacing(t *testing.T) {
	input := "  hello\t\t\t world\t\t"
	reader := strings.NewReader(input)
	sc := NewScanner(reader)

	identifiers := []string{"hello", "world"}
	ident_c := 0
	expected := []Token{IDENT, IDENT, EOF}

	for i, v := range expected {
		token := sc.Advance()
		tok := token.Token
		lit := token.string
		switch v {
		case IDENT:
			if tok != IDENT || lit != identifiers[ident_c] {
				t.Fatalf("'%s' not identified as IDENTIFIER", lit)
			}
			ident_c += 1
		case WS:
			if tok != WS {
				t.Fatalf("'%s' not identified as WHITESPACE", lit)
			}
		case EOF:
			if tok != EOF {
				t.Fatalf("Failed to identify EOF")
			}
			if i != len(expected)-1 {
				t.Fatalf("Failed to read all")
			}
		}
	}
}

func TestIdentifersWeirdAgain(t *testing.T) {
	input := "a\tb c \td\t e  f\t\t g"
	reader := strings.NewReader(input)
	sc := NewScanner(reader)

	identifiers := []string{"a", "b", "c", "d", "e", "f", "g"}
	ident_c := 0
	expected := []Token{IDENT, IDENT, IDENT, IDENT, IDENT, IDENT, IDENT, EOL}

	for i, v := range expected {
		token := sc.Advance()
		tok := token.Token
		lit := token.string
		switch v {
		case IDENT:
			if tok != IDENT || lit != identifiers[ident_c] {
				t.Fatalf("'%s' not identified as IDENTIFIER", lit)
			}
			ident_c += 1
		case WS:
			if tok != WS {
				t.Fatalf("'%s' not identified as WHITESPACE", lit)
			}
		case EOF:
			if tok != EOF {
				t.Fatalf("Failed to identify EOF")
			}
			if i != len(expected)-1 {
				t.Fatalf("Failed to read all")
			}
		}
	}
}

func TestWhitespace(t *testing.T) {
	input := `sug main -o program --optimization="max"`
	reader := strings.NewReader(input)
	sc := NewScanner(reader)

	identifiers := []string{
		"sug",
		"main",
		"o",
		"program",
		"optimization",
		"max",
	}
	ident_c := 0
	expected := []Token{
		IDENT,
		IDENT,
		DASH,
		IDENT,
		IDENT,
		DDASH,
		IDENT,
		ASSIGNMENT,
		STR,
		EOF,
	}

	for i, v := range expected {
		token := sc.Advance()
		tok := token.Token
		lit := token.string
		switch v {
		case IDENT:
			if tok != IDENT || lit != identifiers[ident_c] {
				t.Fatalf("'%s' not identified as IDENTIFIER", lit)
			}
			ident_c += 1
		case WS:
			if tok != WS {
				t.Fatalf("'%s' not identified as WHITESPACE", lit)
			}
		case EOF:
			if tok != EOF {
				t.Fatalf("Failed to identify EOF")
			}
			if i != len(expected)-1 {
				t.Fatalf("Failed to read all")
			}
		}
	}
}
