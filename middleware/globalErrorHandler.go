package middleware

import (
	"errors"
	"log"
	"movies_api/handler"
	"net/http"
)

func GlobalErrorHandler(h handler.AppHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := h(w, r); err != nil {
			handleError(w, err)
		}
	}
}

var (
	ErrMoviesNotFound = errors.New("Movie not found.")
	ErrActorsNotFound = errors.New("Actors not found.")
	ErrGenresNotFound = errors.New("Genres not found.")
	ErrInvalidInput   = errors.New("Invalid input.")
)

func handleError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrActorsNotFound):
		http.Error(w, err.Error(), http.StatusNotFound)
	case errors.Is(err, ErrActorsNotFound):
		http.Error(w, err.Error(), http.StatusNotFound)
	case errors.Is(err, ErrGenresNotFound):
		http.Error(w, err.Error(), http.StatusNotFound)
	case errors.Is(err, ErrInvalidInput):
		http.Error(w, err.Error(), http.StatusBadRequest)

	default:
		log.Printf("internal error: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
