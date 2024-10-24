package main

import (
	"fmt"
	"os"

	"github.com/gbelintani/gator/internal/config"
)

type state struct {
	config *config.Config
}

func main() {
	c, err := config.Read()
	if err != nil {
		panic("error reading config")
	}
	s := state{
		config: &c,
	}

	cmds := commands{
		commands: map[string]func(*state, command) error{},
	}
	cmds.register("login", handlerLogin)

	args := os.Args
	if len(args) < 2 {
		fmt.Printf("Too few args %s", args)
		os.Exit(1)
	}

	cmdS := args[1]
	cmdArgs := args[2:]

	err = cmds.run(&s, command{
		name: cmdS,
		args: cmdArgs,
	})
	if err != nil {
		fmt.Printf("Error running command %s: %s", cmdS, err)
		os.Exit(1)
	}

}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("username is required")
	}

	u := cmd.args[0]
	err := s.config.SetUser(u)
	if err != nil {
		return fmt.Errorf("could not set user: %w", err)
	}

	fmt.Printf("User %s set!\n", u)
	return nil
}
