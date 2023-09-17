package commands

type Command interface {
	execute(args CommandArgument)
}
