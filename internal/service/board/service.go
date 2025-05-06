package board

import (
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/wehw93/kanban-board/internal/lib/jwt"
	"github.com/wehw93/kanban-board/internal/model"
	"github.com/wehw93/kanban-board/internal/storage"
	srv "github.com/wehw93/kanban-board/internal/transport/http"
	"github.com/wehw93/kanban-board/internal/transport/http/response"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	store storage.Store
}

func NewService(store storage.Store) *Service {
	return &Service{
		store: store,
	}
}



func (s *Service) LoginUser(email string, password string) (string, error) {
	const op = "board.service.Login"
	user, err := s.store.User().Login(email)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			slog.Warn("user not found", err)
			return "", fmt.Errorf("%s: %w",op,"Invalid credentials")
		}
		slog.Warn("failed to get user",err)
		return "", fmt.Errorf("%s: %w",op,err)
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Encrypted_password), []byte(password)); err != nil {
		return "", fmt.Errorf("%s : %w", op, err)
	}
	token, err := jwt.NewToken(user, time.Hour, srv.JWTSecret)
	if err != nil {
		return "", fmt.Errorf("%s : %w", op, err)
	}
	return token, nil
}

func (s *Service) CreateUser(user *model.User) error {
	const op = "service.CreateUser"
	err := s.store.User().Create(user)
	if err != nil {
		return fmt.Errorf("%s:%w", op, err)
	}
	return nil
}


func (s *Service) ReadUser(user_id int) (*response.ReadUserResponse, error) {
	const op = "board.service.ReadUser"
	user,err:=s.store.User().GetUserByID(user_id)
	if err!=nil{
		return nil,fmt.Errorf("%s:%w",op,err)
	}
	
}