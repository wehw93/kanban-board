package postgresql

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/wehw93/kanban-board/internal/storage"
)

type Storage struct {
	db                  *sql.DB
	userrepository      *UserRepository
	taskrepository      *TaskRepository
	projectRepository   *ProjectRepository
	columnRepository    *ColumnRepository
	task_log_Repository *Task_log_Repository
}

func New(dsn string) (*Storage, error) {
	const op = "storage.postgresql.new"

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}
	return &Storage{db: db}, nil
}

func (s *Storage) Column() storage.ColumnRepository {
	if s.columnRepository != nil {
		return s.columnRepository
	}
	s.columnRepository = &ColumnRepository{
		store: s,
	}
	return s.columnRepository
}

func (s *Storage) Project() storage.ProjectRepository {
	if s.projectRepository != nil {
		return s.projectRepository
	}
	s.projectRepository = &ProjectRepository{
		store: s,
	}
	return s.projectRepository
}

func (s *Storage) Task_log() storage.Task_log_Repository {
	if s.task_log_Repository != nil {
		return s.task_log_Repository
	}
	s.task_log_Repository = &Task_log_Repository{
		store: s,
	}
	return s.task_log_Repository
}

func (s *Storage) User() storage.UserRepository {
	if s.userrepository != nil {
		return s.userrepository
	}
	s.userrepository = &UserRepository{
		store: s,
	}
	return s.userrepository
}
func (s *Storage) Task() storage.TaskRepository {
	if s.taskrepository != nil {
		return s.taskrepository
	}
	s.taskrepository = &TaskRepository{
		store: s,
	}
	return s.taskrepository
}

func (s *Storage) Close() {
	s.db.Close()
}
