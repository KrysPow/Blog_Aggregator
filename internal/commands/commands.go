package commands

import (
	"fmt"
	"os"

	"github.com/KrysPow/go_blog_aggregator/internal/config"
)

type State struct {
	Config *config.Config
}

type Command struct {
	Name string
	Args []string
}

func HandlerLogin(s *State, cmd Command) error {
	if len(cmd.Args) == 0 {
		fmt.Errorf("login expects an argument, the username")
		os.Exit(1)
	}

	err := s.Config.SetUser(cmd.Args[0])
	if err != nil {
		return err
	}

	fmt.Printf("User has been set to %s\n", cmd.Args[0])
	return nil
}

type Commands struct {
	CommandMap map[string]func(*State, Command) error
}

func (c *Commands) Register(name string, f func(*State, Command) error) {
	if _, ok := c.CommandMap[name]; ok {
		fmt.Printf("%s is already registered\n", name)
	}
	c.CommandMap[name] = f
}

func (c *Commands) Run(s *State, cmd Command) error {
	err := c.CommandMap[cmd.Name](s, cmd)
	if err != nil {
		return err
	}
	return nil
}
