package main

import (
	"fmt"
	"os"

	"www.github.com/a-fleming/blog-aggregator/internal/config"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("gator: error: the following arguments are required: command")
		os.Exit(1)
	}
	commandName := os.Args[1]
	args := os.Args[2:]

	var cfg config.Config
	cfg, err := config.Read()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}

	cliState := state{
		config: &cfg,
	}

	cmds := GetCommands()

	cmdToRun := command{
		name:      commandName,
		arguments: args,
	}

	err = cmds.run(&cliState, cmdToRun)
	if err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}

	cfg, err = config.Read()

	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}
}
