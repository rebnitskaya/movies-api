package handler

import (
	"encoding/json"
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
	name := ""
	h.service.CreateGenre(name)
	w.WriteHeader(http.StatusNoContent)
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
