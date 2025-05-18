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

type CreateUserRequest struct {
	Name     string `json:"name" validate:"required"`
	Password string `json:"password" validate:"required"`
	Email    string `json:"email" validate:"required"`
}



func (s *Server) CreateUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http.CreateUser"
		log := s.Logger.With(slog.String("op", op))
		var req CreateUserRequest
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			log.Error("failed to decode request", sl.Err(err))
			render.JSON(w, r, response.ErrorResponse{
				Status:  http.StatusBadRequest,
				Message: "Invalid request body",
			})
			return
		}
		log.Info("create user request", 
			slog.String("email", req.Email),
		)

		user := &model.User{
			Name:     req.Name,
			Email:    req.Email,
			Password: req.Password,
		}
		if err := user.BeforeCreate(); err != nil {
			log.Error("failed to prepare user", sl.Err(err))
			render.JSON(w, r, response.ErrorResponse{
				Status:  http.StatusBadRequest,
				Message: "Failed to process user data",
			})
			return
		}
		if err := s.Svc.CreateUser(user); err != nil {
			log.Error("failed to create user", sl.Err(err))
			render.JSON(w, r, response.ErrorResponse{
				Status:  http.StatusInternalServerError,
				Message: "Failed to create user",
			})
			return
		}
		log.Info("user created", slog.Int("user_id", user.ID))
		render.JSON(w, r, response.SuccessResponse{
			Status: http.StatusCreated,
			Data:   user,
		})
	}
}

type LoginUserRequest struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func (s *Server) LoginUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http.LoginUser"
		log := s.Logger.With(slog.String("op", op))
		var req LoginUserRequest
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			log.Error("failed to decode request", sl.Err(err))
			render.JSON(w, r, response.ErrorResponse{
				Status:  http.StatusBadRequest,
				Message: "Invalid request body",
			})
			return
		}
		log.Info("login attempt", slog.String("email", req.Email))

		token, err := s.Svc.LoginUser(req.Email, req.Password)
		if err != nil {
			log.Error("login failed", slog.String("email", req.Email), sl.Err(err))
			render.JSON(w, r, response.ErrorResponse{
				Status:  http.StatusInternalServerError,
				Message: "Invalid credentials",
			})
			return
		}
		log.Info("login successful", slog.String("email", req.Email))
		render.JSON(w, r, response.SuccessResponse{
			Status: http.StatusOK,
			Data:   map[string]string{"token": token},
		})
	}
}

func (s *Server) ReadUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http.ReadUser"
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
		log.Info("reading user data", slog.Int("user_id", userID))

		user, err := s.Svc.ReadUser(userID)
		if err != nil {
			log.Error("failed to read user", slog.Int("user_id", userID), sl.Err(err))
			render.JSON(w, r, response.ErrorResponse{
				Status:  http.StatusInternalServerError,
				Message: "User not found",
			})
			return
		}
		render.JSON(w, r, response.SuccessResponse{
			Status: http.StatusOK,
			Data:   user,
		})
	}
}

func (s *Server) DeleteUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http.DeleteUser"
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
		log.Info("deleting user", slog.Int("user_id", userID))

		if err := s.Svc.DeleteUser(userID); err != nil {
			log.Error("failed to delete user", slog.Int("user_id", userID), sl.Err(err))
			render.JSON(w, r, response.ErrorResponse{
				Status:  http.StatusInternalServerError,
				Message: "Failed to delete user",
			})
			return
		}

		render.JSON(w, r, response.SuccessResponse{
			Status:  http.StatusOK,
			Message: "User deleted successfully",
		})
	}
}

type UpdateUserRequest struct {
	Email    *string `json:"email"`
	Password *string `json:"password"`
}

func (s *Server) UpdateUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http.UpdateUser"
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

		var req UpdateUserRequest
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			log.Error("failed to decode request", sl.Err(err))
			render.JSON(w, r, response.ErrorResponse{
				Status:  http.StatusBadRequest,
				Message: "Invalid request body",
			})
			return
		}

		log.Info("updating user", 
			slog.Int("user_id", userID),
			slog.Bool("has_email", req.Email != nil),
			slog.Bool("has_password", req.Password != nil),
		)

		user := model.User{ID: userID}
		var updateErrors []error

		if req.Email != nil {
			user.Email = *req.Email
			if err := s.Svc.UpdateEmail(user); err != nil {
				log.Error("failed to update email", sl.Err(err))
				updateErrors = append(updateErrors, errors.New("failed to update email"))
			}
		}

		if req.Password != nil {
			user.Password = *req.Password
			encrypted, err := model.Encryptstring(user.Password)
			if err != nil {
				log.Error("failed to encrypt password", sl.Err(err))
				updateErrors = append(updateErrors, errors.New("failed to process password"))
			} else {
				user.Encrypted_password = encrypted
				if err := s.Svc.UpdatePassword(user); err != nil {
					log.Error("failed to update password", sl.Err(err))
					updateErrors = append(updateErrors, errors.New("failed to update password"))
				}
			}
		}

		if len(updateErrors) > 0 {
			render.JSON(w, r, response.ErrorResponse{
				Status:  http.StatusInternalServerError,
				Message: "Partial update failed",
			})
			return
		}

		render.JSON(w, r, response.SuccessResponse{
			Status:  http.StatusOK,
			Message: "User updated successfully",
		})
	}
}
