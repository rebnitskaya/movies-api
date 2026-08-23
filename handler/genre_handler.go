package handler

import (
	"encoding/json"
	"fmt"
	m "movies_api/models"
	"net/http"
)

func (h *GenreHandler) GetAllGenres(w http.ResponseWriter, r *http.Request) {
	actors, err := h.service.GetAllGenres()
	if err != nil {
		http.Error(w, "failed to get movies", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(actors)
	if err != nil {
		return
	}
}

func (h *GenreHandler) GetGenre(w http.ResponseWriter, r *http.Request) {
	id := 0
	h.service.GetGenre(id)
	w.WriteHeader(http.StatusNoContent)
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
	id := 0
	h.service.DeleteGenre(id)
	w.WriteHeader(http.StatusNoContent)
}
