package handler

import (
	"encoding/json"
	"fmt"
	"movies_api/models"
	"net/http"
	"strconv"
)

// getAllActors docs for swagger
// @Summary Get all actors
// @Description Returns a paginated list of all actors
// @Tags actors
// @Produce json
// @Param page query int false "Page number"
// @Param size query int false "Number of movies per page"
// @Success 200 {object} models.Actor
// @Router /actors [get]
func (h *ActorHandler) GetAllActors(w http.ResponseWriter, r *http.Request) {
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

func (h *ActorHandler) PostActor(w http.ResponseWriter, r *http.Request) {
	var actorData models.Actor

	err := json.NewDecoder(r.Body).Decode(&actorData)

	ok, err := h.service.CreateActor(actorData)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to make an actor. %s", err), http.StatusBadRequest)
		return
	}

	if ok {
		w.WriteHeader(http.StatusCreated)
	}
}

func (h *ActorHandler) DeleteActor(w http.ResponseWriter, r *http.Request) {
	fmt.Println("deleteActor")

	id := r.PathValue("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to delete an actor. %s", err), http.StatusBadRequest)
		return
	}

	_, err = h.service.DeleteActor(idInt)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to delete an actor. %s", err), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
