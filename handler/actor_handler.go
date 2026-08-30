package handler

import (
	"encoding/json"
	"fmt"
	"movies_api/helper"
	m "movies_api/models"
	"net/http"
	"strconv"
)

// GetAllActors
// @Summary Get all actors
// @Description Returns a paginated list of actors. Use the optional name query parameter to filter actors by name.
// @Tags actors
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param size query int false "Number of actors per page" default(5)
// @Param name query string false "Filter actors by name"
// @Success 200 {object} models.ActorsPaginated
// @Failure 400 {string} string "invalid pagination parameters"
// @Failure 500 {string} string "Internal server error"
// @Router /actors [get]
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

// PostActor
// @Summary Create an actor
// @Description Creates a new actor.
// @Tags actors
// @Accept json
// @Produce json
// @Param actor body models.ActorDto true "Actor data"
// @Success 201 {object} models.ActorInFilmDto
// @Failure 400 {string} string "Invalid request data"
// @Failure 500 {string} string "Internal server error"
// @Router /actors [post]
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

// DeleteActor
// @Summary Delete an actor
// @Description Deletes an actor by ID. By default, an actor associated with movies cannot be deleted. Use force=true to delete the actor and its movie associations.
// @Tags actors
// @Param id path int true "Actor ID"
// @Param force query bool false "Force deletion even if the actor is associated with movies"
// @Success 204
// @Failure 400 {string} invalid input
// @Failure 404 {string} actor not found
// @Failure 500 {string} Internal server error
// @Router /actors/{id} [delete]
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

// GetActor
// @Summary Get an actor
// @Description Returns an actor by ID.
// @Tags actors
// @Produce json
// @Param id path int true "Actor ID"
// @Success 200 {object} models.ActorWithoutMoviesDto
// @Failure 400 {string} invalid input
// @Failure 404 {string} actor not found
// @Failure 500 {string} Internal server error
// @Router /actors/{id} [get]
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

// PatchActor
// @Summary Update an actor
// @Description Updates one or more fields of an existing actor.
// @Tags actors
// @Accept json
// @Produce json
// @Param id path int true "Actor ID"
// @Param actor body map[string]string true "Fields to update"
// @Success 200 {object} models.ActorInFilmDto
// @Failure 400 {string} invalid input
// @Failure 404 {string} actor not found
// @Failure 500 {string} Internal server error
// @Router /actors/{id} [patch]
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
