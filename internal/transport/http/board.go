package http

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/render"
	"github.com/wehw93/kanban-board/internal/lib/http/response"
	"github.com/wehw93/kanban-board/internal/lib/logger/sl"
	"github.com/wehw93/kanban-board/internal/model"
)

type CreateProjectRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
}

// CreateProject godoc
// @Summary Создать новый проект
// @Description Создает новый проект для текущего пользователя
// @Tags Projects
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param input body CreateProjectRequest true "Данные проекта"
// @Success 201 {object} response.SuccessResponse{data=model.Project} "Проект успешно создан"
// @Failure 400 {object} response.ErrorResponse "Неверный формат запроса"
// @Failure 401 {object} response.ErrorResponse "Не авторизован"
// @Failure 422 {object} response.ErrorResponse "Ошибка при создании проекта"
// @Router /api/projects [post]
func (s *Server) CreateProject() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		const op = "transport.http.CreateProject"

		log := s.Logger.With(slog.String("op", op))

		userID, ok := r.Context().Value("userID").(int)
		if !ok {
			log.Error("failed to get userID from context")
			render.JSON(w, r, response.ErrorResponse{
				Status:  http.StatusInternalServerError,
				Message: "Internal server error",
			})
			return
		}

		var req CreateProjectRequest

		if err := render.DecodeJSON(r.Body, &req); err != nil {

			log.Error("failed to decode request body", sl.Err(err))

			render.JSON(w, r, response.ErrorResponse{
				Status:  http.StatusBadRequest,
				Message: "Invalid request body",
			})
			return
		}

		log.Info("create project request",
			slog.Int("user_id", userID),
			slog.String("project_name", req.Name),
		)

		project := &model.Project{
			IDCreator:   int64(userID),
			Name:        req.Name,
			Description: req.Description,
		}

		if err := s.Svc.CreateProject(project); err != nil {
			log.Error("failed to create project",
				sl.Err(err),
			)
			render.JSON(w, r, response.ErrorResponse{
				Status:  http.StatusUnprocessableEntity,
				Message: "Failed to create project",
			})
			return
		}

		log.Info("project created successfully",
			slog.Int("user_id", userID),
			slog.Int64("project_id", project.ID),
		)

		render.JSON(w, r, response.SuccessResponse{
			Status: http.StatusCreated,
			Data:   project,
		})
	}
}

type ReadProjectRequest struct {
	Name string `json:"name" validate:"required"`
}

// ReadProject godoc
// @Summary Получить проект по имени
// @Description Возвращает информацию о проекте по его названию
// @Tags Projects
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param input body ReadProjectRequest true "Имя проекта"
// @Success 200 {object} response.SuccessResponse{data=model.Project} "Успешный запрос"
// @Failure 400 {object} response.ErrorResponse "Неверный формат запроса"
// @Failure 404 {object} response.ErrorResponse "Проект не найден"
// @Failure 500 {object} response.ErrorResponse "Внутренняя ошибка сервера"
// @Router /api/projects/read [post]
func (s *Server) ReadProject() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		const op = "http.ReadProject"

		log := s.Logger.With(slog.String("op", op))

		var req ReadProjectRequest

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request", sl.Err(err))
			render.JSON(w, r, response.ErrorResponse{
				Status:  http.StatusBadRequest,
				Message: "Invalid request body",
			})
			return
		}

		log.Info("reading data of project", slog.String("name", req.Name))

		resp, err := s.Svc.ReadProject(req.Name)
		if err != nil {
			log.Error("failed to read project", slog.String("name", req.Name), sl.Err(err))
			render.JSON(w, r, response.ErrorResponse{
				Status:  http.StatusInternalServerError,
				Message: "Project not found",
			})
			return
		}

		render.JSON(w, r, response.SuccessResponse{
			Status: http.StatusOK,
			Data:   resp,
		})
	}
}

type DeleteProjectRequest struct {
	Name string `json:"name" validate:"required"`
}

