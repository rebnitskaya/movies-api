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

func (h *GenreHandler) GetAllGenres(w http.ResponseWriter, r *http.Request) {

	genres, err := h.service.GetAllGenres()
	if err != nil {
		http.Error(w, "Failed to get genres", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(genres)
}

func (h *GenreHandler) GetGenre(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Wrong genre id format. %s", err), http.StatusBadRequest)
		return
	}

	genre, err := h.service.GetGenre(idInt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Genre not found", http.StatusNotFound)
			return
		}

		http.Error(w, fmt.Sprintf("Failed to get a genre. %s", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(genre)
}

func (h *GenreHandler) PostGenre(w http.ResponseWriter, r *http.Request) {
	var genre m.Genre
	err := json.NewDecoder(r.Body).Decode(&genre)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	createdGenre, err := h.service.CreateGenre(genre.Name)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to make a genre. %s", err), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdGenre)

}

func (h *GenreHandler) PatchGenre(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Wrong genre id format. %s", err), http.StatusBadRequest)
		return
	}

	var genreData m.GenreDto

	err = json.NewDecoder(r.Body).Decode(&genreData)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	genre, err := h.service.PatchGenre(idInt, genreData.Name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Genre not found", http.StatusNotFound)
			return
		}

		http.Error(w, fmt.Sprintf("Failed to update a genre. %s", err), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(genre)

}

func (h *GenreHandler) DeleteGenre(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Wrong genre id format. %s", err), http.StatusBadRequest)
		return
	}

	force := r.URL.Query().Get("force") == "true"

	deleted, err := h.service.DeleteGenre(idInt, force)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Genre not found", http.StatusNotFound)
			return
		}

		http.Error(
			w,
			fmt.Sprintf("Failed to delete genre: %s", err),
			http.StatusBadRequest,
		)
		return
	}

	if !deleted {
		http.Error(w, "Failed to delete genre", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *MovieHandler) PostGenreToMovie(w http.ResponseWriter, r *http.Request) {
	movieID, err := strconv.Atoi(r.PathValue("movieID"))
	if err != nil {
		http.Error(w, "Wrong movie id format", http.StatusBadRequest)
		return
	}
	genreID, err := strconv.Atoi(r.PathValue("genreID"))
	if err != nil {
		http.Error(w, "Wrong genre id format", http.StatusBadRequest)
		return
	}

	movie, err := h.service.AddGenreToMovie(movieID, genreID)
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("Failed to add genre to movie: %s", err),
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(movie)
}

func (h *MovieHandler) DeleteGenreFromMovie(w http.ResponseWriter, r *http.Request) {
	movieID, err := strconv.Atoi(r.PathValue("movieID"))
	if err != nil {
		http.Error(w, "Wrong movie id format", http.StatusBadRequest)
		return
	}
	genreID, err := strconv.Atoi(r.PathValue("genreID"))
	if err != nil {
		http.Error(w, "Wrong genre id format", http.StatusBadRequest)
		return
	}

	movie, err := h.service.DeleteGenreFromMovie(movieID, genreID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Movie or genre relationship not found", http.StatusNotFound)
			return
		}

		http.Error(
			w,
			fmt.Sprintf("Failed to delete genre from movie: %s", err),
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(movie)
}
