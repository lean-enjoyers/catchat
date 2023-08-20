package parser

// This is handwritten to only support posix style command line arguments.

import (
	"strings"

	"github.com/lean-enjoyers/catchat/pkg/commands"
)

type Parser struct {
	scanner Scanner
}

func NewParser(input string) *Parser {
	r := strings.NewReader(input)
	return &Parser{scanner: *NewScanner(r)}
}

func (p *Parser) Parse() commands.CommandArgument {
	args := commands.NewCommandArgument()

	// Retrieve the command name.
	if p.scanner.Expect(IDENT) {
		token := p.scanner.Advance()
		args.SetCommand(token.string)
	} else {
		return *args
	}

	// Parse arguments.
	for {
		token := p.scanner.Advance()

		if token.Token == EOF {
			return *args
		}

		switch token.Token {
		case DASH:

			// Retrieve the flag name.
			var flagname string
			if p.scanner.Expect(IDENT) {
				flagname = p.scanner.Advance().string
			} else {
				return *args
			}

			// Incase of =
			p.scanner.OptionalConsume(ASSIGNMENT)

			valuename := ""

			// Retrieve the flag value
			if p.scanner.Expect(IDENT) || p.scanner.Expect(STR) {
				// Flag value
				valuename = p.scanner.Advance().string

			}
			args.SetShortOption(flagname, valuename)

		case DDASH:

			// Retrieve the flag name.
			var flagname string
			if p.scanner.Expect(IDENT) {
				flagname = p.scanner.Advance().string
			} else {
				return *args
			}

			// Incase of =
			p.scanner.OptionalConsume(ASSIGNMENT)

			valuename := ""

			// Retrieve the flag value
			if p.scanner.Expect(IDENT) || p.scanner.Expect(STR) {
				// Flag value
				valuename = p.scanner.Advance().string
			}

			args.SetLongOption(flagname, valuename)
		}
	}
}
