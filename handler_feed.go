package main

import (
	"context"
	"fmt"
	"time"

	"github.com/gbelintani/gator/internal/database"
	"github.com/gbelintani/gator/internal/rss"
	"github.com/google/uuid"
)

func handlerAgg(s *state, _ command) error {
	feed, err := rss.FetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return fmt.Errorf("could not fetch feed: %w", err)
	}

	fmt.Printf("Feed found: %v\n", feed)
	return nil
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
