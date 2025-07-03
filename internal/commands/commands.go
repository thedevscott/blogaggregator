package commands

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/thedevscott/blogaggregator/internal/config"
	"github.com/thedevscott/blogaggregator/internal/database"
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
	Db  *database.Queries
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

	// Make sure the user is in the DB before setting in config json file
	_, err := s.Db.GetUser(context.Background(), name)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	err = s.Cfg.SetUser(name)
	if err != nil {
		return fmt.Errorf("faild to set current user: %w", err)
	}

	fmt.Println("User login successful!")
	return nil
}

func HandlerRegister(s *State, cmd Command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <name>", cmd.Name)
	}

	name := cmd.Args[0]

	usr, err := s.Db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
	})

	if err != nil {
		return fmt.Errorf("error creating user: %w", err)
	}

	err = s.Cfg.SetUser(usr.Name)
	if err != nil {
		return fmt.Errorf("faild to set current user: %w", err)
	}

	fmt.Printf("Crated user: %s\n", usr.Name)
	printUser(usr)
	return nil
}

func HandlerResetUsers(s *State, cmd Command) error {
	err := s.Db.ResetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("failed to reset users table: %w", err)
	}
	fmt.Println("Users table reset successful!")
	return nil
}

func HandlerGetUsers(s *State, cmd Command) error {
	users, err := s.Db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get users: %w", err)
	}
	fmt.Println("Users table reset successful!")
	for _, user := range users {
		if strings.ToLower(user.Name) != strings.ToLower(s.Cfg.CurrentUserName) {
			fmt.Printf("* %s\n", user.Name)
		} else {
			fmt.Printf("* %s (current)\n", user.Name)
		}
	}
	return nil
}

func printUser(usr database.User) {
	fmt.Printf("\t-ID:\t%v\n", usr.ID)
	fmt.Printf("\t-Name:\t%v\n", usr.Name)
}
