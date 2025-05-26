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
		task.Status = "todo"
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
			task.ID_column=   int64(*req.Id_column)
			if err := s.Svc.UpdateTaskName(task); err != nil {
				log.Error("failed to update name", sl.Err(err))
				updateErrors = append(updateErrors, errors.New("failed to update name"))
			}
		}
		if req.Description!=nil{
			task.Description =*req.Description
			if err:=s.Svc.UpdateTaskDescription(task);err!=nil{
				log.Error("failed to update description", sl.Err(err))
				updateErrors = append(updateErrors, errors.New("failed to update description"))
			}
		}
		if req.Id_column!=nil{
			task.ID_column = int64(*req.Id_column)
			if err:=s.Svc.UpdateTaskColumn(task);err!=nil{
				log.Error("failed to update column id", sl.Err(err))
				updateErrors = append(updateErrors, errors.New("failed to update column id"))
			}
		}
		if len(updateErrors)>0{
			render.JSON(w,r,response.ErrorResponse{
				Status: http.StatusInternalServerError,
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
