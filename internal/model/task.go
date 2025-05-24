package model

import "database/sql"

type Task struct{
	ID int64
	ID_column int64
	Name string
	Description string
	Date_of_create string
	Date_of_execution sql.NullTime
	ID_executor sql.NullInt64
	ID_creator int64
	Status string
}