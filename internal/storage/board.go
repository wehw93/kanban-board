package storage

import "github.com/wehw93/kanban-board/internal/model"

type ProjectRepository interface {
	Create(project *model.Project) error
}
