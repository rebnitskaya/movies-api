package handler

import (
	"encoding/json"
	"fmt"
	m "movies_api/models"
	"net/http"
	"strconv"
)

func (h *ActorHandler) GetAllActors(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	// Retrieve actors filtered by name
	if name := query.Get("name"); name != "" {
		h.service.GetActorsWithName(name)

		return
	}

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
	var actorData m.Actor

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

func (h *ActorHandler) GetActor(w http.ResponseWriter, r *http.Request) {
	actorId := 0
	h.service.GetActor(actorId)
	w.WriteHeader(http.StatusNoContent)
}

func (h *ActorHandler) PatchActor(w http.ResponseWriter, r *http.Request) {
	data := ""
	h.service.PatchActor(data)
	w.WriteHeader(http.StatusNoContent)
}
