package db

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type Executor interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
}

func Connect(info string) (*sql.DB, error) {
	db, err := sql.Open("postgres", info)
	if err != nil {
		return nil, err
	}
	return db, nil
}
