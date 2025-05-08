package storage

import "github.com/wehw93/kanban-board/internal/model"

type UserRepository interface {
	Create(u *model.User) error
	Login(email string) (model.User, error)
	GetByID(user_id int) (model.User, error)
	GetProjects(user_id int) ([]model.Project, error)
	GetTasks(user_id int) ([]model.Task, error)
}
