package database

import (
	"github.com/jmoiron/sqlx"
)

type Database interface {
	UserDB
	TodoDB
}

type database struct {
	conn *sqlx.DB
}

func (d *database) Close() error {
	return d.conn.Close()
}
