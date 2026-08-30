package commands

import (
	"fmt"

	"github.com/FooNameBar/gator/internal/data"
)

func HandlerLogin(s *data.State, cmd Command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("login command requires a username\n")
	}

	name := cmd.Args[0]

	if err := s.Config.SetUser(name); err != nil {
		return err
	}

	fmt.Printf("Username %s has been set\n", name)
	return nil
}

