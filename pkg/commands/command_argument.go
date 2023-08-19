package commands

type CommandArgument struct {
	cmd   string
	largs map[string]string
	sargs map[string]string
}

func (c *CommandArgument) SetCommand(name string) {
	c.cmd = name
}

func (c *CommandArgument) SetLongOption(key string, value string) {
	c.largs[key] = value
}

func (c *CommandArgument) SetShortOption(key string, value string) {
	c.sargs[key] = value
}

func (c *CommandArgument) GetFlag(key string) (string, bool) {
	val1, found1 := c.largs[key]
	val2, found2 := c.sargs[key]

	if found1 {
		return val1, true
	}

	if found2 {
		return val2, true
	}

	return "", false
}
