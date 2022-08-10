package database

import (
	// "context"
	"context"
	"errors"
	"mainApi/models"
	"mainApi/utils"

	_ "github.com/lib/pq"
)

type UserDB interface {
	CreateUser(ctx context.Context, user *models.User) error
	GetUserByCredentials(ctx context.Context, firstname string) (*models.User, error)
}

var ErrUserExists = errors.New("could not create user")

func (d *database) CreateUser(ctx context.Context, user *models.User) error {
	// Get a tx for making transaction requests.
	tx, err := d.conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	sql := "insert into users (firstname, lastname, password) values ($1, $2, $3) returning id"
	{
		stmt, err := tx.PrepareContext(ctx, sql)
		if err != nil {
			tx.Rollback()
			return err
		}
		hashedPassword, err := utils.HashPassword(user.Password)
		if err != nil {
			tx.Rollback()
			return err
		}
		defer stmt.Close()
		if err != nil {
			tx.Rollback()
			return err
		}
		var userId uint32
		err = stmt.QueryRowContext(ctx, user.Firstname, user.Lastname, hashedPassword).Scan(&userId)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	if err = tx.Commit(); err != nil {
		return err
	}
	return err
}

func (d *database) GetUserByCredentials(ctx context.Context, firstname string) (*models.User, error) {
	sql := "SELECT * FROM users WHERE firstname = $1"
	rows, err := d.conn.QueryContext(ctx, sql, firstname)
	if err != nil {
		return &models.User{}, err
	}
	var user models.User
	for rows.Next() {
		err := rows.Scan(&user.Firstname, &user.Lastname, &user.Password, &user.CreatedAt, &user.LastSeenAt)
		if err != nil {
			return &models.User{}, err
		}
	}
	return &user, nil
}
