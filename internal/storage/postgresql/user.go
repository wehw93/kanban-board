package postgresql

import (
	"fmt"

	"github.com/wehw93/kanban-board/internal/model"
)

type UserRepository struct{
	store *Storage
}

func (r *UserRepository) Create(u*model.User) error{
	const op = "storage.postgresql.user.create"
	err:=r.store.db.QueryRow("INSERT INTO users (name,email, encrypted_password) VALUES ($1,$2,$3) RETURNING id",
		u.Name,
		u.Email,
		u.Encrypted_password,
	).Scan(&u.ID)
	
	if err!=nil{
		return fmt.Errorf("%s:%w",op,err)
	}
	return nil
}