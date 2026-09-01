package commands

import (
	"context"
	"fmt"
	"time"

	"github.com/FooNameBar/gator/internal/data"
	"github.com/FooNameBar/gator/internal/database"
	"github.com/google/uuid"
)

func HandlerLogin(s *data.State, cmd Command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("login command requires a username\n")
	}

	name := cmd.Args[0]

	_, err := s.DB.GetUser(context.Background(), name)
	if err != nil {
		return fmt.Errorf("No user found with name '%s'\n", name)
	}

	if err := s.Config.SetUser(name); err != nil {
		return err
	}

	fmt.Printf("%s logged in\n", name)
	return nil
}

func HandlerRegister(s *data.State, cmd Command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("register command requires a name\n")
	}

	name := cmd.Args[0]

	user, err := s.DB.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
	})

	if err != nil {
		return fmt.Errorf("DB.CreateUser(): with name: %s %v\n", name, err)
	}

	if err := s.Config.SetUser(name); err != nil {
		return err
	}

	fmt.Printf("User %s created\n", name)

	fmt.Printf("user: %+v\n", user)
	return nil
}

func HandlerReset(s *data.State, cmd Command) error {
	err := s.DB.Reset(context.Background())
	if err != nil {
		return fmt.Errorf("DB.Reset(): %v\n", err)
	}

	fmt.Println("Database reset!")
	return nil
}

func GetUsers(s *data.State, cmd Command) error {
	users, err := s.DB.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("DB.GetUsers(): %v\n", err)
	}

	for _, u := range users {
		fmt.Printf("* %s", u.Name)
		if u.Name == s.Config.CurrentUserName {
			fmt.Print(" (current)")
		}
		fmt.Println()
	}
	return nil
}
