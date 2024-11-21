package commands

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/KrysPow/go_blog_aggregator/internal/config"
	"github.com/KrysPow/go_blog_aggregator/internal/database"
	"github.com/google/uuid"
)

type State struct {
	DB     *database.Queries
	Config *config.Config
}

type Command struct {
	Name string
	Args []string
}

func HandlerLogin(s *State, cmd Command) error {
	if len(cmd.Args) == 0 {
		log.Fatal("login expects an argument, the username")
	}

	_, err := s.DB.GetUser(context.Background(), sql.NullString{String: cmd.Args[0], Valid: true})
	if err != nil {
		log.Fatal("User does not exist in the database")
	}

	err = s.Config.SetUser(cmd.Args[0])
	if err != nil {
		return err
	}

	fmt.Printf("User has been set to %s\n", cmd.Args[0])
	return nil
}

func HandlerRegister(s *State, cmd Command) error {
	if len(cmd.Args) == 0 {
		log.Fatal("login expects an argument, the username")
	}

	usr_param := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name: sql.NullString{
			String: cmd.Args[0],
			Valid:  true, // Indicating the value is NOT null.
		},
	}

	_, err := s.DB.CreateUser(context.Background(), usr_param)
	if err != nil {
		log.Fatal("User could not be created, ", err)
	}

	_, err = s.DB.GetUser(context.Background(), usr_param.Name)
	if err != nil {
		log.Fatal("User already exists!")
	}

	err = s.Config.SetUser(cmd.Args[0])
	if err != nil {
		return err
	}

	fmt.Printf("User %s has been registered\n", cmd.Args[0])
	return nil
}

func HandlerReset(s *State, cmd Command) error {
	s.DB.DeleteUsers(context.Background())
	fmt.Println("All users have been DELETED!")
	return nil
}

func HandlerUsers(s *State, cmd Command) error {
	users, err := s.DB.GetUsers(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	for _, user := range users {
		if user.Name.String == s.Config.CurrentUserName {
			fmt.Printf("* %s (current)\n", user.Name.String)
		} else {
			fmt.Println("* " + user.Name.String)
		}

	}
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
