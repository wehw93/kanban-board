package service

import "github.com/wehw93/kanban-board/internal/model"

type BoardService interface {
	CreateUser(user *model.User) error
}
