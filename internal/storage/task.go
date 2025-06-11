package storage

import "github.com/wehw93/kanban-board/internal/model"

type TaskRepository interface {
	CreateTask(task *model.Task) error
	ReadTask(task *model.Task) error
	DeleteTask(IDuser int, id int) error
	UpdateTaskName(task *model.Task) error
	UpdateTaskDescription(task *model.Task) error
	UpdateTaskColumn(task *model.Task) error
	GetLogsTask(id_task int) ([]model.Task_log, error)
}
