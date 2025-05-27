package http

import (
	"context"
	"log/slog"
	"net/http"
	"strings"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	httpSwagger "github.com/swaggo/http-swagger"

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
		r.Post("/register", s.CreateUser())
		r.Post("/login", s.LoginUser())
	})
	s.Router.Route("/api", func(r chi.Router) {
		r.Use(middleware.AllowContentType("application/json"))
		r.Use(middleware.SetHeader("Content-Type", "application/json"))
		r.Use(s.AuthentificationUser)

		r.Route("/users", func(r chi.Router) {
			r.Get("/me", s.ReadUser())
			r.Put("/me", s.UpdateUser())
			r.Delete("/me", s.DeleteUser())
		})

		r.Route("/projects", func(r chi.Router) {
			r.Post("/", s.CreateProject())
			r.Get("/read", s.ReadProject())
			r.Delete("/",s.DeleteProject())
			r.Put("/",s.UpdateProject())
			r.Get("/list",s.ListProjects())
		})
		r.Route("/columns",func (r chi.Router){
			r.Post("/",s.CreateColumn())
			r.Get("/",s.ReadColumn())
			r.Delete("/",s.DeleteColumn())
			r.Put("/",s.UpdateColumn())
		})
		r.Route("/tasks",func (r chi.Router){
			r.Post("/",s.CreateTask())
			r.Get("/",s.ReadTask())
			r.Delete("/",s.DeleteTask())
			r.Put("/",s.UpdateTask())
			s.Router.Get("/swagger/*", httpSwagger.Handler(
				httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
			))
		})
		s.Router.Get("/swagger/*", httpSwagger.WrapHandler)
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
