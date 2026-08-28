package handler

import (
	"encoding/json"
	"fmt"
	"movies_api/helper"
	m "movies_api/models"
	"net/http"
	"strconv"
)

func (h *ActorHandler) GetAllActors(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	page, limit, err := helper.GetPaginationParams(r)
	if err != nil {
		return fmt.Errorf("%w: invalid pagination", m.ErrInvalidInput)
	}

	query := r.URL.Query()

	// Retrieve actors filtered by name
	if name := query.Get("name"); name != "" {
		actors, err := h.service.GetActorsWithName(name, ctx)
		if err != nil {
			return err
		}
		w.Header().Set("Content-Type", "application/json")

		return json.NewEncoder(w).Encode(actors)
	}

	actors, err := h.service.GetAllActors(page, limit, ctx)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")

	return json.NewEncoder(w).Encode(actors)
}

func (h *ActorHandler) PostActor(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	var actorData m.ActorDto

	err := json.NewDecoder(r.Body).Decode(&actorData)

	if err != nil {
		return fmt.Errorf("%w: something is wrong with incoming data: %s", m.ErrInvalidInput, err)
	}

	createdActor, err := h.service.CreateActor(actorData, ctx)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	return json.NewEncoder(w).Encode(createdActor)
}

func (h *ActorHandler) DeleteActor(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	id := r.PathValue("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return fmt.Errorf("%w: invalid actor id", m.ErrInvalidInput)
	}

	force := r.URL.Query().Get("force") == "true"

	_, err = h.service.DeleteActor(idInt, force, ctx)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusNoContent)

	return nil
}

func (h *ActorHandler) GetActor(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	id := r.PathValue("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return fmt.Errorf("%w: invalid actor id", m.ErrInvalidInput)
	}

	actor, err := h.service.GetActor(idInt, ctx)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")

	return json.NewEncoder(w).Encode(actor)
}

func (h *ActorHandler) PatchActor(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	id := r.PathValue("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return fmt.Errorf("%w: invalid actor id", m.ErrInvalidInput)
	}

	data := make(map[string]string)

	err = json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		return fmt.Errorf("%w: something is wrong with incoming data: %s", m.ErrInvalidInput, err)
	}

	actor, err := h.service.PatchActor(idInt, data, ctx)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(actor)
}
