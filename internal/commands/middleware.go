package commands

import (
	"context"
	"fmt"

	"github.com/FooNameBar/gator/internal/data"
	"github.com/FooNameBar/gator/internal/database"
)

func MiddlewareLoggedIn(handler func(s *data.State, cmd Command, user database.User) error) func(s *data.State, cmd Command) error {
	return func(s *data.State, cmd Command) error {
		name := s.Config.CurrentUserName
		user, err := s.DB.GetUser(context.Background(), name)
		if err != nil {
			return fmt.Errorf("DB.GetUser: user %s not found: %v\n", name, err)
		}

		return handler(s, cmd, user)
	}
}
