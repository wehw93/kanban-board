package model

import (
	"time"

)

type User struct {
	ID                 int
	name               string
	email              string
	encrypted_password string
}

type Column struct {
	ID         int
	name       string
	ID_project int
}

type Project struct {
	ID          int
	name        string
	id_creator  int
	description string
}

type Log struct {
	ID                int
	ID_task           int
	date_of_operation time.Time
	info string
}

 type Task struct{
	ID int
	ID_column int
	name string
	description string
	date_of_create time.Time
	date_of_execution time.Time
	ID_executor int
	id_creator int
	status string
 }


