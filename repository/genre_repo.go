package repository

import "database/sql"

type GenreRepository struct {
	db *sql.DB
}

type GenreRepositoryInterface interface {
	GetAllGenres()
}
