package commands

import (
	"errors"
	"fmt"
	"strings"

	"github.com/thedevscott/blogaggregator/internal/config"
)

func cleanInput(text string) []string {
	ltext := strings.TrimSpace(strings.ToLower(text))
	splitText := strings.Fields(ltext)

	return splitText
}

type Command struct {
	Name        string
	description string
	Args        []string
}

type Commands struct {
	RegisteredCommands map[string]func(*State, Command) error
}

type State struct {
	Cfg *config.Config
}

func (c *Commands) Register(name string, f func(*State, Command) error) {
	c.RegisteredCommands[name] = f
}

func (c *Commands) Run(s *State, cmd Command) error {
	f, ok := c.RegisteredCommands[cmd.Name]

	if !ok {
		return errors.New("command not found")
	}

	return f(s, cmd)
}

func HandlerLogin(s *State, cmd Command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <name>", cmd.Name)
	}

	name := cmd.Args[0]

	err := s.Cfg.SetUser(name)
	if err != nil {
		return fmt.Errorf("faild to set current user: %w", err)
	}

	fmt.Println("User login successful!")
	return nil
}
