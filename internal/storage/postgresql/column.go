package postgresql

import (
	"fmt"

	"github.com/wehw93/kanban-board/internal/model"
)

type ColumnRepository struct{
	store *Storage
}

func (r*ColumnRepository) CreateColumn(column * model.Column)error{
	const op = "storage.postgresql.column.CreateColumn"

	err:=r.store.db.QueryRow("INSERT INTO columns (name,id_project) VALUES($1,$2) RETURNING id",column.Name,column.ID_project).Scan(&column.ID)
	if err!=nil{
		return fmt.Errorf("%s: %w",op,err)
	}
	return nil
}

