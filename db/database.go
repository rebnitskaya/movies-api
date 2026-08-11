package db

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func OpenDB(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("Failed to open db: %w", err)
	}

	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("Failed to check db with ping: %w", err)
	}

	_, err = db.Exec(`PRAGMA foreign_keys = ON`)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("Failed to enable foreign key: %w", err)
	}

	return db, nil
}
