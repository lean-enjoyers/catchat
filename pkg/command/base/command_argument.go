package command

type CommandArgument struct {
	cmd   string
	args  []string
	lflag map[string]string
	sflag map[string]string
}

func NewCommandArgument() *CommandArgument {
	return &CommandArgument{
		cmd:   "",
		args:  make([]string, 0),
		lflag: make(map[string]string),
		sflag: make(map[string]string),
	}
}

func (c *CommandArgument) SetCommand(name string) {
	c.cmd = name
}

func (c *CommandArgument) SetLongOption(key string, value string) {
	c.lflag[key] = value
}

func (c *CommandArgument) SetShortOption(key string, value string) {
	c.sflag[key] = value
}

func (c *CommandArgument) AddArgument(arg string) {
	c.args = append(c.args, arg)
}

func (c *CommandArgument) GetArguments() []string {
	return c.args
}

func (c *CommandArgument) GetFlag(key string) (string, bool) {
	val1, found1 := c.lflag[key]
	val2, found2 := c.sflag[key]

	if found1 {
		return val1, true
	}

	if found2 {
		return val2, true
	}

	return "", false
}

func (c *CommandArgument) GetCommand() string {
	return c.cmd
}
