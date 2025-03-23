.PHONY: build
build:
	go build -v ./cmd/board

.DEFAULT_GOAL := build