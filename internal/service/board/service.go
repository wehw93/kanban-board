package board

import (
	"fmt"

	"github.com/wehw93/kanban-board/internal/model"
	"github.com/wehw93/kanban-board/internal/storage"
)

type Service struct {
	store storage.Store
}

func NewService(store storage.Store) *Service {
	return &Service{
		store: store,
	}
}

func (s *Service) CreateUser() (*model.User, error) {
	const op = "service.CreateUser"
	user := &model.User{}
	err := s.store.User().Create(user)
	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}
	return user, nil
}
