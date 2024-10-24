package main

import (
	"context"
	"fmt"
	"time"

	"github.com/gbelintani/gator/internal/database"
	"github.com/google/uuid"
)

func handlerFollow(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("wrong args for follow")
	}

	usr, err := s.db.GetUser(context.Background(), s.config.CurrentUserName)
	if err != nil {
		return fmt.Errorf("error get user: %w", err)
	}

	url := cmd.args[0]
	feed, err := s.db.GetByUrl(context.Background(), url)
	if err != nil {
		return fmt.Errorf("error get feed: %w", err)
	}

	follow, err := s.db.CreateFollow(context.Background(),
		database.CreateFollowParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			UserID:    usr.ID,
			FeedID:    feed.ID,
		})

	if err != nil {
		return fmt.Errorf("error creating follow: %w", err)
	}

	fmt.Printf("Feed followed: %s by user:%s\n", follow.FeedName, follow.UserName)

	return nil
}

func handlerUserFollowing(s *state, cmd command) error {
	usr, err := s.db.GetUser(context.Background(), s.config.CurrentUserName)
	if err != nil {
		return fmt.Errorf("error get user: %w", err)
	}

	follows, err := s.db.GetFeedFollowsForUser(context.Background(), usr.ID)
	if err != nil {
		return fmt.Errorf("error get follows: %w", err)
	}

	fmt.Printf("Feeds user %s is following:\n", usr.Name)
	for _, f := range follows {
		fmt.Printf(" - %s\n", f.FeedName)
	}

	return nil
}
