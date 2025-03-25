.PHONY: build
build:
	go build -v ./cmd/board

.PHONY: run_migrations
run_migrations:
	go run ./cmd/migrator --migrations-path=./migrations

.DEFAULT_GOAL := build