package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	m "movies_api/models"
	"net/http"
	"strconv"
)

func (h *MovieHandler) GetMovies(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	w.Header().Set("Content-Type", "application/json")

	//not yet ready
	if genre := query.Get("genre"); genre != "" {
		genreID, err := strconv.Atoi(genre)
		if err != nil {
			http.Error(w, "Wrong genre id format", http.StatusBadRequest)
			return
		}

		movies, err := h.service.GetAllMoviesWithGenre(genreID)
		if err != nil {
			http.Error(w, "Failed to get movies", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(movies)
		return
	}

	if year := query.Get("year"); year != "" {
		year, err := strconv.Atoi(year)
		if err != nil {
			http.Error(w, "Wrong year format.", http.StatusBadRequest)
			return
		}

		movies, err := h.service.GetAllMoviesWithYear(year)
		if err != nil {
			http.Error(w, fmt.Sprintf("Something happend: %s", err), http.StatusBadRequest)
			return
		}

		json.NewEncoder(w).Encode(movies)
		return
	}

	if actor := query.Get("actor"); actor != "" {
		actorID, err := strconv.Atoi(actor)
		if err != nil {
			http.Error(w, "Wrong actor id format", http.StatusBadRequest)
			return
		}

		movies, err := h.service.GetAllMoviesWithActor(actorID)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to get movies: %s", err), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(movies)
		return
	}

	movies, err := h.service.GetAllMovies()
	if err != nil {
		http.Error(w, "Failed to get movies", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(movies)
}

func (h *MovieHandler) GetMovie(w http.ResponseWriter, r *http.Request) {
	movieId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Wrong movie id format", http.StatusBadRequest)
		return
	}

	movie, err := h.service.FindOneMovie(movieId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Movie not found", http.StatusNotFound)
			return
		}

		http.Error(w, "Failed to get movie", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(movie)
	if err != nil {
		return
	}
}

func (h *MovieHandler) PostMovie(w http.ResponseWriter, r *http.Request) {
	movieDto := m.MovieDto{}

	err := json.NewDecoder(r.Body).Decode(&movieDto)
	if err != nil {
		http.Error(w, fmt.Sprintf("Something wrong with incoming data. %s", err), http.StatusBadRequest)
		return
	}

	m, error := h.service.MovieMaker(movieDto)
	if error != nil {
		http.Error(w, fmt.Sprintf("Failed to create a film. %s", error), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(m)
}

func (h *MovieHandler) PatchMovie(w http.ResponseWriter, r *http.Request) {
	movieId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Wrong movie id format", http.StatusBadRequest)
		return
	}
	data := m.MoviePatchDto{}

	err = json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, fmt.Sprintf("Something wrong with incoming data. %s", err), http.StatusBadRequest)
		return
	}

	movie, err := h.service.MoviePatcher(data, movieId)
	if err != nil {
		http.Error(w, fmt.Sprintf("Something wrong during request execution: %s", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "aplication/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(movie)
}

func (h *MovieHandler) DeleteMovie(w http.ResponseWriter, r *http.Request) {
	movieId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Wrong movie id format", http.StatusBadRequest)
		return
	}

	ok, err := h.service.DeleteMovie(movieId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Movie not found", http.StatusNotFound)
			return
		}

		http.Error(w, "Failed to delete movie: %s", http.StatusInternalServerError)
		return
	}

	if !ok {
		http.Error(w, "Failed to delete movie: %s", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *MovieHandler) PostActorToMovie(w http.ResponseWriter, r *http.Request) {
	movieID, err := strconv.Atoi(r.PathValue("movieID"))
	if err != nil {
		http.Error(w, "Wrong movie id format", http.StatusBadRequest)
		return
	}
	actorID, err := strconv.Atoi(r.PathValue("actorID"))
	if err != nil {
		http.Error(w, "Wrong actor id format", http.StatusBadRequest)
		return
	}

	movie, err := h.service.AddActorToMovie(movieID, actorID)
	if err != nil {
		http.Error(w, "Failed to add actor to movie: %s", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "aplication/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(movie)
}

func (h *MovieHandler) DeleteActorFromMovie(w http.ResponseWriter, r *http.Request) {
	movieID, err := strconv.Atoi(r.PathValue("movieID"))
	if err != nil {
		http.Error(w, "Wrong movie id format", http.StatusBadRequest)
		return
	}
	actorID, err := strconv.Atoi(r.PathValue("actorID"))
	if err != nil {
		http.Error(w, "Wrong actor id format", http.StatusBadRequest)
		return
	}

	movie, err := h.service.DeleteActorFromMovie(movieID, actorID)
	if err != nil {
		http.Error(w, "Failed to delete actor from movie: %s", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "aplication/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(movie)
}

// not yet ready
// Retrieve all actors starring in a movie
func (h *MovieHandler) GetActorsInMovie(w http.ResponseWriter, r *http.Request) {
	movieID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Wrong movie id format", http.StatusBadRequest)
		return
	}

	movies, err := h.service.FindActorsInMovie(movieID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Movies with actor not found", http.StatusNotFound)
			return
		}

		http.Error(w, "Failed to delete actor from movie: %s", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "aplication/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(movies)
}
