package model

type Task struct{
	id int64
	id_column int64
	name string
	description string
	date_of_create string
	date_of_execution string
	ID_executor int64
	id_creator int64
}