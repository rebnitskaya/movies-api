package repository

import (
	"database/sql"
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
	FindAllActors() ([]Actor, error)
	CreateActor(Actor) (bool, error)
	FindActorByNameAndBirthDate(string, string) (Actor, error)
	DeleteActorByID(int) (bool, error)
	FindActorByID(int) (Actor, bool)
	ReplaceFieldsInActor(int, map[string]string) (Actor, bool)
	FindActorsByName(string) ([]Actor, error)
}

type GenreRepository interface {
	FindAllGenres() ([]Genre, error)
	CreateGenre(Genre) (bool, error)
	FindGenreByID(int) (Genre, bool)
	ReplaceFieldsInGenre(int, string) (Genre, bool)
	DeleteGenreByID(int) (bool, error)
}

type MovieRepository interface {
	FindAllMovies() ([]Movie, error)
	CreateMovie(Movie) (bool, error)
	FindMovieByID(int) (Movie, bool)
	ReplaceFieldsInMovie(int, map[string]string) (Movie, bool)
	DeleteMovieByID(int) (bool, error)
	FindMoviesByGenre(int) ([]Movie, error)
	FindMoviesByYear(int) ([]Movie, error)
	FindMoviesWithActor(int) ([]Movie, error)
	FindAllActorsInMovie(int) ([]Actor, error)
}
