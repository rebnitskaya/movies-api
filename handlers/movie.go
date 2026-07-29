package handlers

import (
	"encoding/json"
	"net/http"

	"movies-api/models"
)

// to check
var movies = []models.Movie{
	{
		Id:          1,
		Title:       "Titanic",
		ReleaseYear: 1997,
		Duration:    194,
		Genres:      []models.Genre{},
		Actors:      []models.Actor{},
	},
	{
		Id:          2,
		Title:       "Interstellar",
		ReleaseYear: 2014,
		Duration:    169,
		Genres:      []models.Genre{},
		Actors:      []models.Actor{},
	},
}

func GetMoviesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err := json.NewEncoder(w).Encode(movies)
	if err != nil {
		http.Error(w, "failed to show movies", http.StatusInternalServerError)
	}
}
