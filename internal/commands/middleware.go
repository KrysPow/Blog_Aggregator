package commands

import (
	"context"
	"fmt"

	"github.com/KrysPow/go_blog_aggregator/internal/database"
)

func MiddlewareLoggedIn(handler func(s *State, cmd Command, user database.User) error) func(*State, Command) error {
	return func(s *State, cmd Command) error {
		cur_user, err := s.DB.GetUser(context.Background(), s.Config.CurrentUserName)
		if err != nil {
			return fmt.Errorf("failed to get user: %w", err)
		}
		err = handler(s, cmd, cur_user)
		if err != nil {
			return err
		}
		return nil
	}
}
