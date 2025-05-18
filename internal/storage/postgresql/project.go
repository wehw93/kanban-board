package postgresql

import (
	"fmt"
	"log/slog"

	"github.com/wehw93/kanban-board/internal/model"
)

type ProjectRepository struct {
	store *Storage
}

func (r *ProjectRepository) Create(project *model.Project) error {
	const op = "storage.postgresql.user.create"

	err := r.store.db.QueryRow(
		"INSERT INTO projects (name,id_creator,description) VALUES ($1, $2,$3) returning ID",
		project.Name,
		project.IDCreator,
		project.Description,
	).Scan(&project.ID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	slog.Info("project created", slog.Int64("id", int64(project.ID)))
	return nil
}
func (r *ProjectRepository) GetByName(name string) (*model.Project, error) {
	const op = "storage.postgresql.project.getbyname"

	rows, err := r.store.db.Query("SELECT * FROM projects WHERE name = $1", name)
	if err != nil {
		return nil,fmt.Errorf("%s: %w", op, err)
	}
}
