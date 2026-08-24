package models

import "errors"

type AppError struct {
	Status  int
	Message string
	Err     error
}

var (
	ErrMoviesNotFound         = errors.New("Movie not found error: ")
	ErrActorsNotFound         = errors.New("Actors not found error: ")
	ErrGenresNotFound         = errors.New("Genres not found error: ")
	ErrInvalidInput           = errors.New("Invalid input error: ")
	ErrBadRequest             = errors.New("Failed to process request: ")
	ErrMovieHasBeenMadeBefore = errors.New("This movie already has been made before: ")
)
