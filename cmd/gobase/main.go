package main

import (
	"fmt"
	"os"
)

type state struct {
	cfg *Config
}

func main() {
	cfg, err := ReadConfig()
	if err != nil {
		fmt.Fprintln(os.Stdout, err)
		os.Exit(1)
	}

	programState := &state{
		cfg: &cfg,
	}

	cmds := &commands{
		registeredCommands: make(map[string]commandInfo),
	}

	cmds.register("help", func(programState *state, cmd command) error {
		return handlerHelp(cmds, programState, cmd)
	}, "help", "Displays all available commands and their usage")
	cmds.register("init", handlerInit, "init", "Initiate migration")
	cmds.register("migrate", handlerMigrate, "migrate", "Migrate to database")

	if len(os.Args) < 2 {
		fmt.Println("Usage: gobase <command>")
		return
	}

	cmdName := os.Args[1]
	cmdArgs := os.Args[2:]

	err = cmds.run(programState, command{Name: cmdName, Args: cmdArgs})
	if err != nil {
		fmt.Fprintln(os.Stdout, err)
		os.Exit(1)
	}
}
