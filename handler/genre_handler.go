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
		http.Error(w, "failed to get movies", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(genres)
	if err != nil {
		return
	}
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

		http.Error(w, fmt.Sprintf("Failed to get a genre. %s", err), http.StatusBadRequest)
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
	id := 0
	name := ""
	h.service.PatchGenre(id, name)
	w.WriteHeader(http.StatusNoContent)
}

func (h *GenreHandler) DeleteGenre(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Wrong genre id format. %s", err), http.StatusBadRequest)
		return
	}

	ok, err := h.service.DeleteGenre(idInt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Genre not found", http.StatusNotFound)
			return
		}

		http.Error(w, "Failed to delete genre: %s", http.StatusInternalServerError)
		return
	}

	if !ok {
		http.Error(w, "Failed to genre movie: %s", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
