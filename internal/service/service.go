package service

import "github.com/wehw93/kanban-board/internal/model"

type BoardService interface {
	CreateUser() (*model.User, error)
}
