package postgresql

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/wehw93/kanban-board/internal/model"
	"github.com/wehw93/kanban-board/internal/storage"
)

type TaskRepository struct {
	store *Storage
}

const (
	todo = "todo"
	inProgress = "in_progress"
	done = "done"
)

func (r *TaskRepository) CreateTask(task *model.Task) error {
	const op = "storage.postgresql.Task.CreateTask"

	err := r.store.db.QueryRow("INSERT INTO tasks (id_column,name,description,id_creator,status,date_of_create) VALUES ($1,$2,$3,$4,$5,$6) RETURNING id",
		task.ID_column,
		task.Name,
		task.Description,
		task.ID_creator,
		todo,
		task.Date_of_create).Scan(&task.ID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	err=r.logging(int(task.ID),"create task")
	if err!=nil{
		return fmt.Errorf("%s: %w",op,err)
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
	newIdColumn:=task.ID_column
	err := r.store.db.QueryRow("SELECT id_column FROM tasks WHERE id = $1", task.ID).Scan(&task.ID_column)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	var idProject int64
	err = r.store.db.QueryRow(`
		SELECT id_project
		FROM columns
		WHERE id = $1
	`, task.ID_column).Scan(&idProject)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	var newStatus string
	var inProgressColumnID int
	err = r.store.db.QueryRow(`
			SELECT id
			FROM columns
			WHERE id_project = $1
			and name = $2
		`, idProject,
		inProgress).Scan(&inProgressColumnID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	var doneColumnID int
	err = r.store.db.QueryRow(`
		SELECT id
		FROM columns
		WHERE id_project = $1
		and name = $2
	`, idProject,
		done).Scan(&doneColumnID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	switch {
	case newIdColumn == int64(inProgressColumnID):
		newStatus = inProgress
		slog.Info("NEW STATUS: ",newStatus, inProgressColumnID)
		err=r.logging(int(task.ID),"switch status from " + todo + "to "+ newStatus)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
	case newIdColumn == int64(doneColumnID):
		newStatus = done
		task.Date_of_execution = sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		}
		slog.Info("NEW STATUS: ",newStatus, doneColumnID)
		err =r.logging(int(task.ID),"switch status from " + inProgress + "to "+ newStatus)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
	default:
		newStatus = "todo"
		task.Date_of_execution = sql.NullTime{Valid: false}
		slog.Info("NEW STATUS: ",newStatus, newIdColumn)
		err= r.logging(int(task.ID),"switch status from " + task.Status + "to "+ newStatus)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
	}
	res, err := r.store.db.Exec(`UPDATE tasks 
	SET status = $1, 
	id_column = $2, 
	date_of_execution = $3 
	WHERE id = $4`, 
		newStatus,
		newIdColumn,
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

	return nil
}

func (r*TaskRepository)logging (id_task int, info string) error{
	const op = "storage.postgres.Task.logging"
	var id int
	err:=r.store.db.QueryRow(`
		INSERT INTO logs 
		(id_task,date_of_operation,info) 
		VALUES($1,$2,$3) 
		RETURNING id
	`,id_task,
	sql.NullTime{
		Time:  time.Now(),
		Valid: true,
	},info).Scan(&id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (r*TaskRepository) GetLogsTask(id_task int)([]model.Task_log,error){
	const op = "storage.Postgresql.Task.GetLogsTask"

	var logs []model.Task_log
	rows,err:=r.store.db.Query("SELECT * FROM logs WHERE id_task = $1",id_task)
	if err!=nil{
		return nil,fmt.Errorf("%s; %w",op,err)
	}
	defer rows.Close()
	for rows.Next(){
		var l model.Task_log
		if err=rows.Scan(
			&l.ID,
			&l.ID_Task,
			&l.Date_of_operation,
			&l.Info,
			);err!=nil{
			return nil,fmt.Errorf("%s: %w",op,err)
		}
		logs = append(logs, l)
	}
	if err = rows.Err();err!=nil{
		return nil,fmt.Errorf("%s: %w",op,err)
	}
	return logs,nil
}