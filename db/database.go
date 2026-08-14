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

func CreateTables(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS actors (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			birth_date TEXT NOT NULL
		);

		CREATE TABLE IF NOT EXISTS genres (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL UNIQUE
		);

		CREATE TABLE IF NOT EXISTS movies (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			release_year INTEGER NOT NULL,
			duration INTEGER NOT NULL
		);

		CREATE TABLE IF NOT EXISTS movie_actors (
			movie_id INTEGER NOT NULL,
			actor_id INTEGER NOT NULL,

			PRIMARY KEY (movie_id,actor_id),
			FOREIGN KEY (movie_id) REFERENCES movies(id)
			FOREIGN KEY (actor_id) REFERENCES actors(id)
		);

		CREATE TABLE IF NOT EXISTS genres_movies (
			genre_id INTEGER NOT NULL,
			movie_id INTEGER NOT NULL,

			PRIMARY KEY (genre_id,movie_id)
			FOREIGN KEY (genre_id) REFERENCES genres(id)
			FOREIGN KEY (movie_id) REFERENCES movies(id)
		);
	`)

	return err
}
