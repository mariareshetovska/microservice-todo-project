package database

import (
	"context"
	"fmt"

	"mainApi/models"

	"github.com/gofrs/uuid"
)

type TodoDB interface {
	CreateTodo(ctx context.Context, todo *models.Todo) error
	UpdateTodoByID(ctx context.Context, todo models.Todo) (bool, error)
	GetTodoListByUser(ctx context.Context, user_id uint32) ([]*models.Todo, error)
	GetTodoById(ctx context.Context, todoID uuid.UUID) (*models.Todo, error)
	DeleteTodo(ctx context.Context, todoID uuid.UUID) (int64, error)
}

func (d *database) CreateTodo(ctx context.Context, todo *models.Todo) error {
	// Get a tx for making transaction requests.
	tx, err := d.conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	sql := "INSERT INTO todo (name, created_by, overdue) VALUES ($1, $2, $3) RETURNING (id)"
	{
		stmt, err := tx.PrepareContext(ctx, sql)
		if err != nil {
			tx.Rollback()
			return err
		}
		defer stmt.Close()
		if err != nil {
			tx.Rollback()
			return err
		}
		var todoId uuid.UUID
		err = stmt.QueryRowContext(ctx, todo.Name, todo.CreatedBy, todo.Overdue).Scan(&todoId)
		if err != nil {
			tx.Rollback()
			return err
		}
		todo.ID = todoId
		for _, sub := range todo.Subs {
			sql := "insert into sub (name, todo_id, completed) values ($1, $2, $3) returning id"
			{
				stmt, err := tx.PrepareContext(ctx, sql)
				if err != nil {
					tx.Rollback()
					return err
				}
				defer stmt.Close()
				if err != nil {
					tx.Rollback()
					return err
				}
				err = stmt.QueryRowContext(ctx, sub.Name, todoId, sub.Completed).Scan(&sub.ID)
				if err != nil {
					if err2 := tx.Rollback(); err2 != nil {
						return err2
					}
					return err
				}
			}
			sub.TodoId = todoId
		}
		if err = tx.Commit(); err != nil {
			return err
		}
	}
	return err
}

func (d *database) GetTodoListByUser(ctx context.Context, userID uint32) ([]*models.Todo, error) {
	sql := `SELECT * FROM todo WHERE created_by = $1`
	rows, err := d.conn.QueryContext(ctx, sql, userID)
	if err != nil {
		return nil, err
	}
	var todoList []*models.Todo
	for rows.Next() {
		var todo models.Todo
		err := rows.Scan(&todo.ID, &todo.Name, &todo.Deadline, &todo.CreatedAt, &todo.UpdatedAt, &todo.CreatedBy, &todo.Overdue)
		if err != nil {
			return nil, err
		}
		todoList = append(todoList, &todo)

	}
	for _, todo := range todoList {
		sql_sub := `select * from sub where todo_id = $1`
		rs, err := d.conn.QueryContext(ctx, sql_sub, todo.ID)
		if err != nil {
			return nil, err
		}
		for rs.Next() {
			var sub models.Sub
			err := rs.Scan(&sub.ID, &sub.Name, &sub.TodoId, &sub.Completed)
			if err != nil {
				return nil, err
			}
			todo.Subs = append(todo.Subs, &sub)
		}
	}
	return todoList, nil
}

func (d *database) GetTodoById(ctx context.Context, todoID uuid.UUID) (*models.Todo, error) {
	sql := "SELECT * FROM todo WHERE id = $1"
	rows, err := d.conn.QueryContext(ctx, sql, todoID)
	if err != nil {
		return &models.Todo{}, err
	}
	var todo *models.Todo
	for rows.Next() {
		err := rows.Scan(&todo.ID, &todo.Name, &todo.Deadline, &todo.CreatedAt, &todo.UpdatedAt, &todo.CreatedBy, &todo.Overdue)
		if err != nil {
			return &models.Todo{}, err
		}
	}

	sql_sub := `SELECT * FROM sub WHERE todo_id = $1`
	rs_sub, nil := d.conn.QueryContext(ctx, sql_sub, todoID)
	if err != nil {
		return &models.Todo{}, err
	}
	defer rs_sub.Close()
	for rs_sub.Next() {
		var sub models.Sub
		err := rs_sub.Scan(&sub.ID, &sub.Name, &sub.TodoId, &sub.Completed)
		if err != nil {
			return &models.Todo{}, err
		}
		todo.Subs = append(todo.Subs, &sub)
	}
	return todo, nil
}

func (d *database) UpdateTodoByID(ctx context.Context, todo models.Todo) (bool, error) {
	tx, err := d.conn.BeginTx(ctx, nil)
	if err != nil {
		return false, err
	}
	sql := "update todo set name = $1, overdue =$2 where id = $3"
	stmt, err := tx.PrepareContext(ctx, sql)
	if err != nil {
		return false, err
	}
	defer stmt.Close()
	rs, err := stmt.ExecContext(ctx, todo.Name, todo.Overdue, todo.ID)
	if err != nil {
		return false, err
	}

	sqlDel := "delete from sub where todo_id = $1"
	stmtc, err := d.conn.PrepareContext(ctx, sqlDel)
	if err != nil {
		return false, err
	}
	defer stmt.Close()
	rsd, err := stmtc.ExecContext(ctx, todo.ID)
	if err != nil {
		return false, err
	}
	rsd.RowsAffected()

	for _, sub := range todo.Subs {
		sql := "insert into sub (name, todo_id, completed) values ($1, $2, $3) returning id"
		{
			stmt, err := tx.PrepareContext(ctx, sql)
			if err != nil {
				tx.Rollback()
				return false, err
			}
			defer stmt.Close()
			if err != nil {
				tx.Rollback()
				return false, err
			}
			err = stmt.QueryRowContext(ctx, sub.Name, todo.ID, sub.Completed).Scan(&sub.ID)
			if err != nil {
				if err2 := tx.Rollback(); err2 != nil {
					return false, fmt.Errorf("%v; %w", err, err2)
				}
				return false, err
			}

		}
	}
	rs.RowsAffected()
	return true, tx.Commit()
}

func (d *database) DeleteTodo(ctx context.Context, todoID uuid.UUID) (int64, error) {
	sql := "delete from todo where id = $1"
	stmt, err := d.conn.PrepareContext(ctx, sql)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()
	rs, err := stmt.ExecContext(ctx, todoID)
	if err != nil {
		return 0, err
	}
	return rs.RowsAffected()
}
