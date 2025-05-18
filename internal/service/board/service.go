package board

import (
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/wehw93/kanban-board/internal/lib/http/response"
	"github.com/wehw93/kanban-board/internal/lib/jwt"
	"github.com/wehw93/kanban-board/internal/model"
	"github.com/wehw93/kanban-board/internal/storage"
	srv "github.com/wehw93/kanban-board/internal/transport/http"
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
			return "", fmt.Errorf("%s: %w", op, "Invalid credentials")
		}
		slog.Warn("failed to get user", err)
		return "", fmt.Errorf("%s: %w", op, err)
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
	user, err := s.store.User().GetByID(user_id)
	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}
	projects, err := s.store.User().GetProjects(user_id)
	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}
	tasks, err := s.store.User().GetTasks(user_id)
	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}
	resp := &response.ReadUserResponse{
		ID:       uint(user_id),
		Name:     user.Name,
		Email:    user.Email,
		Projects: make([]response.ProjectBrief, 0, len(projects)),
		Tasks:    make([]response.TaskBrief, 0, len(tasks)),
	}
	for _, p := range projects {
		resp.Projects = append(resp.Projects, response.ProjectBrief{
			ID:          uint(p.ID),
			Name:        p.Name,
			Description: p.Description,
		})
	}
	for _, t := range tasks {
		resp.Tasks = append(resp.Tasks, response.TaskBrief{
			ID:     uint(t.ID),
			Name:   t.Name,
			Status: t.Status,
		})
	}
	return resp, nil

}

func (s *Service) DeleteUser(user_id int) error {
	const op = "board.service.deleteuser"
	err := s.store.User().Delete(user_id)
	if err != nil {
		return fmt.Errorf("%s : %w", op, err)
	}
	return nil
}

func (s *Service) UpdateEmail(user model.User) error {
	const op = "board.service.updateEmail"
	err := s.store.User().UpdateEmail(&user)
	if err != nil {
		return fmt.Errorf("%s : %w", op, err)
	}
	return nil
}

func (s *Service) UpdatePassword(user model.User) error {
	const op = "board.service.updatePassword"
	err := s.store.User().UpdatePassword(&user)
	if err != nil {
		return fmt.Errorf("%s : %w", op, err)
	}
	return nil
}

func (s *Service) CreateProject(project *model.Project) error {
	const op = "service.CreateProject"
	err := s.store.Project().Create(project)
	if err != nil {
		return fmt.Errorf("%s:%w", op, err)
	}
	return nil
}

func (s *Service) ReadProject(name string) (*response.ReadProjectResponse, error) {
	const op = "board.service.ReadProject"
	project, err := s.store.Project().GetByName(name)
	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}
	tasks, err := s.store.Project().GetTasks(int(project.ID))
	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}
	resp := &response.ReadProjectResponse{
		ID:       uint(project.ID),
		Name:     project.Name,
		Description: project.Description,
	}
	for _, t := range tasks {
		resp.Tasks = append(resp.Tasks, response.TaskBrief{
			ID:     uint(t.ID),
			Name:   t.Name,
			Status: t.Status,
		})
	}
	return resp, nil

}
