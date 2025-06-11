package storage

import "errors"

type Store interface {
	User() UserRepository
	Project() ProjectRepository
	Column() ColumnRepository
	Task_log() Task_log_Repository
	Task() TaskRepository
}

var (
	ErrUserExists      = errors.New("user already exists")
	ErrUserNotFound    = errors.New("user not found")
	ErrProjectNotFound = errors.New("project not found")
	ErrColumnNotFound  = errors.New("column not found")
	ErrTaskNotFound    = errors.New("task not found")
)
