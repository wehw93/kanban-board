package model

import "database/sql"

type Task struct{
	ID int64
	ID_column int64
	Name string
	Description string
	Date_of_create string
	Date_of_execution sql.NullTime `json:"date_of_execution" swaggertype:"string" format:"date-time"`
	ID_executor sql.NullInt64 `json:"id_executor" swaggertype:"integer"`
	ID_creator int64
	Status string
}