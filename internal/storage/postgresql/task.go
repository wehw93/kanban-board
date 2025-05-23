package postgresql

import (
	"fmt"

	"github.com/wehw93/kanban-board/internal/model"
)

type TaskRepository struct {
	store *Storage
}

func (r *TaskRepository) CreateTask(task * model.Task)error{
	const op = "storage.postgresql.Task.CreateTask"

	err:=r.store.db.QueryRow("INSERT INTO tasks (id_column,name,description,id_creator,status,date_of_create) VALUES ($1,$2,$3,$4,$5,$6) RETURNING id",
	task.ID_column,
	task.Name,
	task.Description,
	task.ID_creator,
	task.Status,
	task.Date_of_create).Scan(&task.ID)
	if err!=nil{
		return fmt.Errorf("%s: %w",op,err)
	}
	return nil
}