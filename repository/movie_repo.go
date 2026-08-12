package repository

import "database/sql"

type MovieRepository struct {
	db *sql.DB
}

type MovieRepositoryInterface interface {
	GetAllMovies()
}
