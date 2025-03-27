package board

import "github.com/wehw93/kanban-board/internal/storage"

type Service struct{
	store storage.Store
}

func NewService(store storage.Store) *Service{
	return &Service{
		store:  store,
	}
}