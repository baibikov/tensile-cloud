package multistore

import (
	"io"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Store interface {
	DB() *sqlx.DB

	io.Closer
}

type store struct {
	db *sqlx.DB
}

func (s store) Close() error {
	return s.db.Close()
}

func (s store) DB() *sqlx.DB {
	return s.db
}

type Config struct {
	Conn string
}

func New(config *Config) (Store, error) {
	db, err := sqlx.Open("postgres", config.Conn)
	if err != nil {
		return nil, err
	}

	return &store{db: db}, nil
}
