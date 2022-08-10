package models

import (
	"errors"

	"github.com/gofrs/uuid"
)

type Sub struct {
	ID        uuid.UUID `json:"_id"`
	Name      string    `json:"name"`
	TodoId    uuid.UUID `json:"todo_id"`
	Completed bool      `json:"completed"`
}

var (
	ErrTodoSubNotFound = errors.New("Todo sub not found")
)

func (s *Sub) Verify() error {
	if s.Name == "" || len(s.Name) == 0 {
		return errors.New("Sub id is required")
	}
	if s.TodoId == uuid.Nil {
		return errors.New("Sub todo_id is required")
	}
	return nil
}
