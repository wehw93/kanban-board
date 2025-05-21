package storage

import "github.com/wehw93/kanban-board/internal/model"

type ColumnRepository interface{
	CreateColumn(column * model.Column)error
	GetID(column model.Column)(int,error)	
	GetTasks(column model.Column)([]model.Task,error)
}