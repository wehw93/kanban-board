package http

import (
	"context"
	"log/slog"
	"net/http"
	"strings"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"

	"github.com/wehw93/kanban-board/internal/config"
	"github.com/wehw93/kanban-board/internal/lib/http/response"
	"github.com/wehw93/kanban-board/internal/lib/jwt/helpers_jwt"
	"github.com/wehw93/kanban-board/internal/lib/logger/sl"
	"github.com/wehw93/kanban-board/internal/service"
)

const JWTSecret = "secret"

type Server struct {
	Server    *http.Server
	Router    *chi.Mux
	Logger    *slog.Logger
	Svc       service.BoardService
	JWTSecret string
}

func NewServer(cfg *config.Config, logger *slog.Logger, svc service.BoardService) *Server {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	return &Server{
		Svc:    svc,
		Router: router,
		Logger: logger,
		Server: &http.Server{
			Addr:        cfg.HTTP_Server.Address,
			Handler:     router,
			ReadTimeout: cfg.HTTP_Server.Timeout,
			IdleTimeout: cfg.HTTP_Server.IdleTimeout,
		},
	}
}

func (s *Server) InitRoutes() {
	s.Router.Route("/auth", func(r chi.Router) {
		r.Use(middleware.AllowContentType("application/json"))
		r.Use(middleware.SetHeader("Content-Type", "application/json"))
		r.Post("/create_user", s.CreateUser())
		r.Post("/login_user", s.LoginUser())
	})
	s.Router.Route("/api", func(r chi.Router) {
		r.Use(middleware.AllowContentType("application/json"))
		r.Use(middleware.SetHeader("Content-Type", "application/json"))
		r.Use(s.AuthentificationUser)
		r.Post("/read_user", s.ReadUser())
		r.Post("/delete_user", s.DeleteUser())
		r.Post("/update_user", s.UpdateUser())
		r.Post("/create_project", s.CreateProject())
	})
}

func (s *Server) AuthentificationUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		const op = "middleware.auth"
		log := s.Logger.With(slog.String("op", op))
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			log.Error("authorization header missing")
			render.JSON(w, r, response.ErrorResponse{
				Status:  http.StatusUnauthorized,
				Message: "Authorization header missing",
			})
			return
		}
		parts := strings.Split(tokenString, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			log.Error("invalid authorization header format")
			render.JSON(w, r, response.ErrorResponse{
				Status:  http.StatusUnauthorized,
				Message: "Invalid authorization header format",
			})
			return
		}
		s.JWTSecret = JWTSecret
		claims, err := helpers_jwt.ParseToken(parts[1], s.JWTSecret)
		if err != nil {
			log.Error("failed to parse token", sl.Err(err))
			render.JSON(w, r, response.ErrorResponse{
				Status:  http.StatusUnauthorized,
				Message: "Invalid token",
			})
			return
		}
		user_id, ok := claims["uid"].(float64)
		if !ok {
			log.Error("invalid user ID in token")
			render.JSON(w, r, response.ErrorResponse{
				Status:  http.StatusUnauthorized,
				Message: "Invalid user ID in token",
			})
			return
		}
		ctx := context.WithValue(r.Context(), "userID", int(user_id))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (s *Server) Start() error {
	return s.Server.ListenAndServe()
}
