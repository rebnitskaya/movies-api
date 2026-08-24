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

func (h *ActorHandler) GetAllActors(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	// Retrieve actors filtered by name
	if name := query.Get("name"); name != "" {
		actors, err := h.service.GetActorsWithName(name)
		if err != nil {
			http.Error(w, "Failed to get actors", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")

		json.NewEncoder(w).Encode(actors)
		return
	}

	actors, err := h.service.GetAllActors()
	if err != nil {
		http.Error(w, "Failed to get actors", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(actors)
}

func (h *ActorHandler) PostActor(w http.ResponseWriter, r *http.Request) {
	var actorData m.ActorDto

	err := json.NewDecoder(r.Body).Decode(&actorData)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	createdActor, err := h.service.CreateActor(actorData)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to make an actor. %s", err), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdActor)
}

func (h *ActorHandler) DeleteActor(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid actor id. %s", err), http.StatusBadRequest)
		return
	}

	ok, err := h.service.DeleteActor(idInt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Actor not found", http.StatusNotFound)
			return
		}

		http.Error(w, "Failed to delete actor: %s", http.StatusInternalServerError)
		return
	}

	if !ok {
		http.Error(w, "Failed to delete actor: %s", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *ActorHandler) GetActor(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Wrong actor id format. %s", err), http.StatusBadRequest)
		return
	}

	actor, err := h.service.GetActor(idInt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Actor not found", http.StatusNotFound)
			return
		}

		http.Error(w, fmt.Sprintf("Failed to get an actor. %s", err), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(actor)

}

func (h *ActorHandler) PatchActor(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid actor id. %s", err), http.StatusBadRequest)
		return
	}

	data := make(map[string]string)

	err = json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid request body. %s", err), http.StatusBadRequest)
		return
	}

	actor, err := h.service.PatchActor(idInt, data)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to update an actor. %s", err), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(actor)
}
