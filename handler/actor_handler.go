package handler

import (
	"encoding/json"
	"fmt"
	"movies_api/models"
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
func (h *ActorHandler) getAllActors(w http.ResponseWriter, r *http.Request) {
	actors, err := h.service.GetAllActors()
	if err != nil {
		http.Error(w, "Failed to get movies", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(actors)
	if err != nil {
		return
	}
}

func (h *ActorHandler) postActor(w http.ResponseWriter, r *http.Request) {
	var actorData models.Actor

	err := json.NewDecoder(r.Body).Decode(&actorData)

	actors, err := h.service.CreateActor(actorData)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to make an actor. %s", err), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(actors)
	if err != nil {
		return
	}
}

func (h *ActorHandler) ActorHandlerRouter(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.getAllActors(w, r)

	case http.MethodPost:
		h.postActor(w, r)

	default:
		http.Error(w, "method is not allowed", http.StatusMethodNotAllowed)
	}
}
