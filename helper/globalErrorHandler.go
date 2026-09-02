package helper

import (
	"errors"
	"log"
	"movies_api/models"
	"net/http"
)

type AppHandler func(w http.ResponseWriter, r *http.Request) error

func GlobalErrorHandler(h AppHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := h(w, r); err != nil {
			handleError(w, err)
		}
	}
}

func handleError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, models.ErrActorNotFound):
		http.Error(w, err.Error(), http.StatusNotFound)
	case errors.Is(err, models.ErrMovieNotFound):
		http.Error(w, err.Error(), http.StatusNotFound)
	case errors.Is(err, models.ErrGenreNotFound):
		http.Error(w, err.Error(), http.StatusNotFound)
	case errors.Is(err, models.ErrInvalidInput):
		http.Error(w, err.Error(), http.StatusBadRequest)
	case errors.Is(err, models.ErrBadRequest):
		http.Error(w, err.Error(), http.StatusBadRequest)
	case errors.Is(err, models.ErrMovieHasBeenMadeBefore):
		http.Error(w, err.Error(), http.StatusBadRequest)
	case errors.Is(err, models.ErrConflict):
		http.Error(w, err.Error(), http.StatusConflict)

	default:
		log.Printf("internal error: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
