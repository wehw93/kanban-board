package model

type Project struct {
	ID          int64  `json:"id" db:"id"`
	Name        string `json:"name" db:"name"`
	IDCreator   int64  `json:"id_creator" db:"id_creator"`
	Description string `json:"description" db:"description"`
}
