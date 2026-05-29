package main

import (
	"fmt"
	"errors"
	"github.com/SlyShamrock/Gator/internal/config"
	"time"
	"github.com/SlyShamrock/Gator/internal/database"
	"github.com/google/uuid"
	"context"
	"os"
	"database/sql"
)

type state struct {
	db *database.Queries
	cfg config.Config
}

type command struct {
	name string
	args []string	
}

type commands struct {
	handlers map[string]func(*state, command) error
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) != 1  {
		return errors.New("username is required\n")
	}
	
	username := cmd.args[0]
	
	_, err := s.db.GetUser(context.Background(), username)
	if err == sql.ErrNoRows {
		return fmt.Errorf("username not found: %s", err)
		os.Exit(1)
	}
	
	if err != nil {
		return fmt.Errorf("failed to get user: %s", err)
	}			

	err = s.cfg.SetUser(username)
	if err != nil {
		return fmt.Errorf("failed to set user: %s", err)
	}
	fmt.Printf("user has been set to : %s\n", s.cfg.CurrentUserName)	
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return errors.New("name is required\n")
	}
	
	user := cmd.args[0]
	now := time.Now()
	newID := uuid.New()

	params := database.CreateUserParams{
		ID: newID,
		CreatedAt: now,
		UpdatedAt: now,
		Name: user,
	}

	u, err := s.db.CreateUser(context.Background(), params)
	if err != nil {
		return fmt.Errorf("failed to create user: %s", err)
	}

	err = s.cfg.SetUser(u.Name)
	if err != nil {
		return fmt.Errorf("failed to set user: %s", err)
	}
	fmt.Printf("new user created: %s\n", s.cfg.CurrentUserName)
	return nil
}

func handlerReset(s *state, cmd command) error {
	err := s.db.DeleteUsers(context.Background())
	if err != nil {
		return fmt.Errorf("failed to delete users from table: %s", err)		
	}
	fmt.Println("successfully deleted all users")
	return nil		
}

func handlerUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		fmt.Errorf("failed to get all users: %s", err)		
	}
	for _, user := range users {
		if user.Name == s.cfg.CurrentUserName {
			fmt.Printf("* %s (current)\n", user.Name)
		} else {		
			fmt.Printf("* %s\n", user.Name)
		}
	}	
	return nil
}

func handlerAgg(s *state, cmd command) error {
	url := "https://www.wagslane.dev/index.xml"
	feed, err := fetchFeed(context.Background(), url)
	if err != nil {
		return fmt.Errorf("failed to fetch from provided url: %s", err)
	}
	fmt.Printf("%+v\n", feed)
	return nil
}

func handlerAddFeed(s *state, cmd command) error {
	if len(cmd.args) != 2 {
		return errors.New("feed name and url require\n")
	}
	username := s.cfg.CurrentUserName
	userData, err := s.db.GetUser(context.Background(), username)
	if err == sql.ErrNoRows {
		return fmt.Errorf("username not found: %s", err)		
	}	
	
	now := time.Now()
	feedName := cmd.args[0]
	feedUrl := cmd.args[1]
	newId := uuid.New()
	
	feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID: newId,
		CreatedAt: now,
		UpdatedAt: now,
		Name: feedName,
		Url: feedUrl,
		UserID: userData.ID,
	})
	if err != nil {
		return fmt.Errorf("failed to create feed: %s", err)
	}
	fmt.Printf("%+v\n", feed)
	return nil
}

func handlerFeeds(s *state, cmd command) error {
	feeds, err := s.db.DisplayFeeds(context.Background())
	if err != nil {
			return fmt.Errorf("failed to retrieve feeds: %s", err)
		}
	for _, feed := range feeds {		
		fmt.Println(feed.FeedName)
		fmt.Println(feed.FeedUrl)
		fmt.Println(feed.UserName)
	}
	return nil
}

func (c *commands) run(s *state, cmd command) error {
	value, ok := c.handlers[cmd.name]
	if !ok {
		return errors.New("command not found")
	}
	return value(s, cmd)	
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.handlers[name] = f
}

