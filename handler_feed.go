package main

import (
	"context"
	"fmt"
	"time"

	"github.com/gbelintani/gator/internal/database"
	"github.com/google/uuid"
)

func handlerAgg(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("wrong args")
	}

	interval, err := time.ParseDuration(cmd.args[0])
	if err != nil {
		return fmt.Errorf("error on parsing duration: %w", err)
	}

	fmt.Printf("Collecting feeds every %v", interval)

	ticker := time.NewTicker(interval)
	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}
}

func handlerAddFeed(s *state, cmd command, usr database.User) error {
	if len(cmd.args) != 2 {
		return fmt.Errorf("wrong args for add feed")
	}

	feed, err := s.db.CreateFeed(context.Background(),
		database.CreateFeedParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name:      cmd.args[0],
			Url:       cmd.args[1],
			UserID:    usr.ID,
		})
	if err != nil {
		return fmt.Errorf("error creating feed on db: %w", err)
	}

	_, err = s.db.CreateFollow(context.Background(),
		database.CreateFollowParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			UserID:    usr.ID,
			FeedID:    feed.ID,
		})
	if err != nil {
		return fmt.Errorf("error creating follow on db: %w", err)
	}

	fmt.Printf("Feed Added: %v\n", feed)
	return nil
}

func handlerListFeed(s *state, _ command) error {
	feeds, err := s.db.ListFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("error listing feeds: %w", err)
	}

	for _, f := range feeds {
		fmt.Printf("Name: %s\n", f.Name)
		fmt.Printf("URL: %s\n", f.Url)
		fmt.Printf("User: %s\n", f.UserName)
	}

	return nil
}
