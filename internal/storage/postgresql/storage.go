package postgresql

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/wehw93/kanban-board/internal/storage"
)

type Storage struct {
	db             *sql.DB
	userrepository *UserRepository
}

func New(dsn string) (*Storage, error) {
	const op = "storage.postgresql.new"

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}
	return &Storage{db: db}, nil
}

func (s *Storage) User() storage.UserRepository {
	if s.userrepository != nil {
		return s.userrepository
	}
	s.userrepository = &UserRepository{
		store: s,
	}
	return s.userrepository
}

func (s *Storage) Close() {
	s.db.Close()
}
