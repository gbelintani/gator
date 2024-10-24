package main

import (
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
