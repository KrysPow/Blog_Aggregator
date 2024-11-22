package commands

import (
	"context"
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

	_, err := s.DB.GetUser(context.Background(), cmd.Args[0])
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
		Name:      cmd.Args[0],
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
		if user.Name == s.Config.CurrentUserName {
			fmt.Printf("* %s (current)\n", user.Name)
		} else {
			fmt.Println("* " + user.Name)
		}

	}
	return nil
}

func HandlerAgg(s *State, cmd Command) error {
	url := "https://www.wagslane.dev/index.xml"

	rss_feed, err := FetchFeed(context.Background(), url)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print(*rss_feed)
	return nil
}

func HandlerAddFeed(s *State, cmd Command, user database.User) error {
	if len(cmd.Args) != 2 {
		log.Fatal("You need 2 arguments, the name and the url")
	}

	feed, err := s.DB.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.Args[0],
		Url:       cmd.Args[1],
		UserID:    user.ID,
	})
	if err != nil {
		log.Fatal("Feed could not be created: ", err)
	}

	_, err = s.DB.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		log.Fatal("Feed follow could not be created: ", err)
	}

	fmt.Println(feed)
	return nil
}

func HandlerFeeds(s *State, cmd Command) error {
	feeds_data, err := s.DB.GetFeedsNamesUrlsUserName(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	for _, feed := range feeds_data {
		fmt.Println(feed.Name)
		fmt.Println(feed.Url)
		fmt.Println(feed.Name_2.String)
	}
	return nil
}

func HandlerFollow(s *State, cmd Command, user database.User) error {
	if len(cmd.Args) != 1 {
		log.Fatal("Follow need one argument, the URL")
	}

	feed, err := s.DB.GetFeedByUrl(context.Background(), cmd.Args[0])
	if err != nil {
		log.Fatal(err)
	}

	feed_follow, err := s.DB.CreateFeedFollow(context.Background(),
		database.CreateFeedFollowParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			UserID:    user.ID,
			FeedID:    feed.ID,
		})
	if err != nil {
		log.Fatal("feed_follow error: ", err)
	}

	fmt.Println(feed_follow.FeedName)
	fmt.Println(feed_follow.UserName)
	return nil
}

func HandlerFollowing(s *State, cmd Command, user database.User) error {
	if len(cmd.Args) != 0 {
		log.Fatal("Following does not need any argument")
	}

	feeds, err := s.DB.GetFeedFollowForUser(context.Background(), user.ID)
	if err != nil {
		log.Fatal("Feeds could not be querried: ", err)
	}

	for _, feed := range feeds {
		fmt.Println(feed.FeedName)
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
