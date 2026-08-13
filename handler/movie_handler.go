package handler

import (
	"encoding/json"
	"movies_api/repository"
	"net/http"
)

func (h *MovieHandler) GetMovies(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	if genre := query.Get("genre"); genre != "" {
		// Retrieve movies filtered by genre
		genreId := 0
		h.service.GetAllMoviesWithGenre(genreId)
		return
	}

	if year := query.Get("year"); year != "" {
		// Retrieve movies filtered by release year
		year := 0
		h.service.GetAllMoviesWithYear(year)
		return
	}

	if actor := query.Get("actor"); actor != "" {
		// Retrieve movies that the actor with the given id has starred in
		actorId := 0
		h.service.GetAllMoviesWithActor(actorId)
		return
	}

	movies, err := h.service.GetAllMovies()
	if err != nil {
		http.Error(w, "failed to get movies", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(movies)
	if err != nil {
		return
	}
}

func (h *MovieHandler) GetMovie(w http.ResponseWriter, r *http.Request) {
	movieId := 0
	h.service.FindOneMovie(movieId)
	w.WriteHeader(http.StatusNoContent)

}

func (h *MovieHandler) PostMovie(w http.ResponseWriter, r *http.Request) {
	data := repository.Movie{}
	h.service.MovieMaker(data)
	w.WriteHeader(http.StatusNoContent)
}

func (h *MovieHandler) PatchMovie(w http.ResponseWriter, r *http.Request) {
	data := repository.Movie{}
	h.service.MoviePatcher(data)
	w.WriteHeader(http.StatusNoContent)
}

func (h *MovieHandler) DeleteMovie(w http.ResponseWriter, r *http.Request) {
	movieId := 0
	h.service.DeleteMovie(movieId)
	w.WriteHeader(http.StatusNoContent)
}

// Retrieve all actors starring in a movie
func (h *MovieHandler) GetActorsInMovie(w http.ResponseWriter, r *http.Request) {
	movieId := 0
	h.service.FindActorsInMovie(movieId)
	w.WriteHeader(http.StatusNoContent)
}
