package storage

import "github.com/wehw93/kanban-board/internal/model"

type ProjectRepository interface {
	Create(project *model.Project) error
	GetByName(name string) (*model.Project, error)
	GetTasks(projectID int) ([]model.Task, error)
	Delete(userID int, name string) error
	UpdateName(name string, project model.Project) error
	UpdateDescription(project model.Project) error
	ListProjects() ([]model.Project, error)
}
