package handler

import (
	"encoding/json"
	"fmt"
	"movies_api/helper"
	m "movies_api/models"
	"net/http"
	"strconv"
)

// @Failure 400 {string} invalid input
// @Failure 404 {string} movie not found

// GetAllGenres
// @Summary Get all genres
// @Description Returns a paginated list of genres, including their associated movies.
// @Tags genres
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Number of genres per page" default(10)
// @Success 200 {object} models.GenresPaginated
// @Failure 400 {string} invalid input
// @Failure 404 {string} genre not found
// @Failure 500 {string} Internal server error
// @Router /genres [get]
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

// GetGenre
// @Summary Get genre by ID
// @Description Returns a genre by ID, including movies associated with the genre.
// @Tags genres
// @Produce json
// @Param id path int true "Genre ID"
// @Success 200 {object} models.GenreDto
// @Failure 400 {string} invalid input
// @Failure 404 {string} genre not found
// @Failure 500 {string} Internal server error
// @Router /genres/{id} [get]
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

// PostGenre
// @Summary Create a genre
// @Description Creates a new genre.
// @Tags genres
// @Accept json
// @Produce json
// @Param genre body models.GenreDto true "Genre data"
// @Success 201 {object} models.CreateGenreDto
// @Failure 400 {string} invalid input
// @Failure 404 {string} genre not found
// @Failure 500 {string} Internal server error
// @Router /genres [post]
func (h *GenreHandler) PostGenre(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	var genre m.CreateGenreDto
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

// PatchGenre
// @Summary Update a genre
// @Description Updates the name of an existing genre.
// @Tags genres
// @Accept json
// @Produce json
// @Param id path int true "Genre ID"
// @Param genre body models.GenreDto true "Genre data"
// @Success 200 {object} models.GenreDto
// @Failure 400 {string} invalid input
// @Failure 404 {string} genre not found
// @Failure 500 {string} Internal server error
// @Router /genres/{id} [patch]
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

// DeleteGenre
// @Summary Delete a genre
// @Description Deletes a genre by ID. If the genre is associated with movies, deletion requires force=true. Associated movie relationships are removed through cascading foreign keys.
// @Tags genres
// @Param id path int true "Genre ID"
// @Param force query bool true "Force deletion when the genre has associated movies"
// @Success 204
// @Failure 400 {string} invalid input
// @Failure 404 {string} genre not found
// @Failure 500 {string} Internal server error
// @Router /genres/{id} [delete]
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

// PostGenreToMovie
// @Summary Add genre to movie
// @Description Associates an existing genre with an existing movie and returns the updated movie.
// @Tags movies
// @Produce json
// @Param movieID path int true "Movie ID"
// @Param genreID path int true "Genre ID"
// @Success 200 {object} models.MovieDto
// @Failure 400 {string} invalid input
// @Failure 404 {string} genre not found
// @Failure 500 {string} Internal server error
// @Router /movies/{movieID}/genres/{genreID} [post]
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

// DeleteGenreFromMovie
// @Summary Remove genre from movie
// @Description Removes the association between a genre and a movie and returns the updated movie.
// @Tags movies
// @Produce json
// @Param movieID path int true "Movie ID"
// @Param genreID path int true "Genre ID"
// @Success 200 {object} models.MovieDto
// @Failure 400 {string} invalid input
// @Failure 404 {string} genre not found
// @Failure 500 {string} Internal server error
// @Router /movies/{movieID}/genres/{genreID} [delete]
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
