package http

import (
	"log/slog"
	"net/http"
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
			ID_column: int64(req.IDColumn),
			Name: req.Name,
			Description: req.Description,
			ID_creator: int64(creator_id),
		}
		task.Status = "todo"
		task.Date_of_create = time.Now().Format("2006-01-02")
		err:=s.Svc.CreateTask(task)
		if err!=nil{
			log.Error("failed to create task", sl.Err(err))
			render.JSON(w, r, response.ErrorResponse{
				Status:  http.StatusUnprocessableEntity,
				Message: "failed to create project",
			})
			return
		}
		render.JSON(w,r,response.SuccessResponse{
			Status: http.StatusOK,
			Data: task,
		})
	}

}
