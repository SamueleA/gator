package middleware

import (
	"context"

	"github.com/samuelea/gator/internal/config"
	"github.com/samuelea/gator/internal/database"
)

func MiddlewareLoggedIn(handler func(s *config.State, cmd config.Command, user database.User) error) func(*config.State, config.Command) error {
	return func(s *config.State, cmd config.Command) error {
		user, err := s.DbQueries.GetUser(context.Background(), s.Config.CurrentUserName)
		if err != nil {
			return err
		}
		
		return handler(s, cmd, user)
	}
}