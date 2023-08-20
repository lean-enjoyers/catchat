package parser

import (
	"bufio"
	"bytes"
	"io"
)

type TokenStruct struct {
	Token
	string
}

// Lexical scanner
type Scanner struct {
	r      *bufio.Reader
	tokens []TokenStruct
	pos    uint32
}

func NewScanner(r io.Reader) *Scanner {
	s := &Scanner{r: bufio.NewReader(r)}
	s.tokens = s.scanTokens()
	return s
}

// Returns the next character
func (s *Scanner) read() rune {
	ch, _, err := s.r.ReadRune()

	if err != nil {
		return eof
	}
	return ch
}

// Place the previously read rune back on the reader
func (s *Scanner) unreadRune() {
	_ = s.r.UnreadRune()
}

// Consumes the current rune if it matches what is expected
func (s *Scanner) expect(expectedRune rune) bool {
	if ch := s.read(); ch == expectedRune {
		return true
	} else {
		s.unreadRune()
		return false
	}
}

func (s *Scanner) Read() TokenStruct {
	return s.tokens[s.pos]
}

func (s *Scanner) Advance() TokenStruct {
	t := s.Read()
	s.pos += 1
	return t
}

func (s *Scanner) Peek(offset uint32) TokenStruct {
	if int(offset+s.pos) >= len(s.tokens) {
		return TokenStruct{EOF, ""}
	}

	return s.tokens[offset+s.pos]
}

func (s *Scanner) Expect(token Token) bool {
	return s.Read().Token == token
}

func (s *Scanner) OptionalConsume(token Token) {
	if s.Expect(token) {
		s.Advance()
	}
}

func (s *Scanner) scanTokens() []TokenStruct {
	var tokens []TokenStruct

	for {
		tok, lex := s.advanceRune()
		tokens = append(tokens, TokenStruct{tok, lex})

		if tok == EOF {
			break
		}
	}

	return tokens
}

// Returns the token and the string literal
func (s *Scanner) advanceRune() (token Token, literal string) {
	// Get rid of all whitespace
	ch := s.read()
	s.unreadRune()

	if isWhiteSpace(ch) {
		s.scanWhitespace()
	}

	ch = s.read()

	if isString(ch) {
		s.unreadRune()
		return s.scanIdentifier()
	} else {
		// Individual character
		switch ch {
		case eof:
			return EOF, ""
		case equal:
			return ASSIGNMENT, "="
		case dash:
			if s.expect('-') {
				return DDASH, "--"
			} else {
				return DASH, "-"
			}
		}
	}

	return ILLEGAL, string(ch)
}

func (s *Scanner) scanWhitespace() {
	// Consume all contiguous runes (as long as they are whitespace)
	for {
		if ch, _, _ := s.r.ReadRune(); ch == eof {
			break
		} else if !isWhiteSpace(ch) {
			s.unreadRune()
			break
		}
		// continue by default...
	}
}

func (s *Scanner) scanIdentifier() (token Token, literal string) {
	var buf bytes.Buffer

	firstCh := s.read()
	var hasQuote bool

	// Consume the first letter if it isn't quote
	if !isQuote(firstCh) {
		buf.WriteRune(firstCh)
	} else {
		hasQuote = true
	}

	// Consume all contiguous letter runes
	for {
		if ch := s.read(); isEOF(ch) {
			s.unreadRune()
			break
		} else if hasQuote && isQuote(ch) {
			break
		} else if !hasQuote && !isLetter(ch) && !isDigit(ch) && ch != '_' {
			// end quote or not letter/digit/_
			s.unreadRune()
			break
		} else {
			buf.WriteRune(ch)
		}
	}

	if hasQuote {
		return STR, buf.String()
	} else {
		return IDENT, buf.String()
	}
}
