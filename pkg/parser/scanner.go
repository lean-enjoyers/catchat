package parser

import (
	"bufio"
	"bytes"
	"io"
)

// Lexical scanner
type Scanner struct {
	r *bufio.Reader
}

func NewScanner(r io.Reader) *Scanner {
	return &Scanner{r: bufio.NewReader(r)}
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
func (s *Scanner) unread() {
	_ = s.r.UnreadRune()
}

// Consumes the current rune if it matches what is expected
func (s *Scanner) expect(expectedRune rune) bool {
	if ch := s.read(); ch == expectedRune {
		return true
	}
	s.unread()
	return false
}

// Returns the token and the string literal
func (s *Scanner) Scan() (token Token, literal string) {
	// Read the next rune
	ch := s.read()

	if isWhiteSpace(ch) {
		return s.scanWhitespace()
	} else if isString(ch) {
		s.unread()
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

func (s *Scanner) scanWhitespace() (token Token, literal string) {
	// Consume all contiguous runes (as long as they are whitespace)
	for {
		if ch := s.read(); ch == eof {
			break
		} else if !isWhiteSpace(ch) {
			s.unread()
			break
		}
		// continue by default...
	}

	return WS, " "
}

func (s *Scanner) scanIdentifier() (token Token, literal string) {
	var buf bytes.Buffer

	firstCh := s.read()

	// Consume the first letter if it isn't quote
	if !isQuote(firstCh) {
		buf.WriteRune(firstCh)
	}

	// Consume all contiguous letter runes
	for {
		if ch := s.read(); isEOF(ch) || isQuote(ch) {
			break
		} else if !isLetter(ch) && !isDigit(ch) && ch != '_' {
			// end quote or not letter/digit/_
			s.unread()
			break
		} else {
			buf.WriteRune(ch)
		}
	}
	return IDENT, buf.String()
}
