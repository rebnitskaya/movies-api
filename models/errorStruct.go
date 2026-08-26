package models

import "errors"

type AppError struct {
	Status  int
	Message string
	Err     error
}

var (
	ErrMovieNotFound          = errors.New("movie not found")
	ErrActorNotFound          = errors.New("actors not found")
	ErrGenreNotFound          = errors.New("genres not found")
	ErrInvalidInput           = errors.New("invalid input")
	ErrBadRequest             = errors.New("failed to process request")
	ErrMovieHasBeenMadeBefore = errors.New("this movie already has been made before")
	ErrInternalIssue          = errors.New("error during program execution")
)
