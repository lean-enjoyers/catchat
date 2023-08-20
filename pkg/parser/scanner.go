package parser

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
)

// Lexical scanner
type Scanner struct {
	r              *bufio.Reader
	offset         uint32
	lastWordLength uint32
}

func NewScanner(r io.Reader) *Scanner {
	return &Scanner{r: bufio.NewReader(r), offset: 0, lastWordLength: 0}
}

// Returns the next character
func (s *Scanner) read() rune {
	ch, _, err := s.r.ReadRune()
	s.offset += 1

	if err != nil {
		return eof
	}
	return ch
}

// Place the previously read rune back on the reader
func (s *Scanner) unread() {
	_ = s.r.UnreadRune()
	s.offset -= 1
}

func (s *Scanner) unreadMany(count uint32) {
	fmt.Printf("count: %d\n", count)
	for count != 0 {
		fmt.Println("UNREAD")
		s.unread()
		count -= 1
	}
}

// Consumes the current rune if it matches what is expected
func (s *Scanner) expect(expectedRune rune) bool {
	if ch := s.read(); ch == expectedRune {
		return true
	} else {
		s.unread()
		return false
	}
}

func (s *Scanner) Read() (token Token, literal string) {
	oldOffset := s.offset
	tok, lit := s.Advance()
	s.unreadMany(s.offset - oldOffset)

	return tok, lit
}

func (s *Scanner) Expect(token Token) bool {
	tok, l := s.Read()

	fmt.Printf("Token: %s\n", l)

	tok, l = s.Read()
	fmt.Printf("Token: %s\n", l)

	if tok == token {
		return true
	} else {
		log.Printf("Expected %d token code, found %d\n", token, tok)
		return false
	}
}

func (s *Scanner) OptionalConsume(token Token) {
	if s.Expect(token) {
		s.Advance()
	}
}

// Returns the token and the string literal
func (s *Scanner) Advance() (token Token, literal string) {
	// Read the next rune
	ch := s.read()

	if isWhiteSpace(ch) {
		s.scanWhitespace()
	} else {
		s.unread()
	}

	ch = s.read()

	if isString(ch) {
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

func (s *Scanner) scanWhitespace() {
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
