package main

import (
	"context"
	"fmt"
	"time"

	"github.com/gbelintani/gator/internal/database"
	"github.com/google/uuid"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("username is required")
	}

	u := cmd.args[0]
	_, err := s.db.GetUser(context.Background(), u)
	if err != nil {
		return fmt.Errorf("error get user: %w", err)
	}

	err = s.config.SetUser(u)
	if err != nil {
		return fmt.Errorf("could not set user: %w", err)
	}

	fmt.Printf("User %s set!\n", u)
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("name is required")
	}

	u := cmd.args[0]
	dbUser, err := s.db.CreateUser(context.Background(),
		database.CreateUserParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name:      u,
		})
	if err != nil {
		return fmt.Errorf("error creating user: %w", err)
	}

	s.config.SetUser(u)
	fmt.Printf("User %s(%v) created!\n", u, dbUser.ID)
	return nil
}

func handlerListUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("error listing users %w", err)
	}

	for _, u := range users {
		fmt.Printf("* %s", u.Name)
		if u.Name == s.config.CurrentUserName {
			fmt.Print(" (current)")
		}
		fmt.Println()
	}

	return nil
}
