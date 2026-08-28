package handler

import (
	"encoding/json"
	"fmt"
	"movies_api/helper"
	m "movies_api/models"
	"net/http"
	"strconv"
)

func (h *GenreHandler) GetAllGenres(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	page, limit, err := helper.GetPaginationParams(r)
	if err != nil {
		return fmt.Errorf("%w: invalid pagination", m.ErrInvalidInput)
	}

	genres, err := h.service.GetAllGenres(page, limit, ctx)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")

	return json.NewEncoder(w).Encode(genres)
}

func (h *GenreHandler) GetGenre(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	id := r.PathValue("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return fmt.Errorf("%w: invalid genre id", m.ErrInvalidInput)
	}

	genre, err := h.service.GetGenre(idInt, ctx)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(genre)
}

func (h *GenreHandler) PostGenre(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	var genre m.Genre
	err := json.NewDecoder(r.Body).Decode(&genre)
	if err != nil {
		return fmt.Errorf("%w: invalid request body", m.ErrBadRequest)
	}

	createdGenre, err := h.service.CreateGenre(genre.Name, ctx)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	return json.NewEncoder(w).Encode(createdGenre)
}

func (h *GenreHandler) PatchGenre(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	id := r.PathValue("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return fmt.Errorf("%w: invalid genre id", m.ErrInvalidInput)
	}

	var genreData m.GenreDto

	err = json.NewDecoder(r.Body).Decode(&genreData)
	if err != nil {
		return fmt.Errorf("%w: invalid request body", m.ErrBadRequest)
	}

	genre, err := h.service.PatchGenre(idInt, genreData.Name, ctx)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(genre)
}

func (h *GenreHandler) DeleteGenre(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	id := r.PathValue("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return fmt.Errorf("%w: invalid genre id", m.ErrInvalidInput)
	}

	force := r.URL.Query().Get("force") == "true"

	_, err = h.service.DeleteGenre(idInt, force, ctx)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusNoContent)

	return nil
}

func (h *MovieHandler) PostGenreToMovie(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	movieID, err := strconv.Atoi(r.PathValue("movieID"))
	if err != nil {
		return fmt.Errorf("%w: invalid movie id", m.ErrInvalidInput)
	}
	genreID, err := strconv.Atoi(r.PathValue("genreID"))
	if err != nil {
		return fmt.Errorf("%w: invalid genre id", m.ErrInvalidInput)
	}

	movie, err := h.service.AddGenreToMovie(movieID, genreID, ctx)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	return json.NewEncoder(w).Encode(movie)
}

func (h *MovieHandler) DeleteGenreFromMovie(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	movieID, err := strconv.Atoi(r.PathValue("movieID"))
	if err != nil {
		return fmt.Errorf("%w: invalid movie id", m.ErrInvalidInput)
	}
	genreID, err := strconv.Atoi(r.PathValue("genreID"))
	if err != nil {
		return fmt.Errorf("%w: invalid genre id", m.ErrInvalidInput)
	}

	movie, err := h.service.DeleteGenreFromMovie(movieID, genreID, ctx)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(movie)
}
