package parser

// This is handwritten to only support posix style command line arguments.

import (
	"strings"

	"github.com/lean-enjoyers/catchat/pkg/command/args"
)

type Parser struct {
	Scanner Scanner
}

func NewParser(input string) *Parser {
	r := strings.NewReader(input)
	return &Parser{Scanner: *NewScanner(r)}
}

func (p *Parser) Parse() command_args.CommandArgument {
	args := command_args.NewCommandArgument()

	// Retrieve the command name.
	if p.Scanner.Expect(IDENT) {
		token := p.Scanner.Advance()
		args.SetCommand(token.Val)
	} else {
		return *args
	}

	// Parse arguments.
	for {
		token := p.Scanner.Advance()

		if token.Tok == EOF {
			return *args
		}

		switch token.Tok {

		case DASH:

			// Retrieve the flag name.
			var flagname string
			if p.Scanner.Expect(IDENT) {
				flagname = p.Scanner.Advance().Val
			} else {
				return *args
			}

			// Incase of =
			p.Scanner.OptionalConsume(ASSIGNMENT)

			valuename := ""

			// Retrieve the flag value
			if p.Scanner.Expect(IDENT) || p.Scanner.Expect(STR) {
				// Flag value
				valuename = p.Scanner.Advance().Val

			}
			args.SetShortOption(flagname, valuename)

		case DDASH:

			// Retrieve the flag name.
			var flagname string
			if p.Scanner.Expect(IDENT) {
				flagname = p.Scanner.Advance().Val
			} else {
				return *args
			}

			// Incase of =
			p.Scanner.OptionalConsume(ASSIGNMENT)

			valuename := ""

			// Retrieve the flag value
			if p.Scanner.Expect(IDENT) || p.Scanner.Expect(STR) {
				// Flag value
				valuename = p.Scanner.Advance().Val
			}

			args.SetLongOption(flagname, valuename)

		case IDENT:
			args.AddArgument(token.Val)

		case STR:
			args.AddArgument(token.Val)
		}
	}
}
