package postgresql

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"github.com/wehw93/kanban-board/internal/model"
	"github.com/wehw93/kanban-board/internal/storage"
)

type UserRepository struct {
	store *Storage
}

func (r *UserRepository) Create(u *model.User) error {
	const op = "storage.postgresql.user.create"
	err := r.store.db.QueryRow("INSERT INTO users (name,email, encrypted_password) VALUES ($1,$2,$3) RETURNING id",
		u.Name,
		u.Email,
		u.Encrypted_password,
	).Scan(&u.ID)
	slog.Info("returning id:", string(u.ID))
	if err != nil {
		return fmt.Errorf("%s:%w", op, err)
	}
	return nil
}

func (r *UserRepository) Login(email string) (model.User, error) {
	const op = "storage.postgresql.user.Login"
	res := r.store.db.QueryRow("SELECT id,name,encrypted_password FROM users WHERE email = $1", email)
	var user model.User
	err := res.Scan(&user.ID, &user.Name, &user.Encrypted_password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.User{}, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)

		}
		return model.User{}, fmt.Errorf("%s: %w", op, err)

	}
	user.Email = email
	return user, nil
}

func (r *UserRepository) Delete(user_id int) error {
	const op = "storage.postgresql.user.delete"
	res, err := r.store.db.Exec("DELETE FROM users WHERE id = $1", user_id)
	if err != nil {
		return fmt.Errorf("%s:%w", op, err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s:%w", op, err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("%s:%w", op, storage.ErrUserNotFound)
	}
	return nil
}

func (r *UserRepository) GetByID(user_id int) (model.User, error) {
	const op = "storage.postgresql.user.getuserbyid"
	res := r.store.db.QueryRow("SELECT name,email FROM users WHERE id = $1", user_id)
	var user model.User
	err := res.Scan(&user.Name, &user.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.User{}, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)

		}
		return model.User{}, fmt.Errorf("%s: %w", op, err)
	}
	return user, nil
}

func (r *UserRepository) GetProjects(user_id int) ([]model.Project, error) {
	const op = "storage.postgresql.user.GetProjects"

	var projects []model.Project

	rows, err := r.store.db.Query(`
        SELECT id, name, description 
        FROM projects 
        WHERE id_creator = $1`,
		user_id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	for rows.Next() {
		var p model.Project
		if err := rows.Scan(
			&p.ID,
			&p.Name,
			&p.Description,
		); err != nil {
			return nil, fmt.Errorf("%s: scan error: %w", op, err)
		}
		projects = append(projects, p)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: rows error: %w", op, err)
	}

	return projects, nil
}
func (r *UserRepository) GetTasks(user_id int) ([]model.Task, error) {
	const op = "storage.postgresql.user.getusertasks"
	var tasks []model.Task
	rows, err := r.store.db.Query("SELECT id,name,description,status FROM tasks WHERE id_executor = $1", user_id)
	defer rows.Close()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	for rows.Next() {
		var t model.Task
		if err := rows.Scan(&t.ID, &t.Name, &t.Description, &t.Status); err != nil {
			return nil, fmt.Errorf("%s:%w", op, err)
		}
		tasks = append(tasks, t)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}
	return tasks, nil
}

func (r *UserRepository) UpdatePassword(u *model.User) error {
	const op = "storage.postgresql.user.UpdatePassword"
	res, err := r.store.db.Exec("UPDATE users SET encrypted_password = $1 WHERE id = $2", u.Encrypted_password, u.ID)
	if err != nil {
		return fmt.Errorf("%s:%w", op, err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s:%w", op, err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("%s:%w", op, storage.ErrUserNotFound)
	}
	return nil
}


func (r *UserRepository) UpdateEmail(u *model.User) error {
	const op = "storage.postgresql.user.UpdateEmail"
	res, err := r.store.db.Exec("UPDATE users SET email = $1 WHERE id = $2", u.Email, u.ID)
	if err != nil {
		return fmt.Errorf("%s:%w", op, err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s:%w", op, err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("%s:%w", op, storage.ErrUserNotFound)
	}
	return nil
}