package models

import (
	"errors"
	"time"

	"github.com/gofrs/uuid"
)

type Todo struct {
	ID        uuid.UUID `json:"_id"`
	Name      string    `json:"name"`
	Deadline  time.Time `json:"deadline"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedBy uint32    `json:"created_by"`
	Overdue   bool      `json:"overdue"`
	Subs      []*Sub    `json:"subs"`
}

func (t *Todo) Verify() error {
	if t.Name == "" || len(t.Name) == 0 {
		return errors.New("Todo id is required")
	}
	if t.CreatedBy == 0 {
		return errors.New("Todo created_by is required")
	}
	return nil
}
