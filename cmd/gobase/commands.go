package main

import "errors"

type command struct {
	Name string
	Args []string
}

type commandInfo struct {
	handler     func(*state, command) error
	usage       string
	description string
}

type commands struct {
	registeredCommands map[string]commandInfo
}

func (c *commands) register(
	name string,
	handler func(*state, command) error,
	usage, description string,
) {
	c.registeredCommands[name] = commandInfo{
		handler:     handler,
		usage:       usage,
		description: description,
	}
}

func (c *commands) run(s *state, cmd command) error {
	ci, ok := c.registeredCommands[cmd.Name]
	if !ok {
		return errors.New("command not found")
	}

	return ci.handler(s, cmd)
}
