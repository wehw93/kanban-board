package storage

import "github.com/wehw93/kanban-board/internal/model"

type ColumnRepository interface{
	CreateColumn(column * model.Column)error	
}