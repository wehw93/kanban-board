package service

import (
	"github.com/wehw93/kanban-board/internal/lib/http/response"
	"github.com/wehw93/kanban-board/internal/model"
)

type BoardService interface {
	CreateUser(user *model.User) error
	LoginUser(email string, password string) (string, error)
	ReadUser(user_id int) (*response.ReadUserResponse, error)
}
