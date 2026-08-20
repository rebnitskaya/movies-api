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
	FindActorByID(int) (m.Actor, error)
	ReplaceFieldsInActor(int, map[string]string) (m.Actor, error)
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
	FindAllMovies() ([]m.MovieDto, error)
	CreateMovie(m.MovieDto) (m.Movie, error)
	FindMovieByID(int) (m.MovieDto, error)
	ReplaceFieldsInMovie(int, map[string]any) (m.Movie, error)
	DeleteMovieByID(int) (bool, error)
	FindMoviesByGenre(int) ([]m.Movie, error)
	FindMoviesByYear(int) ([]m.Movie, error)
	FindMoviesWithActor(int) ([]m.MovieDto, error)
	FindAllActorsInMovie(int) ([]m.ActorInFilmDto, error)
	FindMovieByTitleAndYear(string, int) (m.Movie, error)
	AddActorToMovie(int, int) error
	RemoveActorFromMovie(int, int) error
}
