package http

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/render"
	"github.com/wehw93/kanban-board/internal/lib/http/response"
	"github.com/wehw93/kanban-board/internal/lib/logger/sl"
	"github.com/wehw93/kanban-board/internal/model"
)

type CreateColumnRequest struct {
	Name      string `json:"name" validate:"required"`
	ProjectID int    `json:"id_project" validate:"required"`
}

func (s *Server) CreateColumn() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "transport.http.CreateColumn"
		log := s.Logger.With(slog.String("op", op))
		var req CreateColumnRequest
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			log.Error("failed to decode request body", sl.Err(err))
			render.JSON(w, r, response.ErrorResponse{
				Status:  http.StatusBadRequest,
				Message: "Invalid request body",
			})
			return
		}
		log.Info("create column request",
			slog.Int("project_id", req.ProjectID),
			slog.String("project_name", req.Name),
		)
		column := &model.Column{
			Name:       req.Name,
			ID_project: int64(req.ProjectID),
		}
		if err := s.Svc.CreateColumn(column); err != nil {
			log.Error("failed to create column",
				sl.Err(err),
			)
			render.JSON(w, r, response.ErrorResponse{
				Status:  http.StatusUnprocessableEntity,
				Message: "Failed to create column",
			})
			return
		}
		log.Info("column created successfully",
			slog.Int("column_id", int(column.ID)),
			slog.Int64("project_id", column.ID_project),
		)

		render.JSON(w, r, response.SuccessResponse{
			Status: http.StatusCreated,
			Data:   map[string]int64{"id": column.ID},
		})
	}
}

type ReadColumnRequest struct{
	Name string `json:"name" validate:"required"`
	IDProject int `json:"id_project" validate:"required"`
}


func (s*Server)ReadColumn()http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http.ReadColumn"

		log:=s.Logger.With("op",op)
		var req ReadColumnRequest
		err:=render.DecodeJSON(r.Body,&req)
		if err!=nil{
			log.Error("failed to decode request body",sl.Err(err))
			render.JSON(w,r,response.ErrorResponse{
				Status: http.StatusBadRequest,
				Message: "invalid request body",
			})
		}
		column:=model.Column{
			Name: req.Name,
			ID_project: int64(req.IDProject),
		}
		resp,err:=s.Svc.ReadColumn(column)
		if err!=nil{
			log.Error("failed to read column",sl.Err(err))
			render.JSON(w,r,response.ErrorResponse{
				Status: http.StatusBadRequest,
				Message: "failed to read column",
			})
		}
		render.JSON(w,r,response.SuccessResponse{
			Status: http.StatusOK,
			Data: resp,
		})
	}
}
