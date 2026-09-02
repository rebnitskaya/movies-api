package handler

import (
	"encoding/json"
	"fmt"
	"movies_api/helper"
	m "movies_api/models"
	"net/http"
	"strconv"
)

// GetMovies
// @Summary Get movies
// @Description Returns a paginated list of movies. Results can optionally be filtered by genre, release year, or actor.
// @Tags movies
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param size query int false "Number of movies per page" default(5)
// @Param genre query int false "Filter movies by genre ID"
// @Param year query int false "Filter movies by release year"
// @Param actor query int false "Filter movies by actor ID"
// @Success 200 {object} models.MoviesPaginated
// @Failure 400 {string} invalid input
// @Failure 404 {string} movie not found
// @Failure 500 {string} Internal server error
// @Router /movies [get]
func (h *MovieHandler) GetMovies(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	page, limit, err := helper.GetPaginationParams(r)
	if err != nil {
		return fmt.Errorf("%w: invalid pagination", m.ErrInvalidInput)
	}

	query := r.URL.Query()
	w.Header().Set("Content-Type", "application/json")

	//not yet ready
	if genre := query.Get("genre"); genre != "" {
		genreID, err := strconv.Atoi(genre)
		if err != nil {
			return fmt.Errorf("%w: invalid genre id", m.ErrInvalidInput)
		}

		movies, err := h.service.GetAllMoviesWithGenre(genreID, page, limit, ctx)
		if err != nil {
			return err
		}

		return json.NewEncoder(w).Encode(movies)
	}

	if year := query.Get("year"); year != "" {
		year, err := strconv.Atoi(year)
		if err != nil {
			return fmt.Errorf("%w: invalid year format", m.ErrInvalidInput)
		}

		movies, err := h.service.GetAllMoviesWithYear(year, page, limit, ctx)
		if err != nil {
			return err
		}

		return json.NewEncoder(w).Encode(movies)
	}

	if actor := query.Get("actor"); actor != "" {
		actorID, err := strconv.Atoi(actor)
		if err != nil {
			return fmt.Errorf("%w: invalid actor id", m.ErrInvalidInput)
		}

		movies, err := h.service.GetAllMoviesWithActor(actorID, page, limit, ctx)
		if err != nil {
			return err
		}

		return json.NewEncoder(w).Encode(movies)
	}

	movies, err := h.service.GetAllMovies(page, limit, ctx)
	if err != nil {
		return err
	}

	return json.NewEncoder(w).Encode(movies)
}

// SearchByTitle
// @Summary Search movies by title
// @Description Returns movies whose titles contain the specified search text. The search is case-insensitive and supports pagination.
// @Tags movies
// @Produce json
// @Param title query string true "Text to search for in movie titles"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Number of movies per page" default(10)
// @Success 200 {array} models.MovieDto
// @Failure 400 {string} invalid input
// @Failure 404 {string} movie not found
// @Failure 500 {string} Internal server error
// @Router /movies/search [get]
func (h *MovieHandler) SearchByTitle(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	page, limit, err := helper.GetPaginationParams(r)
	if err != nil {
		return fmt.Errorf("%w: invalid pagination", m.ErrInvalidInput)
	}

	query := r.URL.Query()

	w.Header().Set("Content-Type", "application/json")

	//not yet ready
	if title := query.Get("title"); title != "" {
		movies, err := h.service.GetAllMoviesWithTitle(title, page, limit, ctx)
		if err != nil {
			return err
		}

		return json.NewEncoder(w).Encode(movies)
	}

	return fmt.Errorf("%w: invalid search", m.ErrBadRequest)
}

// GetMovie
// @Summary Get movie by ID
// @Description Returns a movie by its ID, including its associated actors and genres.
// @Tags movies
// @Produce json
// @Param id path int true "Movie ID"
// @Success 200 {object} models.MovieDto
// @Failure 400 {string} invalid input
// @Failure 404 {string} movie not found
// @Failure 500 {string} Internal server error
// @Router /movies/{id} [get]
func (h *MovieHandler) GetMovie(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	movieId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		return fmt.Errorf("%w: invalid movie id", m.ErrInvalidInput)
	}

	movie, err := h.service.FindOneMovie(movieId, ctx)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")

	return json.NewEncoder(w).Encode(movie)
}

// PostMovie
// @Summary Create a movie
// @Description Creates a new movie and optionally associates it with existing actors and genres.
// @Tags movies
// @Accept json
// @Produce json
// @Param movie body models.CreateMovieDto true "Movie data"
// @Success 201 {object} models.MovieDto
// @Failure 400 {string} invalid input
// @Failure 404 {string} movie not found
// @Failure 500 {string} Internal server error
// @Router /movies [post]
func (h *MovieHandler) PostMovie(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	movieDto := m.CreateMovieDto{}

	err := json.NewDecoder(r.Body).Decode(&movieDto)
	if err != nil {
		return fmt.Errorf("%w: something is wrong with incoming data: %s", m.ErrInvalidInput, err)
	}

	movie, err := h.service.MovieMaker(movieDto, ctx)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	return json.NewEncoder(w).Encode(movie)
}

