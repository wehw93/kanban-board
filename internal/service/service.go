package service

import (
	"github.com/wehw93/kanban-board/internal/lib/http/response"
	"github.com/wehw93/kanban-board/internal/model"
)

type BoardService interface {
	CreateUser(user *model.User) error
	LoginUser(email string, password string) (string, error)
	ReadUser(user_id int) (*response.ReadUserResponse, error)
	DeleteUser(user_id int) error
	UpdateEmail(user model.User) error
	UpdatePassword(user model.User) error
	CreateProject(project *model.Project) error
	ReadProject(name string) (*response.ReadProjectResponse, error)
	DeleteProject(userID int,name string) error
	UpdateProjectDescription(project model.Project)error
	UpdateProjectName(name string,project model.Project)error
	ListProjects()([]model.Project,error)
}
