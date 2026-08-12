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
func (h *ActorHandler) GetAllActors(w http.ResponseWriter, r *http.Request) {
	actors, err := h.service.GetAllActors()
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

// bith_date >= '1895-12-28'
// AND birth_date <=date('now')
