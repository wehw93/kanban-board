package http

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/go-chi/render"
	"github.com/wehw93/kanban-board/internal/http/responce"
	"github.com/wehw93/kanban-board/internal/lib/jwt/helpers_jwt"
	"github.com/wehw93/kanban-board/internal/lig/logger/sl"
	"github.com/wehw93/kanban-board/internal/model"
)

type CreateUserRequest struct {
	Name     string `json:"name" validate:"required"`
	Password string `json:"password" validate:"required"`
	Email    string `json:"email" validate:"required"`
}

func (s *Server) CreateUser() http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		const op = "transport.http.CreateUser"
		log := s.Logger.With(slog.String("op:", op))
		var r CreateUserRequest
		err := render.DecodeJSON(req.Body, &r)
		if err != nil {
			log.Error("Failed to prepare user", sl.Err(err))
			render.JSON(resp, req, responce.Error("Failed to decode request"))
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
			render.JSON(resp, req, responce.Error("failed to decode request"))
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

func (s *Server) Login() http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		const op = "transport.http.Login"
		log := s.Logger.With(slog.String("op:", op))
		var r LoginUserRequest
		err := render.DecodeJSON(req.Body, &r)
		if err != nil {
			log.Error("Failed to prepare user", sl.Err(err))
			render.JSON(resp, req, responce.Error("Failed to decode request"))
			return
		}
		log.Info("LoginUser", slog.Any("request", r))

		token, err := s.Svc.Login(r.Email, r.Password)
		if err != nil {
			s.error(resp, req, http.StatusUnprocessableEntity, err)
			return
		}
		res, err := helpers_jwt.ParseToken(token, "secret")
		s.respond(resp, req, http.StatusCreated, res)
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
