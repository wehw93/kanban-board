package main

// @title Kanban Board API
// @version 1.0
// @description API для управления проектами и задачами

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @host localhost:8080
import (
	"log/slog"
	"os"

	_ "github.com/wehw93/kanban-board/docs"
	"github.com/wehw93/kanban-board/internal/config"
	"github.com/wehw93/kanban-board/internal/lib/logger/sl"
	"github.com/wehw93/kanban-board/internal/service/auth"
	"github.com/wehw93/kanban-board/internal/service/board"
	"github.com/wehw93/kanban-board/internal/storage/postgresql"
	server "github.com/wehw93/kanban-board/internal/transport/http"
)

const (
	jwtSecret = "secret"
	env_local = "local"
	env_prod  = "prod"
)

func main() {
	cfg := config.MustLoad()

	log := SetupLogger(cfg.Env)
	log.Info("starting server")

	store, err := postgresql.New(cfg.DB.GetDSN())
	if err != nil {
		panic(err)
	}
	log.Info("postgres port: ", cfg.DB.Port)
	defer store.Close()

	svcAuth := auth.NewService(jwtSecret)

	svcBoard := board.NewService(store, jwtSecret)

	srv := server.NewServer(cfg, log, svcBoard, svcAuth)

	srv.InitRoutes()

	log.Info("starting server", slog.String("addr: ", cfg.HTTP_Server.Address))
	if err := srv.Start(); err != nil {
		log.Error("failed to start server", sl.Err(err))
		os.Exit(1)
	}
}

func SetupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case env_local:
		log = slog.New(slog.NewTextHandler(os.Stdin, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case env_prod:
		log = slog.New(slog.NewTextHandler(os.Stdin, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}
