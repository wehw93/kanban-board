package postgresql

import (
	"fmt"

	"github.com/wehw93/kanban-board/internal/model"
	"github.com/wehw93/kanban-board/internal/storage"
)

type ColumnRepository struct {
	store *Storage
}

func (r *ColumnRepository) CreateColumn(column *model.Column) error {
	const op = "storage.postgresql.column.CreateColumn"

	err := r.store.db.QueryRow("INSERT INTO columns (name,id_project) VALUES($1,$2) RETURNING id", column.Name, column.ID_project).Scan(&column.ID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (r *ColumnRepository) GetID(column model.Column) (int, error) {
	const op = "storage.postgresql.column.GetID"
	var id int
	err := r.store.db.QueryRow("SELECT id FROM columns WHERE name = $1 and id_project = $2", column.Name, column.ID_project).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return id, nil
}

func (r *ColumnRepository) GetTasks(column model.Column) ([]model.Task, error) {
	const op = "storage.postgresql.column.GetTasks"

	rows, err := r.store.db.Query("SELECT id,name,description,status FROM columns WHERE id_column = $1", column.ID)
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

func (r *ColumnRepository) DeleteColumn(id int) error {
	const op = "storage.postgesql.column.DeleteColumn"
	res, err := r.store.db.Exec("DELETE FROM column WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("%s: %w", op, storage.ErrColumnNotFound)
	}
	return nil
}

func (r*ColumnRepository)UpdateColumnName(column model.Column,name string)error{
	const op = "storage.postgreqsql.column,UpdateColumnName"
	res,err:=r.store.db.Exec("UPDATE columns SET name = $1 WHERE id = $2",name,column.ID)
	if err!=nil{
		return fmt.Errorf("%s: %w",op,err)
	}
	rowsAffectd,err:=res.RowsAffected()
	if err!=nil{
		return fmt.Errorf("%s: %w",op,err)
	}
	if rowsAffectd==0{
		return storage.ErrColumnNotFound
	}
	return nil
}