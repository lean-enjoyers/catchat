package parser

// This is handwritten to only support posix style command line arguments.

import (
	"strings"
)

type Parser struct {
	Scanner Scanner
}

func NewParser(input string) *Parser {
	r := strings.NewReader(input)
	return &Parser{Scanner: *NewScanner(r)}
}
