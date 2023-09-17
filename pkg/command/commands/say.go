package commands

import (
	command_args "github.com/lean-enjoyers/catchat/pkg/command/args"
	command "github.com/lean-enjoyers/catchat/pkg/command/base"
)

type SayCommand struct{}

func (s *SayCommand) Execute(args command_args.CommandArgument, hub command.IHub) {
	arguments := args.GetArguments()

	ok := len(arguments) > 0
	for _, v := range args.GetArguments() {
		hub.BroadcastMessage(v)
	}

	msg, ok1 := args.GetFlag("message")

	if ok1 {
		hub.BroadcastMessage(msg)
	}

	msg1, ok2 := args.GetFlag("m")

	if ok2 {
		hub.BroadcastMessage(msg1)
	}

	// Neither specified
	if !(ok || ok1 || ok2) {
		hub.BroadcastMessage("Say Error: No message.")
	}
}

func (s *SayCommand) Name() string {
	return "say"
}

func init() {
	sayCommand := &SayCommand{}
	command.RegisterCommand(sayCommand)
}
