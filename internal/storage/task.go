package storage

import "github.com/wehw93/kanban-board/internal/model"

type TaskRepository interface{
	CreateTask(task * model.Task)error
}