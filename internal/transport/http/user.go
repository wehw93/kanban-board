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

const JWTSecret = "secret"

func (s *Server) CreateUser() http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		const op = "transport.http.CreateUser"
		log := s.Logger.With(slog.String("op:", op))
		var r CreateUserRequest
		err := render.DecodeJSON(req.Body, &r)
		if err != nil {
			log.Error("Failed to prepare user", sl.Err(err))
			render.JSON(resp, req, response.Error("Failed to decode request"))
			return
		}
		log.Info("createUser", slog.Any("request", r))

		user := &model.User{
			Name:     r.Name,
			Email:    r.Email,
			Password: r.Password,
		}
		err = user.BeforeCreate()
		if err != nil {
			log.Error("failed to prepare user", sl.Err(err))
			render.JSON(resp, req, response.Error("failed to decode request"))
			return
		}
		log.Info("creating user with data", slog.String("Email", user.Email), slog.String("Password", user.Password))

		err = s.Svc.CreateUser(user)
		if err != nil {
			s.error(resp, req, http.StatusUnprocessableEntity, err)
			return
		}
		s.respond(resp, req, http.StatusCreated, user)

	}
}

type LoginUserRequest struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func (s *Server) LoginUser() http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		const op = "transport.http.Login"
		log := s.Logger.With(slog.String("op:", op))
		var r LoginUserRequest
		err := render.DecodeJSON(req.Body, &r)
		if err != nil {
			log.Error("Failed to decode request", sl.Err(err))
			render.JSON(resp, req, response.Error("Failed to decode request"))
			return
		}
		log.Info("succes decode request", slog.Any("request", r))

		token, err := s.Svc.LoginUser(r.Email, r.Password)
		if err != nil {
			s.error(resp, req, http.StatusUnprocessableEntity, err)
			return
		}
		s.JWTSecret = JWTSecret
		//res, err := helpers_jwt.ParseToken(token, s.JWTSecret)
		s.respond(resp, req, http.StatusCreated, token)
	}
}

func (s *Server) ReadUser() http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		const op = "transport.http.ReadUser"
		log := s.Logger.With(slog.String("op:", op))
		user_id, ok := req.Context().Value("userID").(int)
		if !ok {
			log.Error("failed to get userID from context")
			s.error(resp, req, http.StatusInternalServerError, errors.New("internal server error"))
			return
		}
		userData, err := s.Svc.ReadUser(int(user_id))
		if err != nil {
			log.Error("failed to read user data", slog.Int("user_id", int(user_id)), sl.Err(err))
			s.error(resp, req, http.StatusInternalServerError, err)
			return
		}
		s.respond(resp, req, http.StatusOK, userData)
	}
}

func (s *Server) DeleteUser() http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		const op = "transport.http.DeleteUser"
		log := s.Logger.With(slog.String("op:", op))
		user_id, ok := req.Context().Value("userID").(int)
		if !ok {
			log.Error("failed to get userID from context")
			s.error(resp, req, http.StatusInternalServerError, errors.New("internal server error"))
			return
		}
		err := s.Svc.DeleteUser(int(user_id))
		if err != nil {
			log.Error("failed to delete user ", slog.Int("user_id", int(user_id)), sl.Err(err))
			s.error(resp, req, http.StatusInternalServerError, err)
			return
		}
		s.respond(resp, req, http.StatusOK, "user deleted")
	}
}

type UpdateUserRequest struct {
	Email    *string `json:"email"`
	Password *string `json:"password"`
}

func (s *Server) UpdateUser() http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		const op = "transport.http.UpdateUser"
		log := s.Logger.With(slog.String("op:", op))
		user_id, ok := req.Context().Value("userID").(int)
		if !ok {
			log.Error("failed to get userID from context")
			s.error(resp, req, http.StatusInternalServerError, errors.New("internal server error"))
			return
		}
		var r UpdateUserRequest
		render.DecodeJSON(req.Body, &r)
		var user model.User
		user.ID = user_id
		var err error
		if r.Email != nil {
			user.Email = *r.Email
			err := s.Svc.UpdateEmail(user)
			if err != nil {
				log.Error("failed to update email")
				s.error(resp, req, http.StatusInternalServerError, errors.New("internal server error"))
				return
			}
		}
		if r.Password != nil {
			user.Password = *r.Password
			user.Encrypted_password, err = model.Encryptstring(user.Password)
			if err != nil {
				log.Error("failed to encrypt password")
				s.error(resp, req, http.StatusInternalServerError, errors.New("internal server error"))
				return
			}
			err = s.Svc.UpdatePassword(user)
			if err != nil {
				log.Error("failed to update password")
				s.error(resp, req, http.StatusInternalServerError, errors.New("internal server error"))
				return
			}
		}
		s.respond(resp, req, http.StatusOK, "user data updated")
	}
}
