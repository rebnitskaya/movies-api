package handler

import (
	"encoding/json"
	"fmt"
	m "movies_api/models"
	"net/http"
	"strconv"
)

func (h *MovieHandler) GetMovies(w http.ResponseWriter, r *http.Request) error {
	query := r.URL.Query()

	w.Header().Set("Content-Type", "application/json")

	//not yet ready
	if genre := query.Get("genre"); genre != "" {
		genreID, err := strconv.Atoi(genre)
		if err != nil {
			return fmt.Errorf("%w: invalid genre id", m.ErrInvalidInput)
		}

		movies, err := h.service.GetAllMoviesWithGenre(genreID)
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

		movies, err := h.service.GetAllMoviesWithYear(year)
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

		movies, err := h.service.GetAllMoviesWithActor(actorID)
		if err != nil {
			return err
		}

		return json.NewEncoder(w).Encode(movies)
	}

	movies, err := h.service.GetAllMovies()
	if err != nil {
		return err
	}

	return json.NewEncoder(w).Encode(movies)
}

func (h *MovieHandler) SearchByTitle(w http.ResponseWriter, r *http.Request) error {
	query := r.URL.Query()

	w.Header().Set("Content-Type", "application/json")

	//not yet ready
	if title := query.Get("title"); title != "" {
		movies, err := h.service.GetAllMoviesWithTitle(title)
		if err != nil {
			return err
		}

		return json.NewEncoder(w).Encode(movies)
	}

	return fmt.Errorf("%w: invalid search", m.ErrBadRequest)
}

func (h *MovieHandler) GetMovie(w http.ResponseWriter, r *http.Request) error {
	movieId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		return fmt.Errorf("%w: invalid movie id", m.ErrInvalidInput)
	}

	movie, err := h.service.FindOneMovie(movieId)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")

	return json.NewEncoder(w).Encode(movie)
}

func (h *MovieHandler) PostMovie(w http.ResponseWriter, r *http.Request) error {
	movieDto := m.CreateMovieDto{}

	err := json.NewDecoder(r.Body).Decode(&movieDto)
	if err != nil {
		return fmt.Errorf("%w: something is wrong with incoming data: %s", m.ErrInvalidInput, err)
	}

	movie, err := h.service.MovieMaker(movieDto)
	if err != nil {
		return err
	}

	message := fmt.Sprintf("New movie created. Movie id: %d", movie.Id)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	return json.NewEncoder(w).Encode(message)
}

func (h *MovieHandler) PatchMovie(w http.ResponseWriter, r *http.Request) error {
	movieId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		return fmt.Errorf("%w: invalid movie id", m.ErrInvalidInput)
	}
	data := m.MoviePatchDto{}

	err = json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		return fmt.Errorf("%w: something wrong with incoming data: %s", m.ErrInvalidInput, err)
	}

	movie, err := h.service.MoviePatcher(data, movieId)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "aplication/json")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(movie)
}

func (h *MovieHandler) DeleteMovie(w http.ResponseWriter, r *http.Request) error {
	movieId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		return fmt.Errorf("%w: invalid movie id", m.ErrInvalidInput)
	}

	force := r.URL.Query().Get("force") == "true"

	ok, err := h.service.DeleteMovie(movieId, force)
	if err != nil {
		return err
	}

	if !ok {
		return fmt.Errorf("Failed to delete movie.")
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}

func (h *MovieHandler) PostActorToMovie(w http.ResponseWriter, r *http.Request) error {
	movieID, err := strconv.Atoi(r.PathValue("movieID"))
	if err != nil {
		return fmt.Errorf("%w: invalid movie id", m.ErrInvalidInput)
	}
	actorID, err := strconv.Atoi(r.PathValue("actorID"))
	if err != nil {
		return fmt.Errorf("%w: invalid actor id", m.ErrInvalidInput)
	}

	movie, err := h.service.AddActorToMovie(movieID, actorID)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(movie)
}

func (h *MovieHandler) DeleteActorFromMovie(w http.ResponseWriter, r *http.Request) error {
	movieID, err := strconv.Atoi(r.PathValue("movieID"))
	if err != nil {
		return fmt.Errorf("%w: invalid movie id", m.ErrInvalidInput)
	}
	actorID, err := strconv.Atoi(r.PathValue("actorID"))
	if err != nil {
		return fmt.Errorf("%w: invalid actor id", m.ErrInvalidInput)
	}

	movie, err := h.service.DeleteActorFromMovie(movieID, actorID)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(movie)
}

func (h *MovieHandler) GetActorsInMovie(w http.ResponseWriter, r *http.Request) error {
	movieID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		return fmt.Errorf("%w: invalid movie id", m.ErrInvalidInput)
	}

	movies, err := h.service.FindActorsInMovie(movieID)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "aplication/json")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(movies)
}
