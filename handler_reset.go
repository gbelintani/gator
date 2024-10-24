package main

import (
	"context"
	"fmt"
)

func handlerReset(s *state, _ command) error {
	err := s.db.CleanUsers(context.Background())
	if err != nil {
		return fmt.Errorf("could not reset users: %w", err)
	}
	return nil
}
