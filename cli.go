package main

import (
	"fmt"

	"www.github.com/a-fleming/blog-aggregator/internal/config"
)

type state struct {
	config *config.Config
}

type command struct {
	name      string
	arguments []string
}

type commands struct {
	cliCommands map[string]func(*state, command) error
}

func GetCommands() commands {
	cmds := commands{
		cliCommands: map[string]func(*state, command) error{},
	}
	cmds.register("login", handlerLogin)
	return cmds
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.arguments) == 0 {
		return fmt.Errorf("gator login: error: the following arguments are required: username")
	}
	userName := cmd.arguments[0]
	err := s.config.SetUser(userName)
	if err != nil {
		return err
	}
	fmt.Printf("username '%s' has been set\n", userName)
	return nil
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.cliCommands[name] = f
}

func (c *commands) run(s *state, cmd command) error {
	cmdFunc, exists := c.cliCommands[cmd.name]
	if !exists {
		return fmt.Errorf("Unknown command '%s'\n", cmd.name)
	}

	return cmdFunc(s, cmd)
}
