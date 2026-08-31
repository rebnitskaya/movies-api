package repository

import (
	"context"
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
	FindAllActors(int, int, context.Context) ([]m.Actor, error)
	CreateActor(m.Actor, context.Context) (m.Actor, error)
	FindActorByNameAndBirthDate(string, string, context.Context) (m.Actor, error)
	DeleteActorByID(int, context.Context) (bool, error)
	FindActorByID(int, context.Context) (m.Actor, error)
	ReplaceFieldsInActor(int, map[string]string, context.Context) (m.Actor, error)
	FindActorsByName(string, context.Context) ([]m.ActorWithoutMoviesDto, error)
	CountActors(context.Context) (int, error)
}

type GenreRepository interface {
	FindAllGenres(int, int, context.Context) ([]m.Genre, error)
	CreateGenre(m.Genre, context.Context) (m.Genre, error)
	FindGenreByID(int, context.Context) (m.Genre, error)
	ReplaceFieldsInGenre(int, string, context.Context) (m.Genre, error)
	DeleteGenreByID(int, context.Context) (bool, error)
	FindGenreByName(string, context.Context) (m.Genre, error)
	CountGenres(context.Context) (int, error)
}

type MovieRepository interface {
	FindAllMovies(bool, string, int, int, context.Context) ([]m.MovieDto, error)
	CreateMovie(m.CreateMovieDto, context.Context) (m.Movie, error)
	FindMovieByID(int, context.Context) (m.MovieDto, error)
	ReplaceFieldsInMovie(int, map[string]any, context.Context) (m.Movie, error)
	DeleteMovieByID(int, context.Context) (bool, error)
	FindMoviesByGenre(int, int, int, context.Context) ([]m.MovieDto, error)
	FindMoviesByYear(int, int, int, context.Context) ([]m.MovieDto, error)
	FindMoviesWithActor(int, int, int, context.Context) ([]m.MovieDto, error)
	FindAllActorsInMovie(int, context.Context) ([]m.ActorInFilmDto, error)
	FindMovieByTitleAndYear(string, int, context.Context) (m.Movie, error)
	AddActorToMovie(int, int, context.Context) error
	RemoveActorFromMovie(int, int, context.Context) error
	AddGenreToMovie(int, int, context.Context) error
	RemoveGenreFromMovie(int, int, context.Context) error
	FindGenresInMovie(int, context.Context) ([]m.GenreWithoutMovies, error)
	CountMovies(context.Context) (int, error)
	CountMoviesByGenre(int, context.Context) (int, error)
	CountMoviesByYear(int, context.Context) (int, error)
	CountMoviesByActor(int, context.Context) (int, error)
}
