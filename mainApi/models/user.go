package models

import (
	"errors"

	"time"
)

type User struct {
	ID         uint32    `json:"_id"`
	Firstname  string    `json:"firstname"`
	Lastname   string    `json:"lastname"`
	Password   string    `json:"password"`
	CreatedAt  time.Time `json:"created_at"`
	LastSeenAt time.Time `json:"last_seen_at"`
}

func (u *User) Verify() error {
	if u.Firstname == "" || len(u.Firstname) == 0 {
		return errors.New("User firstname is required")
	}
	if u.Lastname == "" || len(u.Lastname) == 0 {
		return errors.New("User lastname is required")
	}
	return nil
}
