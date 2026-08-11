package handlers

import (
	"encoding/json"
	"net/http"

	"movies_api/models"
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

func CreateMovieHandler(w http.ResponseWriter, r *http.Request) {
	var newMovie models.Movie

	err := json.NewDecoder(r.Body).Decode(&newMovie)
	if err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if newMovie.Title == "" {
		http.Error(w, "title is missing", http.StatusBadRequest)
		return
	}

	if newMovie.ReleaseYear <= 0 {
		http.Error(w, "release year is missing or invalid number", http.StatusBadRequest)
		return
	}

	if newMovie.Duration <= 0 {
		http.Error(w, "duration is missing or invalid number", http.StatusBadRequest)
		return
	}

	newMovie.Id = int64(len(movies) + 1)

	if newMovie.Genres == nil {
		newMovie.Genres = []models.Genre{}
	}

	if newMovie.Actors == nil {
		newMovie.Actors = []models.Actor{}
	}

	movies = append(movies, newMovie)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	err = json.NewEncoder(w).Encode(newMovie)
	if err != nil {
		http.Error(w, "failed to add a movie", http.StatusInternalServerError)
		return
	}
}

func MoviesHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		GetMoviesHandler(w, r)

	case http.MethodPost:
		CreateMovieHandler(w, r)

	default:
		http.Error(w, "method is not allowed", http.StatusMethodNotAllowed)
	}
}
