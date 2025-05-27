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
			return "", fmt.Errorf("%s: %w", op, errors.New("Invalid credentials"))
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

func (s *Service) DeleteProject(userID int,name string)error{
	const op = "board.service.DeleteProject"

	err:=s.store.Project().Delete(userID,name)
	if err!=nil{
		return fmt.Errorf("%s: %w",op,err)
	}
	return nil
}

func (s*Service)UpdateProjectName(name string,project model.Project)error{
	const op = "board.service.UpdateProjectName"

	err:=s.store.Project().UpdateName(name,project)
	if err!=nil{
		return fmt.Errorf("%s: %w",op,err)
	}
	return nil
}

func (s*Service)UpdateProjectDescription(project model.Project)error{
	const op = "board.service.UpdateProjectDescription"

	err:=s.store.Project().UpdateDescription(project)
	if err!=nil{
		return fmt.Errorf("%s: %w",op,err)
	}
	return nil
}

func (s*Service)ListProjects()([]model.Project,error){
	const op = "board.service.ListProjects"

	listProjects,err:=s.store.Project().ListProjects()
	if err!=nil{
		return nil,fmt.Errorf("%s: %w",op,err)
	}
	return listProjects,nil
}
func (s*Service)CreateColumn(column * model.Column)error{
	const op = "board.service.CreateColumn"

	err:=s.store.Column().CreateColumn(column)
	if err!=nil{
		return fmt.Errorf("%s: %w",op,err)
	}
	return nil
}

func (s*Service)ReadColumn(column model.Column) (*response.ReadColumnResponse,error){
	const op = "board.service.ReadColumn"

	id,err:=s.store.Column().GetID(column)
	if err!=nil{
		return nil,fmt.Errorf("%s: %w",op,err)
	}
	column.ID = int64(id)
	resp:=&response.ReadColumnResponse{
		ID: id,
		Name: column.Name,
	}
	tasks,err:=s.store.Column().GetTasks(column)
	if err!=nil{
		return nil,fmt.Errorf("%s: %w",op,err)
	}
	for _,t:=range tasks{
		resp.Tasks = append(resp.Tasks, response.TaskBrief{
			ID: uint(t.ID),
			Name: t.Name,
			Status: t.Status,
		})
	}
	return resp,nil
}

func (s*Service)DeleteColumn(id int)error{
	const op = "board.service.DeleteColumn"

	err:=s.store.Column().DeleteColumn(id)
	if err!=nil{
		return fmt.Errorf("%s: %w",op,err)
	}
	return nil
}

func (s*Service)UpdateColumnName(column model.Column, name string)error{
	const op = "board.service.UpdateColumnName"

	err:=s.store.Column().UpdateColumnName(column,name)
	if err!=nil{
		return fmt.Errorf("%s: %w",op,err)
	}
	return nil
}

func (s*Service)CreateTask(task *model.Task)error{
	const op = "board.service.CreateTask"

	err:=s.store.Task().CreateTask(task)
	if err!=nil{
		return fmt.Errorf("%s: %w",op,err)
	}
	return nil
}

func (s*Service)ReadTask(task *model.Task)error{
	const op = "board.service.ReadTask"

	err:=s.store.Task().ReadTask(task)
	if err!=nil{
		return fmt.Errorf("%s: %w",op,err)
	}
	return nil
}

func (s*Service)DeleteTask(IDuser int,id int)error{
	const op = "board.service.DeleteTask"

	err:=s.store.Task().DeleteTask(IDuser,id)
	if err!=nil{
		return fmt.Errorf("%s: %w",op,err)
	}
	return nil
}

func (s*Service) UpdateTaskName(task *model.Task)error{
	const op = "board.service.UpdateTaskName"

	err:=s.store.Task().UpdateTaskName(task)
	if err!=nil{
		return fmt.Errorf("%s: %w",op,err)
	}
	return nil
}

func (s*Service) UpdateTaskColumn(task *model.Task)error{
	const op = "board.service.UpdateTaskColumn"

	err:=s.store.Task().UpdateTaskColumn(task)
	if err!=nil{
		return fmt.Errorf("%s: %w",op,err)
	}
	return nil
}

func (s*Service) UpdateTaskDescription(task *model.Task)error{
	const op = "board.service.UpdateTaskDescription"

	err:=s.store.Task().UpdateTaskDescription(task)
	if err!=nil{
		return fmt.Errorf("%s: %w",op,err)
	}
	return nil
}

func (s*Service)GetLogsTask(id_task int) ([]model.Task_log,error){
	const op = "service.board.GetLogsTask"
	logs,err:=s.store.Task().GetLogsTask(id_task)
	if err!=nil{
		return nil,fmt.Errorf("%s: %w",op,err)
	}
	return logs,nil
}