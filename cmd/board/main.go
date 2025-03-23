package main

import (
	"log/slog"
	"os"

	"github.com/wehw93/kanban-board/internal/config"
	"github.com/wehw93/kanban-board/internal/storage/postgresql"
)

const (
	env_local = "local"
	env_prod  = "prod"
)

func main() {
	cfg := config.MustLoad()
	log := SetupLogger(cfg.Env)
	log.Info("starting server")
	store, err := postgresql.New(cfg.DB.GetDSN())
	if err!=nil{
		panic(err)
	}
	log.Info("postgres port: ",cfg.DB.Port)
	defer store.Close()
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
