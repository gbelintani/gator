package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/gbelintani/gator/internal/config"
	"github.com/gbelintani/gator/internal/database"
	"github.com/gbelintani/gator/internal/rss"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type state struct {
	db     *database.Queries
	config *config.Config
}

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, c command) error {
		usr, err := s.db.GetUser(context.Background(), s.config.CurrentUserName)
		if err != nil {
			return fmt.Errorf("error get user: %w", err)
		}
		return handler(s, c, usr)
	}
}

func scrapeFeeds(s *state) error {
	next, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return fmt.Errorf("error get next: %w", err)
	}

	err = s.db.MarkFeedFetched(context.Background(), next.ID)
	if err != nil {
		return fmt.Errorf("error marking as fetched: %w", err)
	}

	res, err := rss.FetchFeed(context.Background(), next.Url)
	if err != nil {
		return fmt.Errorf("error fetching: %w", err)
	}

	fmt.Printf("\nSaving from: %s\n", res.Channel.Title)
	for _, item := range res.Channel.Item {
		pubDate, _ := time.Parse(time.RFC3339Nano, item.PubDate)

		post, err := s.db.CreatePost(context.Background(),
			database.CreatePostParams{
				ID:          uuid.New(),
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
				Title:       item.Title,
				Description: sql.NullString{String: item.Description, Valid: true},
				Url:         item.Link,
				PublishedAt: pubDate,
				FeedID:      next.ID,
			})
		if err != nil {
			fmt.Printf("SAVED(%v): %s\n", post.ID, item.Title)
		}
	}
	fmt.Printf("Doen Saving from: %s\n", res.Channel.Title)

	return nil
}

func main() {
	c, err := config.Read()
	if err != nil {
		panic("error reading config")
	}
	db, err := sql.Open("postgres", c.DbURL)
	if err != nil {
		panic("error openning db connection")
	}
	dbQueries := database.New(db)
	s := &state{
		config: &c,
		db:     dbQueries,
	}

	cmds := commands{
		commands: map[string]func(*state, command) error{},
	}
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerListUsers)
	cmds.register("agg", handlerAgg)
	cmds.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	cmds.register("feeds", handlerListFeed)
	cmds.register("follow", middlewareLoggedIn(handlerFollow))
	cmds.register("following", middlewareLoggedIn(handlerUserFollowing))
	cmds.register("unfollow", middlewareLoggedIn(handlerUnfollow))
	cmds.register("browse", hanlderBrowse)

	args := os.Args
	if len(args) < 2 {
		fmt.Printf("Too few args %s", args)
		os.Exit(1)
	}

	cmdS := args[1]
	cmdArgs := args[2:]

	err = cmds.run(s, command{
		name: cmdS,
		args: cmdArgs,
	})
	if err != nil {
		fmt.Printf("Error running command %s: %s", cmdS, err)
		os.Exit(1)
	}

}
