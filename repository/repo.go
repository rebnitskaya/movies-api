package repository

import (
	"database/sql"
	m "movies_api/models"
)

type actorRepository struct {
	db *sql.DB
}

type genreRepository struct {
	db *sql.DB
}

type movieRepository struct {
	db *sql.DB
}

func NewActorRepository(db *sql.DB) ActorRepository {
	return actorRepository{db: db}
}

func NewGenreRepository(db *sql.DB) GenreRepository {
	return genreRepository{db: db}
}

func NewMovieRepository(db *sql.DB) MovieRepository {
	return movieRepository{db: db}
}

type ActorRepository interface {
	FindAllActors() ([]m.Actor, error)
	CreateActor(m.Actor) (bool, error)
	FindActorByNameAndBirthDate(string, string) (m.Actor, error)
	DeleteActorByID(int) (bool, error)
	FindActorByID(int) (m.Actor, bool)
	ReplaceFieldsInActor(int, map[string]string) (m.Actor, bool)
	FindActorsByName(string) ([]m.Actor, error)
}

type GenreRepository interface {
	FindAllGenres() ([]m.Genre, error)
	CreateGenre(m.Genre) (bool, error)
	FindGenreByID(int) (m.Genre, bool)
	ReplaceFieldsInGenre(int, string) (m.Genre, bool)
	DeleteGenreByID(int) (bool, error)
}

type MovieRepository interface {
	FindAllMovies() ([]m.Movie, error)
	CreateMovie(m.MovieDto) (m.Movie, error)
	FindMovieByID(int) (m.Movie, bool)
	ReplaceFieldsInMovie(int, map[string]string) (m.Movie, bool)
	DeleteMovieByID(int) (bool, error)
	FindMoviesByGenre(int) ([]m.Movie, error)
	FindMoviesByYear(int) ([]m.Movie, error)
	FindMoviesWithActor(int) ([]m.Movie, error)
	FindAllActorsInMovie(int) ([]m.Actor, error)
	FindMovieByTitleAndYear(string, int) (m.Movie, error)
}
