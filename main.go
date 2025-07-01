package main

import (
	"log"
	"os"

	"github.com/thedevscott/blogaggregator/internal/commands"
	"github.com/thedevscott/blogaggregator/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}

	programState := &commands.State{
		Cfg: &cfg,
	}

	cmds := commands.Commands{
		RegisteredCommands: make(map[string]func(*commands.State, commands.Command) error),
	}

	cmds.Register("login", commands.HandlerLogin)

	if len(os.Args) < 2 {
		log.Fatal("Usage: cli <command> [args...]")
	}

	cmdName := os.Args[1]
	cmdArgs := os.Args[2:]

	err = cmds.Run(programState, commands.Command{Name: cmdName, Args: cmdArgs})

	if err != nil {
		log.Fatal(err)
	}
}
