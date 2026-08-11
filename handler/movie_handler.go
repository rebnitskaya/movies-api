package handler

import (
	"encoding/json"
	"net/http"
)

// GetMovies godoc
// @Summary Get all movies
// @Description Returns a paginated list of movies
// @Tags movies
// @Produce json
// @Param page query int false "Page number"
// @Param size query int false "Number of movies per page"
// @Success 200 {object} models.Movie
// @Router /movies [get]
func (h *MovieHandler) GetAllMovies(w http.ResponseWriter, r *http.Request) {
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
