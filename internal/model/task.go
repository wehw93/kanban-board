package model

type Task struct{
	ID int64
	ID_column int64
	Name string
	Description string
	Date_of_create string
	Date_of_execution string
	ID_executor int64
	ID_creator int64
	Status string
}