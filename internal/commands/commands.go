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
	"github.com/thedevscott/blogaggregator/internal/feed"
)

type Command struct {
	Name string
	Args []string
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

func MiddlewareLoggedIn(handler func(s *State, cmd Command, user database.User) error) func(*State, Command) error {
	return func(s *State, cmd Command) error {
		user, err := s.Db.GetUser(context.Background(), s.Cfg.CurrentUserName)
		if err != nil {
			return err
		}
		return handler(s, cmd, user)
	}
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

func HandlerAggregate(s *State, cmd Command) error {
	feed, err := feed.FetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return fmt.Errorf("failed to fetch feed: %w", err)
	}
	fmt.Printf("Feed: %+v\n", feed)
	return nil
}

func HandlerAddFeed(s *State, cmd Command, user database.User) error {

	if len(cmd.Args) != 2 {
		return fmt.Errorf("usage: %s <name> <url>", cmd.Name)
	}

	name := cmd.Args[0]
	url := cmd.Args[1]

	feed, err := s.Db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		Name:      name,
		Url:       url,
	})

	if err != nil {
		return fmt.Errorf("failed to create feed: %w", err)
	}

	feedFollow, err := s.Db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return fmt.Errorf("failed to create feed follow: %w", err)
	}

	fmt.Println("Feed created successfully:")
	printFeed(feed, user)
	fmt.Println("Feed followed successfully:")
	printFeedFollow(feedFollow.UserName, feedFollow.FeedName)
	fmt.Println("\n=================================")
	return nil

}

func HandlerGetFeed(s *State, cmd Command) error {
	feeds, err := s.Db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get feeds: %w", err)
	}

	if len(feeds) == 0 {
		fmt.Println("No feed found.")
		return nil
	}

	fmt.Printf("Feeds found: %d\n", len(feeds))

	for _, feed := range feeds {
		user, err := s.Db.GetUserById(context.Background(), feed.UserID)

		if err != nil {
			return fmt.Errorf("failed to get user name from DB: %w", err)
		}

		fmt.Printf("Feed Name: %s\n", feed.Name)
		fmt.Printf("Feed URL: %s\n", feed.Url)
		fmt.Printf("Feed User: %s\n", user.Name)
	}

	return nil
}

func HandlerFollow(s *State, cmd Command, user database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <url_of_feed>", cmd.Name)
	}

	url := cmd.Args[0]

	feed, err := s.Db.GetFeedByURL(context.Background(), url)
	if err != nil {
		return fmt.Errorf("faild to get feed: %w", err)
	}

	feedFollowRow, err := s.Db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})

	if err != nil {
		return fmt.Errorf("failed to create feed follow: %w", err)
	}

	fmt.Println("Feed Follow created:")
	printFeedFollow(feedFollowRow.UserName, feedFollowRow.FeedName)

	return nil
}

func HandlerFollowing(s *State, cmd Command, user database.User) error {
	feedFollows, err := s.Db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("failed to get feed follows: %w", err)
	}

	if len(feedFollows) == 0 {
		fmt.Println("No feeds followed by this user.")
		return nil
	}

	fmt.Printf("Feeds Followed by user %s:\n", user.Name)
	for _, feed := range feedFollows {
		fmt.Printf("* %s\n", feed.FeedName)
	}

	return nil
}

func printFeedFollow(username, feedname string) {
	fmt.Printf("* User:        %s\n", username)
	fmt.Printf("* Feed:        %s\n", feedname)
}

func printFeed(feed database.Feed, user database.User) {
	fmt.Printf("* ID:            %s\n", feed.ID)
	fmt.Printf("* Created:       %v\n", feed.CreatedAt)
	fmt.Printf("* Updated:       %v\n", feed.UpdatedAt)
	fmt.Printf("* Name:          %s\n", feed.Name)
	fmt.Printf("* URL:           %s\n", feed.Url)
	fmt.Printf("* User:          %s\n", user.Name)
}

func printUser(usr database.User) {
	fmt.Printf("\t-ID:\t%v\n", usr.ID)
	fmt.Printf("\t-Name:\t%v\n", usr.Name)
}
