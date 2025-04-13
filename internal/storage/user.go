package storage

import "github.com/wehw93/kanban-board/internal/model"

type UserRepository interface {
	Create(u *model.User) error
}
