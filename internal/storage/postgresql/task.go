package postgresql

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/wehw93/kanban-board/internal/model"
	"github.com/wehw93/kanban-board/internal/storage"
)

type TaskRepository struct {
	store *Storage
}

func (r *TaskRepository) CreateTask(task *model.Task) error {
	const op = "storage.postgresql.Task.CreateTask"

	err := r.store.db.QueryRow("INSERT INTO tasks (id_column,name,description,id_creator,status,date_of_create) VALUES ($1,$2,$3,$4,$5,$6) RETURNING id",
		task.ID_column,
		task.Name,
		task.Description,
		task.ID_creator,
		task.Status,
		task.Date_of_create).Scan(&task.ID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (r *TaskRepository) ReadTask(task *model.Task) error {
	const op = "storage.postgresql.Task.ReadTask"
	err := r.store.db.QueryRow("SELECT * FROM tasks WHERE id = $1", task.ID).Scan(
		&task.ID,
		&task.ID_column,
		&task.Name,
		&task.Description,
		&task.Date_of_create,
		&task.Date_of_execution,
		&task.ID_executor,
		&task.ID_creator,
		&task.Status,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("%s: %w", op, storage.ErrTaskNotFound)
		}
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (r *TaskRepository) DeleteTask(IDuser int, id int) error {
	const op = "storage.postgresql.Task.DeleteTask"

	res, err := r.store.db.Exec("DELETE FROM tasks WHERE id = $1 and id_creator = $2", id, IDuser)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("%s: %w", op, storage.ErrTaskNotFound)
	}
	return nil
}

func (r *TaskRepository) UpdateTaskName(task *model.Task) error {
	const op = "storage.postgresql.Task.UpdateTaskName"

	res, err := r.store.db.Exec("UPDATE tasks SET name = $1 WHERE id = $2", task.Name, task.ID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if rowsAffected == 0 {
		return storage.ErrTaskNotFound
	}
	return nil
}

func (r *TaskRepository) UpdateTaskDescription(task *model.Task) error {
	const op = "storage.postgresql.Task.UpdateTaskDescription"

	res, err := r.store.db.Exec("UPDATE tasks SET description = $1 WHERE id = $2", task.Description, task.ID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if rowsAffected == 0 {
		return storage.ErrTaskNotFound
	}
	return nil
}
func (r *TaskRepository) UpdateTaskColumn(task *model.Task) error {
	const op = "storage.postgresql.Task.UpdateTaskColumn"
	var newStatus string
	err := r.store.db.QueryRow("SELECT name FROM columns WHERE id = $1", task.ID_column).Scan(&newStatus)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if newStatus == "in_progress" {
		res, err := r.store.db.Exec("UPDATE tasks SET id_column = $1 and id_executor = $2  and status = $3 WHERE id = $4",
			task.ID_column,
			task.Description,
			"in_progress",
			task.ID)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
		rowsAffected, err := res.RowsAffected()
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
		if rowsAffected == 0 {
			return storage.ErrTaskNotFound
		}
	} else {
		task.Date_of_execution = sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		}
		res, err := r.store.db.Exec("UPDATE tasks SET id_column = $1 and id_executor = $2 and status = $3 and date_of_execution = $4 WHERE id = $5",
			task.ID_column,
			task.Description,
			"done",
			task.Date_of_execution,
			task.ID)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
		rowsAffected, err := res.RowsAffected()
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
		if rowsAffected == 0 {
			return storage.ErrTaskNotFound
		}
	}

	return nil
}