// DeleteProject godoc
// @Summary Удалить проект
// @Description Удаляет проект по его названию (только для создателя проекта)
// @Tags Projects
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param input body DeleteProjectRequest true "Имя проекта"
// @Success 200 {object} response.SuccessResponse "Проект успешно удален"
// @Failure 400 {object} response.ErrorResponse "Неверный формат запроса"
// @Failure 401 {object} response.ErrorResponse "Не авторизован"
// @Failure 403 {object} response.ErrorResponse "Нет прав на удаление"
// @Failure 500 {object} response.ErrorResponse "Внутренняя ошибка сервера"
// @Router /api/projects [delete]
func (s *Server) DeleteProject() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		const op = "http.DeleteProject"

		log := s.Logger.With(slog.String("op", op))

		userID, ok := r.Context().Value("userID").(int)
		if !ok {
			log.Error("failed to get userID from context")
			render.JSON(w, r, response.ErrorResponse{
				Status:  http.StatusInternalServerError,
				Message: "Internal server error",
			})
			return
		}

		var req DeleteProjectRequest

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request", sl.Err(err))
			render.JSON(w, r, response.ErrorResponse{
				Status:  http.StatusBadRequest,
				Message: "Invalid request body",
			})
			return
		}

		log.Info("deleting project",
			slog.String("name", req.Name),
			slog.Int("user_id", userID),
		)

		if err := s.Svc.DeleteProject(userID, req.Name); err != nil {
			log.Error("failed to delete project", sl.Err(err))
			render.JSON(w, r, response.ErrorResponse{
				Status:  http.StatusInternalServerError,
				Message: "Failed to delete project",
			})
			return
		}

		render.JSON(w, r, response.SuccessResponse{
			Status:  http.StatusOK,
			Message: "project deleted successfully",
		})
	}
}

type UpdateProjectRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
}

// UpdateProject godoc
// @Summary Обновить проект
// @Description Обновляет данные проекта (название и/или описание)
// @Tags Projects
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param name query string true "Название проекта для обновления"
// @Param input body UpdateProjectRequest true "Новые данные проекта"
// @Success 200 {object} response.SuccessResponse "Проект успешно обновлен"
// @Failure 400 {object} response.ErrorResponse "Неверный формат запроса"
// @Failure 401 {object} response.ErrorResponse "Не авторизован"
// @Failure 403 {object} response.ErrorResponse "Нет прав на обновление"
// @Failure 500 {object} response.ErrorResponse "Внутренняя ошибка сервера"
// @Router /api/projects [put]
func (s *Server) UpdateProject() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		const op = "http.UpdateProject"

		log := s.Logger.With(slog.String("op", op))

		userID, ok := r.Context().Value("userID").(int)
		if !ok {
			log.Error("failed to get userID from context")
			render.JSON(w, r, response.ErrorResponse{
				Status:  http.StatusInternalServerError,
				Message: "Internal server error",
			})
			return
		}

		name := r.URL.Query().Get("name")
		if name == "" {
			log.Error("empty project name in URL")
			render.JSON(w, r, response.ErrorResponse{
				Status:  http.StatusBadRequest,
				Message: "Project name is required in URL",
			})
			return
		}

		var req UpdateProjectRequest

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request", sl.Err(err))
			render.JSON(w, r, response.ErrorResponse{
				Status:  http.StatusBadRequest,
				Message: "Invalid request body",
			})
			return
		}

		log.Info("updating project",
			slog.String("original_name", name),
			slog.Int("user_id", userID),
			slog.Any("new_data", req),
		)

		var updateErrors []error

		project := model.Project{IDCreator: int64(userID), Name: name}

		if req.Name != nil {
			if err := s.Svc.UpdateProjectName(*req.Name, project); err != nil {
				log.Error("failed to update name", sl.Err(err))
				updateErrors = append(updateErrors, errors.New("failed to update name"))
			}
			project.Name = *req.Name
		}

		if req.Description != nil {
			project.Description = *req.Description
			if err := s.Svc.UpdateProjectDescription(project); err != nil {
				log.Error("failed to update description", sl.Err(err))
				updateErrors = append(updateErrors, errors.New("failed to update description"))
			}
		}

		if len(updateErrors) > 0 {
			render.JSON(w, r, response.ErrorResponse{
				Status:  http.StatusInternalServerError,
				Message: "partial update failed",
			})
			return
		}

		render.JSON(w, r, response.SuccessResponse{
			Status:  http.StatusOK,
			Message: "Project update succssfully",
		})
	}
}

// ListProjects godoc
// @Summary Список проектов
// @Description Возвращает список всех проектов пользователя
// @Tags Projects
// @Security BearerAuth
// @Produce json
// @Success 200 {object} response.SuccessResponse{data=[]model.Project} "Список проектов"
// @Failure 401 {object} response.ErrorResponse "Не авторизован"
// @Failure 500 {object} response.ErrorResponse "Внутренняя ошибка сервера"
// @Router /api/projects/list [get]
func (s *Server) ListProjects() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		const op = "http.ListProjects"

		log := s.Logger.With(slog.String("op", op))

		listProjects, err := s.Svc.ListProjects()

		if err != nil {
			log.Error("failed to read list of projects", sl.Err(err))
			render.JSON(w, r, response.ErrorResponse{
				Status:  http.StatusInternalServerError,
				Message: "failed to read list of projects",
			})
			return
		}

		render.JSON(w, r, response.SuccessResponse{
			Status:  http.StatusOK,
			Message: "projects:",
			Data:    listProjects,
		})
	}
}
