package command

import (
	"github.com/lean-enjoyers/catchat/pkg/parser"
)

type IHub interface {
	BroadcastToClient(payload []byte, targetUserID string)
	BroadcastMessage(message string)
	HandleCommand(cmd string)
	Run()
	HandleMessage(message string)
}

type Command interface {
	Execute(args CommandArgument, hub IHub)
	Name() string
}

type CommandMap struct {
	commands map[string]Command
}

var Commands *CommandMap

func init() {
	Commands = &CommandMap{commands: make(map[string]Command)}
}

func RegisterCommand(command Command) {
	Commands.commands[command.Name()] = command
}

func (c *CommandMap) Get(cmd string) Command {
	return c.commands[cmd]
}

func GetArgs(cmd string) CommandArgument {
	p := parser.NewParser(cmd)
	args := NewCommandArgument()

	// Retrieve the command name.
	if p.Scanner.Expect(parser.IDENT) {
		token := p.Scanner.Advance()
		args.SetCommand(token.Val)
	} else {
		return *args
	}

	// Parse arguments.
	for {
		token := p.Scanner.Advance()

		if token.Tok == parser.EOF {
			return *args
		}

		switch token.Tok {

		case parser.DASH:

			// Retrieve the flag name.
			var flagname string
			if p.Scanner.Expect(parser.IDENT) {
				flagname = p.Scanner.Advance().Val
			} else {
				return *args
			}

			// Incase of =
			p.Scanner.OptionalConsume(parser.ASSIGNMENT)

			valuename := ""

			// Retrieve the flag value
			if p.Scanner.Expect(parser.IDENT) || p.Scanner.Expect(parser.STR) {
				// Flag value
				valuename = p.Scanner.Advance().Val

			}
			args.SetShortOption(flagname, valuename)

		case parser.DDASH:

			// Retrieve the flag name.
			var flagname string
			if p.Scanner.Expect(parser.IDENT) {
				flagname = p.Scanner.Advance().Val
			} else {
				return *args
			}

			// Incase of =
			p.Scanner.OptionalConsume(parser.ASSIGNMENT)

			valuename := ""

			// Retrieve the flag value
			if p.Scanner.Expect(parser.IDENT) || p.Scanner.Expect(parser.STR) {
				// Flag value
				valuename = p.Scanner.Advance().Val
			}

			args.SetLongOption(flagname, valuename)

		case parser.IDENT:
			args.AddArgument(token.Val)

		case parser.STR:
			args.AddArgument(token.Val)
		}
	}
}
