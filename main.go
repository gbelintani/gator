package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/gbelintani/gator/internal/config"
	"github.com/gbelintani/gator/internal/database"
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
