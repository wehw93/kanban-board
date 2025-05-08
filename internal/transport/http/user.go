package http

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/go-chi/render"
	"github.com/wehw93/kanban-board/internal/lib/http/response"
	"github.com/wehw93/kanban-board/internal/lib/jwt/helpers_jwt"
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

type ReadUserRequest struct {
	JWTToken string `json:"jwt-token" validate:"required"`
}

func (s *Server) ReadUser() http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		const op = "transport.http.ReadUser"
		log := s.Logger.With(slog.String("op:", op))
		var r ReadUserRequest
		tokenString := req.Header.Get("Authorization")
		if tokenString == "" {
			log.Error("authorization header missing")
			s.error(resp, req, http.StatusUnauthorized, errors.New("authorization header required"))
			return
		}
		log.Info("succes decode request", slog.Any("request", r))
		s.JWTSecret = JWTSecret
		parts := strings.Split(tokenString, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			log.Error("invalid authorization header format")
			s.error(resp, req, http.StatusUnauthorized, errors.New("invalid authorization header format"))
			return
		}
		claims, err := helpers_jwt.ParseToken(parts[1], s.JWTSecret)
		if err != nil {
			log.Error("failed to parse token", sl.Err(err))
			s.error(resp, req, http.StatusUnauthorized, errors.New("invalid token"))
			return
		}
		user_id, ok := claims["uid"].(float64)
		if !ok {
			log.Error("invalid user ID in token")
			s.error(resp, req, http.StatusUnauthorized, errors.New("invalid user ID in token"))
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

/*
type DeleteUserRequest struct{

}

	func (s*Server) DeleteUser() http.HandlerFunc{
		return func (resp http.ResponseWriter, req*http.Request){
			const op = "transport.http.DeleteUser"
			log :=s.Logger.With(slog.String("op:",op))
			var r DeleteUserRequest
			err:=render.DecodeJSON(req.Body,&r)
		}
	}
*/
func (s *Server) error(w http.ResponseWriter, r *http.Request, code int, err error) {
	s.respond(w, r, code, map[string]string{"error": err.Error()})
}

func (s *Server) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}
