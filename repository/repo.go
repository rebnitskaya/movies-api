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
	FindAllActors(int, int) ([]m.Actor, error)
	CreateActor(m.Actor) (m.Actor, error)
	FindActorByNameAndBirthDate(string, string) (m.Actor, error)
	DeleteActorByID(int) (bool, error)
	FindActorByID(int) (m.Actor, error)
	ReplaceFieldsInActor(int, map[string]string) (m.Actor, error)
	FindActorsByName(string) ([]m.Actor, error)
}

type GenreRepository interface {
	FindAllGenres() ([]m.Genre, error)
	CreateGenre(m.Genre) (m.Genre, error)
	FindGenreByID(int) (m.Genre, error)
	ReplaceFieldsInGenre(int, string) (m.Genre, error)
	DeleteGenreByID(int) (bool, error)
	FindGenreByName(string) (m.Genre, error)
}

type MovieRepository interface {
	FindAllMovies(bool, string) ([]m.MovieDto, error)
	CreateMovie(m.CreateMovieDto) (m.Movie, error)
	FindMovieByID(int) (m.MovieDto, error)
	ReplaceFieldsInMovie(int, map[string]any) (m.Movie, error)
	DeleteMovieByID(int) (bool, error)
	FindMoviesByGenre(int) ([]m.MovieDto, error)
	FindMoviesByYear(int) ([]m.Movie, error)
	FindMoviesWithActor(int) ([]m.MovieDto, error)
	FindAllActorsInMovie(int) ([]m.ActorInFilmDto, error)
	FindMovieByTitleAndYear(string, int) (m.Movie, error)
	AddActorToMovie(int, int) error
	RemoveActorFromMovie(int, int) error
	AddGenreToMovie(int, int) error
	RemoveGenreFromMovie(int, int) error
	FindGenresInMovie(int) ([]m.GenreWithoutMovies, error)
}
