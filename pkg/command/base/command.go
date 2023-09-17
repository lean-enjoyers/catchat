package command

import (
	command_args "github.com/lean-enjoyers/catchat/pkg/command/args"
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
	Execute(args command_args.CommandArgument, hub IHub)
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

func GetArgs(cmd string) command_args.CommandArgument {
	p := parser.NewParser(cmd)
	return p.Parse()
}
