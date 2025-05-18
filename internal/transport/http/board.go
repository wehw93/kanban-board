package http

import (
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
				slog.Int("user_id", userID),
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
			Data:   map[string]int64{"id": project.ID},
		})
	}
}

type ReadProjectRequest struct {
	Name string `json:"name" validate:"requires"`
}

func (s *Server) ReadProject() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http.ReadProject"
		log := s.Logger.With(slog.String("op", op))
		var req ReadProjectRequest
		err := render.DecodeJSON(r.Body, &req)

		log.Info("reading data of project", slog.String("name",req.Name))
		resp,err := s.Svc.ReadProject(req.Name)
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
