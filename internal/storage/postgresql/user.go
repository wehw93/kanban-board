package postgresql

import (
	"fmt"
	"log/slog"

	"github.com/wehw93/kanban-board/internal/model"
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
	const op = "storage.postgresql.user.User"
	res := r.store.db.QueryRow("SELECT id,name,encrypted_password WHERE email = ?", email)
	var user model.User
	err := res.Scan(&user.ID, &user.Name, &user.Encrypted_password)
	if err != nil {
		return model.User{}, err
	}
	user.Email = email
	return user, nil
}
