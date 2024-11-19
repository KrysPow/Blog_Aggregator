package main

import (
	"fmt"
	"os"

	"github.com/KrysPow/go_blog_aggregator/internal/commands"
	"github.com/KrysPow/go_blog_aggregator/internal/config"
)

func main() {
	conf, err := config.Read()
	if err != nil {
		fmt.Println(err)
	}

	state := commands.State{
		Config: &conf,
	}

	cmds := commands.Commands{
		CommandMap: make(map[string]func(*commands.State, commands.Command) error),
	}

	cmds.Register("login", commands.HandlerLogin)

	args := os.Args
	if len(args) < 2 {
		fmt.Println("Command required")
		os.Exit(1)
	}

	cmd := commands.Command{
		Name: args[1],
		Args: args[2:],
	}

	err = cmds.Run(&state, cmd)
	if err != nil {
		fmt.Println(err)
	}
}
