package postgresql

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"github.com/wehw93/kanban-board/internal/model"
	"github.com/wehw93/kanban-board/internal/storage"
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
	project := &model.Project{
		Name: name,
	}
	err := r.store.db.QueryRow("SELECT * FROM projects WHERE name = $1", name).Scan(&project.ID, &project.Name, &project.IDCreator, &project.Description)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, storage.ErrProjectNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return project, nil
}

func (r *ProjectRepository) GetTasks(projectID int) ([]model.Task, error) {
	const op = "storage.postgresql.project.get_tasks"

	rows, err := r.store.db.Query(
		`SELECT t.id, t.name, t.description,t.status FROM tasks t 
		JOIN columns c ON t.id_column = c.id 
		WHERE c.id_project = $1`,
		projectID,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()
	var tasks []model.Task
	for rows.Next() {
		var t model.Task
		if err := rows.Scan(
			&t.ID,
			&t.Name,
			&t.Description,
			&t.Status,
		); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		tasks = append(tasks, t)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return tasks, nil
}

func (r *ProjectRepository) Delete(userID int, name string) error {
	const op = "storage.postgresql.project.delete"

	res, err := r.store.db.Exec(
		"DELETE FROM projects WHERE id_creator = $1 and name = $2",
		userID,
		name,
	)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("%s: %w", op, storage.ErrProjectNotFound)
	}

	return nil
}

func (r *ProjectRepository) UpdateName(name string, project model.Project) error {
	const op = "storage.postgresql.project.updateName"

	res, err := r.store.db.Exec("UPDATE projects SET name = $1 WHERE name = $2 and id_creator = $3", name, project.Name, project.IDCreator)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if rowsAffected == 0 {
		return storage.ErrProjectNotFound
	}
	return nil
}

func (r *ProjectRepository) UpdateDescription(project model.Project) error {
	const op = "storage.postgresql.project.UpdateDescription"

	res, err := r.store.db.Exec("UPDATE projects SET description = $1 WHERE name = $2 and id_creator = $3",
		project.Description,
		project.Name,
		project.IDCreator)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if rowsAffected == 0 {
		return storage.ErrProjectNotFound
	}
	return nil
}

func (r *ProjectRepository) ListProjects() ([]model.Project, error) {
	const op = "storage.postgresql.project.ListProjects"

	var listProjects []model.Project
	rows, err := r.store.db.Query("SELECT * FROM projects")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()
	for rows.Next() {
		var p model.Project
		if err := rows.Scan(
			&p.ID,
			&p.Name,
			&p.IDCreator,
			&p.Description,
		); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		listProjects = append(listProjects, p)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return listProjects, nil
}
