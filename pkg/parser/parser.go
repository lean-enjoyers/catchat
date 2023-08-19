package parser

// This is handwritten to only support posix style command line arguments.

import (
	"log"
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
	var args commands.CommandArgument

	// Retrieve the command name.
	if p.scanner.Expect(IDENT) {
		_, cmd := p.scanner.Advance()
		args.SetCommand(cmd)
	} else {
		log.Println("Invalid command name, not specified.")
		return args
	}

	// Parse arguments.
	for {
		token, _ := p.scanner.Advance()

		if token == EOL {
			return args
		}

		switch token {
		case DASH:

			// Retrieve the flag name.
			var flagname string
			if p.scanner.Expect(IDENT) {
				_, flagname = p.scanner.Advance()
			} else {
				log.Println("Error: expected flag.")
				return args
			}

			// Incase of =
			p.scanner.OptionalConsume(ASSIGNMENT)

			// Retrieve the flag value
			if p.scanner.Expect(IDENT) {
				// Flag value
				_, valuename := p.scanner.Advance()

				args.SetShortOption(flagname, valuename)
			}
		}
	}
}
