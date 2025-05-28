package http

import (
	"database/sql"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/render"
	"github.com/wehw93/kanban-board/internal/lib/http/response"
	"github.com/wehw93/kanban-board/internal/lib/logger/sl"
	"github.com/wehw93/kanban-board/internal/model"
)

type CreateTaskRequest struct {
	IDColumn    int    `json:"id_column" validate:"required"`
	Name        string `json:"name" validate:"required"`
	Description string `json:"description" validate:"required"`
}

// CreateTask godoc
// @Summary Создание новой задачи
// @Description Создает новую задачу в колонке
// @Tags Tasks
// @Accept json
// @Produce json
// @Param input body CreateTaskRequest true "Данные задачи"
// @Success 200 {object} response.SuccessResponse{data=model.Task} "Задача успешно создана"
// @Failure 400 {object} response.ErrorResponse "Неверный формат запроса"
// @Failure 422 {object} response.ErrorResponse "Ошибка при создании задачи"
// @Security BearerAuth
// @Router /api/tasks [post]
func (s *Server) CreateTask() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http.CreateTask"
		log := s.Logger.With(slog.String("op", op))

		creator_id, ok := r.Context().Value("userID").(int)
		if !ok {
			log.Error("failed to get userID from context")
			render.JSON(w, r, response.ErrorResponse{
				Status:  http.StatusInternalServerError,
				Message: "Internal server error",
			})
			return
		}

		var req CreateTaskRequest
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			log.Error("failed to decode request body", sl.Err(err))
			render.JSON(w, r, response.ErrorResponse{
				Status:  http.StatusBadRequest,
				Message: "Invalid request body",
			})
			return
		}
		log.Info("create task request",
			slog.Int("creator_id", creator_id),
			slog.String("name of task", req.Name))
		task := &model.Task{
			ID_column:   int64(req.IDColumn),
			Name:        req.Name,
			Description: req.Description,
			ID_creator:  int64(creator_id),
		}
		task.Date_of_create = time.Now().Format("2006-01-02")
		err := s.Svc.CreateTask(task)
		if err != nil {
			log.Error("failed to create task", sl.Err(err))
			render.JSON(w, r, response.ErrorResponse{
				Status:  http.StatusUnprocessableEntity,
				Message: "failed to create task",
			})
			return
		}
		render.JSON(w, r, response.SuccessResponse{
			Status: http.StatusOK,
			Data:   task,
		})
	}

}

type ReadTaskRequest struct {
	ID int `json:"id" validate:"required"`
}

// ReadTask godoc
// @Summary Получение задачи
// @Description Возвращает информацию о задаче по ID
// @Tags Tasks
// @Accept json
// @Produce json
// @Param input body ReadTaskRequest true "ID задачи"
// @Success 200 {object} response.SuccessResponse{data=model.Task} "Информация о задаче"
// @Failure 400 {object} response.ErrorResponse "Неверный формат запроса"
// @Failure 500 {object} response.ErrorResponse "Задача не найдена"
// @Security BearerAuth
// @Router /api/tasks [get]
func (s *Server) ReadTask() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http.ReadTask"
		log := s.Logger.With("op", op)
		var req ReadTaskRequest
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request", sl.Err(err))
			render.JSON(w, r, response.ErrorResponse{
				Status:  http.StatusBadRequest,
				Message: "invalid request body",
			})
			return
		}
		task := &model.Task{
			ID: int64(req.ID),
		}
		log.Info("reading data of task", slog.Int("id", req.ID))
		err = s.Svc.ReadTask(task)
		if err != nil {
			log.Error("failed to read task", sl.Err(err))
			render.JSON(w, r, response.ErrorResponse{
				Status:  http.StatusInternalServerError,
				Message: "task not found",
			})
			return
		}
		render.JSON(w, r, response.SuccessResponse{
			Status: http.StatusOK,
			Data:   task,
		})
	}
}

type DeleteTaskRequest struct {
	ID int `json:"id" validate:"required"`
}

// DeleteTask godoc
// @Summary Удаление задачи
// @Description Удаляет задачу по ID
// @Tags Tasks
// @Accept json
// @Produce json
// @Param input body DeleteTaskRequest true "ID задачи"
// @Success 200 {object} response.SuccessResponse "Задача успешно удалена"
// @Failure 400 {object} response.ErrorResponse "Неверный формат запроса"
// @Failure 500 {object} response.ErrorResponse "Ошибка при удалении задачи"
// @Security BearerAuth
// @Router /api/tasks [delete]
func (s *Server) DeleteTask() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http.DeleteTask"
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
		var req DeleteTaskRequest
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request", sl.Err(err))
			render.JSON(w, r, response.ErrorResponse{
				Status:  http.StatusBadRequest,
				Message: "Invalid request body",
			})
			return
		}
		log.Info("deleting task",
			slog.Int("id", req.ID),
			slog.Int("user_id", userID),
		)

		if err := s.Svc.DeleteTask(userID, req.ID); err != nil {
			log.Error("failed to delete task", sl.Err(err))
			render.JSON(w, r, response.ErrorResponse{
				Status:  http.StatusInternalServerError,
				Message: "Failed to delete task",
			})
			return
		}
		render.JSON(w, r, response.SuccessResponse{
			Status:  http.StatusOK,
			Message: "task deleted successfully",
		})
	}
}

type UpdateTaskRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	Id_column   *int    `json:"id_column"`
}

// UpdateTask godoc
// @Summary Обновление задачи
// @Description Обновляет имя, описание или колонку задачи
// @Tags Tasks
// @Accept json
// @Produce json
// @Param id query int true "ID задачи"
// @Param input body UpdateTaskRequest true "Обновленные данные задачи"
// @Success 200 {object} response.SuccessResponse "Задача успешно обновлена"
// @Failure 400 {object} response.ErrorResponse "Неверный формат запроса"
// @Failure 500 {object} response.ErrorResponse "Ошибка при обновлении задачи"
// @Security BearerAuth
// @Router /api/tasks [put]
func (s *Server) UpdateTask() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http.UpdateTask"
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

		id, err := strconv.Atoi(r.URL.Query().Get("id"))
		if err != nil {
			log.Error("failed to conv id")
			render.JSON(w, r, response.ErrorResponse{
				Status:  http.StatusBadRequest,
				Message: "bad request",
			})
			return
		}
		if id == 0 {
			log.Error("empty task id in URL")
			render.JSON(w, r, response.ErrorResponse{
				Status:  http.StatusBadRequest,
				Message: "task id is required in URL",
			})
			return
		}
		var req UpdateTaskRequest
		err = render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request", sl.Err(err))
			render.JSON(w, r, response.ErrorResponse{
				Status:  http.StatusBadRequest,
				Message: "Invalid request body",
			})
			return
		}
		log.Info("updating task",
			slog.Int("id", id),
			slog.Int("user_id", userID),
			slog.Any("new_data", req),
		)
		var updateErrors []error

		task := &model.Task{
			ID:          int64(id),
			ID_executor: sql.NullInt64{Int64: int64(userID), Valid: userID != 0},
		}
		if req.Name != nil {
			task.ID_column = int64(*req.Id_column)
			if err := s.Svc.UpdateTaskName(task); err != nil {
				log.Error("failed to update name", sl.Err(err))
				updateErrors = append(updateErrors, errors.New("failed to update name"))
			}
		}
		if req.Description != nil {
			task.Description = *req.Description
			if err := s.Svc.UpdateTaskDescription(task); err != nil {
				log.Error("failed to update description", sl.Err(err))
				updateErrors = append(updateErrors, errors.New("failed to update description"))
			}
		}
		if req.Id_column != nil {
			task.ID_column = int64(*req.Id_column)
			if err := s.Svc.UpdateTaskColumn(task); err != nil {
				log.Error("failed to update column id", sl.Err(err))
				updateErrors = append(updateErrors, errors.New("failed to update column id"))
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
			Message: "task update succssfully",
		})
	}
}

// GetLogsTask godoc
// @Summary Получение логов задачи
// @Description Возвращает логи действий по задаче
// @Tags Tasks
// @Accept json
// @Produce json
// @Param id query int true "ID задачи"
// @Success 200 {object} response.SuccessResponse{data=[]model.Task_log} "Логи задачи"
// @Failure 400 {object} response.ErrorResponse "Неверный ID"
// @Failure 500 {object} response.ErrorResponse "Ошибка при получении логов"
// @Security BearerAuth
// @Router /api/tasks/logs [get]
func (s *Server) GetLogsTask() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http.GetLogsTask"

		log := s.Logger.With(slog.String("op", op))

		id_task, err := strconv.Atoi(r.URL.Query().Get("id"))
		if err != nil {
			log.Error("failed to conv id", sl.Err(err))
			render.JSON(w, r, response.ErrorResponse{
				Status:  http.StatusBadRequest,
				Message: "bad request",
			})
			return
		}
		logs, err := s.Svc.GetLogsTask(id_task)
		if err != nil {
			log.Error("failed to get logs task", sl.Err(err))
			render.JSON(w, r, response.ErrorResponse{
				Status:  http.StatusInternalServerError,
				Message: "failed to get logs task",
			})
			return
		}
		render.JSON(w, r, response.SuccessResponse{
			Status:  http.StatusOK,
			Message: "logs of task id:" + strconv.Itoa(id_task),
			Data:    logs,
		})
	}
}