// PatchMovie
// @Summary Update a movie
// @Description Updates one or more fields of an existing movie. Only provided fields are modified.
// @Tags movies
// @Accept json
// @Produce json
// @Param id path int true "Movie ID"
// @Param movie body models.MoviePatchDto true "Movie fields to update"
// @Success 200 {object} models.MovieDto
// @Failure 400 {string} invalid input
// @Failure 404 {string} movie not found
// @Failure 500 {string} Internal server error
// @Router /movies/{id} [patch]
func (h *MovieHandler) PatchMovie(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	movieId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		return fmt.Errorf("%w: invalid movie id", m.ErrInvalidInput)
	}
	data := m.MoviePatchDto{}

	err = json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		return fmt.Errorf("%w: something wrong with incoming data: %s", m.ErrInvalidInput, err)
	}

	movie, err := h.service.MoviePatcher(data, movieId, ctx)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(movie)
}

// DeleteMovie
// @Summary Delete a movie
// @Description Deletes a movie by ID. Deletion requires the force query parameter to be set to true. Associated actor and genre relationships are removed through cascading foreign keys.
// @Tags movies
// @Param id path int true "Movie ID"
// @Param force query bool true "Must be true to allow deletion"
// @Success 204
// @Failure 400 {string} invalid input
// @Failure 404 {string} movie not found
// @Failure 500 {string} Internal server error
// @Router /movies/{id} [delete]
func (h *MovieHandler) DeleteMovie(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	movieId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		return fmt.Errorf("%w: invalid movie id", m.ErrInvalidInput)
	}

	force := r.URL.Query().Get("force") == "true"

	ok, err := h.service.DeleteMovie(movieId, force, ctx)
	if err != nil {
		return err
	}

	if !ok {
		return fmt.Errorf("Failed to delete movie.")
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}

// PostActorToMovie
// @Summary Add actor to movie
// @Description Associates an existing actor with an existing movie and returns the updated movie.
// @Tags movies
// @Produce json
// @Param movieID path int true "Movie ID"
// @Param actorID path int true "Actor ID"
// @Success 200 {object} models.MovieDto
// @Failure 400 {string} invalid input
// @Failure 404 {string} movie not found
// @Failure 500 {string} Internal server error
// @Router /movies/{movieID}/actors/{actorID} [post]
func (h *MovieHandler) PostActorToMovie(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	movieID, err := strconv.Atoi(r.PathValue("movieID"))
	if err != nil {
		return fmt.Errorf("%w: invalid movie id", m.ErrInvalidInput)
	}
	actorID, err := strconv.Atoi(r.PathValue("actorID"))
	if err != nil {
		return fmt.Errorf("%w: invalid actor id", m.ErrInvalidInput)
	}

	movie, err := h.service.AddActorToMovie(movieID, actorID, ctx)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(movie)
}

// DeleteActorFromMovie
// @Summary Remove actor from movie
// @Description Removes the association between an actor and a movie and returns the movie.
// @Tags movies
// @Produce json
// @Param movieID path int true "Movie ID"
// @Param actorID path int true "Actor ID"
// @Success 200 {object} models.MovieDto
// @Failure 400 {string} invalid input
// @Failure 404 {string} movie not found
// @Failure 500 {string} Internal server error
// @Router /movies/{movieID}/actors/{actorID} [delete]
func (h *MovieHandler) DeleteActorFromMovie(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	movieID, err := strconv.Atoi(r.PathValue("movieID"))
	if err != nil {
		return fmt.Errorf("%w: invalid movie id", m.ErrInvalidInput)
	}
	actorID, err := strconv.Atoi(r.PathValue("actorID"))
	if err != nil {
		return fmt.Errorf("%w: invalid actor id", m.ErrInvalidInput)
	}

	movie, err := h.service.DeleteActorFromMovie(movieID, actorID, ctx)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(movie)
}

// GetActorsInMovie
// @Summary Get actors in a movie
// @Description Returns all actors associated with the specified movie.
// @Tags movies
// @Produce json
// @Param id path int true "Movie ID"
// @Success 200 {array} models.ActorInFilmDto
// @Failure 400 {string} invalid input
// @Failure 404 {string} movie not found
// @Failure 500 {string} Internal server error
// @Router /movies/{id}/actors [get]
func (h *MovieHandler) GetActorsInMovie(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	movieID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		return fmt.Errorf("%w: invalid movie id", m.ErrInvalidInput)
	}

	movies, err := h.service.FindActorsInMovie(movieID, ctx)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(movies)
}
