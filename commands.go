package main
import (
	"fmt"
	"errors"
	"github.com/SlyShamrock/Gator/internal/config"
)

type state struct {
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
	
	err := s.cfg.SetUser(username)
	if err != nil {
		return fmt.Errorf("failed to set user: %s", err)
	}
	fmt.Printf("user has been set to : %s\n", s.cfg.CurrentUserName)
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