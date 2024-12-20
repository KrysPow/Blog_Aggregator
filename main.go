package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/KrysPow/go_blog_aggregator/internal/commands"
	"github.com/KrysPow/go_blog_aggregator/internal/config"
	"github.com/KrysPow/go_blog_aggregator/internal/database"
	_ "github.com/lib/pq"
)

func main() {
	conf, err := config.Read()
	if err != nil {
		fmt.Println(err)
	}

	db, err := sql.Open("postgres", conf.DBurl)
	if err != nil {
		log.Fatal("Could not connect to database: ", err)
	}

	state := commands.State{
		DB:     database.New(db),
		Config: &conf,
	}

	cmds := commands.Commands{
		CommandMap: make(map[string]func(*commands.State, commands.Command) error),
	}

	cmds.Register("login", commands.HandlerLogin)
	cmds.Register("register", commands.HandlerRegister)
	cmds.Register("reset", commands.HandlerReset)
	cmds.Register("users", commands.HandlerUsers)
	cmds.Register("agg", commands.HandlerAgg)
	cmds.Register("addfeed", commands.MiddlewareLoggedIn(commands.HandlerAddFeed))
	cmds.Register("feeds", commands.HandlerFeeds)
	cmds.Register("follow", commands.MiddlewareLoggedIn(commands.HandlerFollow))
	cmds.Register("following", commands.MiddlewareLoggedIn(commands.HandlerFollowing))
	cmds.Register("unfollow", commands.MiddlewareLoggedIn(commands.HandlerUnfollow))
	cmds.Register("browse", commands.MiddlewareLoggedIn(commands.HandlerBrowse))

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
