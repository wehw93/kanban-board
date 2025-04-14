package http

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"

	"github.com/wehw93/kanban-board/internal/config"
	"github.com/wehw93/kanban-board/internal/service"
)

type Server struct {
	Server *http.Server
	Router *chi.Mux
	Logger *slog.Logger
	Svc    service.BoardService
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
	s.Router.Route("/create_user", func (r chi.Router){
		r.Use(middleware.AllowContentType("application/json"))
		r.Use(middleware.SetHeader("Content-Type", "application/json"))
		r.Post("/",s.CreateUser())
	})
}

func (s *Server) Start() error {
	return s.Server.ListenAndServe()
}
