package http

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/render"
	"github.com/wehw93/kanban-board/internal/lib/http/response"
	"github.com/wehw93/kanban-board/internal/lib/logger/sl"
	"github.com/wehw93/kanban-board/internal/model"
)

type CreateColumnRequest struct {
	Name      string `json:"name" validate:"required"`
	ProjectID int    `json:"id_project" validate:"required"`
}

// CreateColumn godoc
// @Summary Создание новой колонки
// @Description Создает новую колонку в указанном проекте
// @Tags Columns
// @Accept json
// @Produce json
// @Param input body CreateColumnRequest true "Данные колонки"
// @Success 201 {object} response.SuccessResponse{data=model.Column} "Колонка успешно создана"
// @Failure 400 {object} response.ErrorResponse "Неверный формат запроса"
// @Failure 422 {object} response.ErrorResponse "Ошибка при создании колонки"
// @Security BearerAuth
// @Router /api/columns [post]
func (s *Server) CreateColumn() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		const op = "transport.http.CreateColumn"

		log := s.logger.With(slog.String("op", op))

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
			slog.String("column_name", req.Name),
		)

		column := &model.Column{
			Name:       req.Name,
			ID_project: int64(req.ProjectID),
		}

		if err := s.boardSvc.CreateColumn(column); err != nil {
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
			Data:   column,
		})
	}
}

type ReadColumnRequest struct {
	Name      string `json:"name" validate:"required"`
	IDProject int    `json:"id_project" validate:"required"`
}

// ReadColumn godoc
// @Summary Получение информации о колонке
// @Description Возвращает информацию о колонке по имени и ID проекта
// @Tags Columns
// @Accept json
// @Produce json
// @Param input body ReadColumnRequest true "Параметры запроса"
// @Success 200 {object} response.SuccessResponse{data=model.Column} "Информация о колонке"
// @Failure 400 {object} response.ErrorResponse "Неверный формат запроса"
// @Failure 404 {object} response.ErrorResponse "Колонка не найдена"
// @Security BearerAuth
// @Router /api/columns [get]
func (s *Server) ReadColumn() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		const op = "http.ReadColumn"

		log := s.logger.With("op", op)

		var req ReadColumnRequest

		err := render.DecodeJSON(r.Body, &req)

		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))
			render.JSON(w, r, response.ErrorResponse{
				Status:  http.StatusBadRequest,
				Message: "invalid request body",
			})
		}

		column := model.Column{
			Name:       req.Name,
			ID_project: int64(req.IDProject),
		}

		resp, err := s.boardSvc.ReadColumn(column)
		if err != nil {
			log.Error("failed to read column", sl.Err(err))
			render.JSON(w, r, response.ErrorResponse{
				Status:  http.StatusBadRequest,
				Message: "failed to read column",
			})
		}

		render.JSON(w, r, response.SuccessResponse{
			Status: http.StatusOK,
			Data:   resp,
		})
	}
}

type DeleteColumnRequest struct {
	ID int `json:"id" validate:"required"`
}

// DeleteColumn godoc
// @Summary Удаление колонки
// @Description Удаляет колонку по её ID
// @Tags Columns
// @Accept json
// @Produce json
// @Param input body DeleteColumnRequest true "ID колонки"
// @Success 200 {object} response.SuccessResponse "Колонка успешно удалена"
// @Failure 400 {object} response.ErrorResponse "Неверный формат запроса"
// @Failure 500 {object} response.ErrorResponse "Ошибка сервера при удалении колонки"
// @Security BearerAuth
// @Router /api/columns [delete]
func (s *Server) DeleteColumn() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		const op = "http.DeleteColumn"

		log := s.logger.With(slog.String("op", op))

		var req DeleteColumnRequest

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request", sl.Err(err))
			render.JSON(w, r, response.ErrorResponse{
				Status:  http.StatusBadRequest,
				Message: "Invalid request body",
			})
			return
		}

		log.Info("deleting column",
			slog.Int("column_id", req.ID),
		)

		if err := s.boardSvc.DeleteColumn(req.ID); err != nil {
			log.Error("failed to delete column", sl.Err(err))
			render.JSON(w, r, response.ErrorResponse{
				Status:  http.StatusInternalServerError,
				Message: "Failed to delete column",
			})
			return
		}

		render.JSON(w, r, response.SuccessResponse{
			Status:  http.StatusOK,
			Message: "Column deleted successfully",
		})
	}
}

type UpdateColumnRequest struct {
	Name *string `json:"name"`
}

// UpdateColumn godoc
// @Summary Обновление информации о колонке
// @Description Обновляет данные колонки по её ID
// @Tags Columns
// @Accept json
// @Produce json
// @Param id query int true "ID колонки"
// @Param input body UpdateColumnRequest true "Новые данные колонки"
// @Success 200 {object} response.SuccessResponse "Колонка успешно обновлена"
// @Failure 400 {object} response.ErrorResponse "Неверный формат запроса"
// @Failure 500 {object} response.ErrorResponse "Ошибка сервера при обновлении колонки"
// @Security BearerAuth
// @Router /api/columns [put]
func (s *Server) UpdateColumn() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		const op = "http.UpdateColumn"

		log := s.logger.With(slog.String("op", op))

		id, err := strconv.Atoi(r.URL.Query().Get("id"))
		if err != nil {
			log.Error("failed to get id from url", err)
			render.JSON(w, r, response.ErrorResponse{
				Status:  http.StatusInternalServerError,
				Message: "failed to get id rom url",
			})
			return
		}

		if id == 0 {
			log.Error("empty column id in URL")
			render.JSON(w, r, response.ErrorResponse{
				Status:  http.StatusBadRequest,
				Message: "column id is required in URL",
			})
			return
		}

		var req UpdateColumnRequest

		err = render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request", sl.Err(err))
			render.JSON(w, r, response.ErrorResponse{
				Status:  http.StatusBadRequest,
				Message: "Invalid request body",
			})
			return
		}

		log.Info("updating column",
			slog.Int("id", id),
			slog.Any("new_data", req),
		)

		var updateErrors []error

		column := model.Column{ID: int64(id)}

		if req.Name != nil {
			if err := s.boardSvc.UpdateColumnName(column, *req.Name); err != nil {
				log.Error("failed to update name", sl.Err(err))
				updateErrors = append(updateErrors, errors.New("failed to update name"))
			}
			column.Name = *req.Name
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
			Message: "column update succssfully",
		})
	}
}
