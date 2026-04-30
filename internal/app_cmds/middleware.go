package app_cmds

import (
	"gator/internal/database"
	"context"
	"fmt"
)

func MiddlewareLoggedIn(handler func(s *State, cmd Command, user database.User) error) func(*State, Command) error {
	return func(s *State, cmd Command) error {
		user_name := s.ConfigPtr.Current_user_name

		ctx := context.Background()
		user, err := s.DbQPtr.GetUser(ctx, user_name)
		if fmt.Sprintf("%v", err) == fmt.Sprintf("sql: no rows in result set") {
			return fmt.Errorf("no current user found")
		} else if err != nil {
			return err
		}
		return handler(s, cmd, user)
	}
}