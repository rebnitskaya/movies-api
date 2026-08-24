package middleware

import (
	"errors"
	"log"
	"movies_api/models"
	"net/http"
)

func GlobalErrorHandler(h models.AppHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := h(w, r); err != nil {
			handleError(w, err)
		}
	}
}

func handleError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, models.ErrActorsNotFound):
		http.Error(w, err.Error(), http.StatusNotFound)
	case errors.Is(err, models.ErrMoviesNotFound):
		http.Error(w, err.Error(), http.StatusNotFound)
	case errors.Is(err, models.ErrGenresNotFound):
		http.Error(w, err.Error(), http.StatusNotFound)
	case errors.Is(err, models.ErrInvalidInput):
		http.Error(w, err.Error(), http.StatusBadRequest)
	case errors.Is(err, models.ErrBadRequest):
		http.Error(w, err.Error(), http.StatusBadRequest)
	case errors.Is(err, models.ErrMovieHasBeenMadeBefore):
		http.Error(w, err.Error(), http.StatusBadRequest)

	default:
		log.Printf("internal error: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
