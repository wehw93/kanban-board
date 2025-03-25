package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	_"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/wehw93/kanban-board/internal/config"
)

func main() {
	var migrationsPath string

	flag.StringVar(&migrationsPath, "migrations-path", "", "path to migrations")

	flag.Parse()

	cfg := config.MustLoad()

	connURL, err := convertConnStringToURL(cfg.DB.GetDSN())
	if err != nil {
		panic(err)
	}

	m, err := migrate.New("file://"+
		migrationsPath,
		connURL)
	if err != nil {
		log.Fatal(err)
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("no migrations to apply")
			return
		}
		log.Fatal(err)
	}
	fmt.Println("migrations applied succesfully")

}

// convertConnStringToURL преобразует формат "host=... port=..." в URL "postgres://..."
func convertConnStringToURL(connStr string) (string, error) {
	parts := strings.Fields(connStr)
	params := make(map[string]string)

	for _, part := range parts {
		keyValue := strings.SplitN(part, "=", 2)
		if len(keyValue) != 2 {
			return "", fmt.Errorf("invalid connection string format")
		}
		params[keyValue[0]] = keyValue[1]
	}

	requiredKeys := []string{"host", "port", "dbname", "user", "password", "sslmode"}
	for _, key := range requiredKeys {
		if _, ok := params[key]; !ok {
			return "", fmt.Errorf("missing required key: %s", key)
		}
	}

	urlStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		url.QueryEscape(params["user"]),
		url.QueryEscape(params["password"]),
		params["host"],
		params["port"],
		params["dbname"],
		params["sslmode"],
	)

	return urlStr, nil
}
