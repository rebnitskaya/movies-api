package service

import (
	"fmt"
	m "movies_api/models"
	r "movies_api/repository"
	"slices"
)

type MovieService struct {
	repo r.MovieRepository
}

func (s *MovieService) GetAllMovies() ([]m.MovieDto, error) {
	movies, err := s.repo.FindAllMovies()
	if err != nil {
		return []m.MovieDto{}, err
	}

	sortMovies(movies)

	return movies, nil
}

func (s *MovieService) FindOneMovie(id int) (m.MovieDto, error) {
	movie, err := s.repo.FindMovieByID(id)
	if err != nil {
		return m.MovieDto{}, err
	}

	return movie, nil
}

func (s *MovieService) GetAllMoviesWithGenre(genreID int) ([]m.MovieDto, error) {
	movies, err := s.repo.FindMoviesByGenre(genreID)
	if err != nil {
		return []m.MovieDto{}, err
	}

	sortMovies(movies)

	return movies, nil
}

func (s *MovieService) GetAllMoviesWithYear(year int) ([]m.Movie, error) {
	if year < 1885 || year > 2050 || year == 0 {
		return []m.Movie{}, fmt.Errorf("Invalid release year.")
	}

	movies, err := s.repo.FindMoviesByYear(year)
	if err != nil {
		return []m.Movie{}, err
	}
	return movies, nil
}

func (s *MovieService) GetAllMoviesWithActor(actorID int) ([]m.MovieDto, error) {
	movies, err := s.repo.FindMoviesWithActor(actorID)
	if err != nil {
		return []m.MovieDto{}, err
	}

	sortMovies(movies)

	return movies, nil
}

func (s *MovieService) MovieMaker(movieData m.MovieDto) (m.Movie, error) {
	_, err := movieData.Validate()
	if err != nil {
		return m.Movie{}, err
	}

	res, err := s.repo.FindMovieByTitleAndYear(movieData.Title, movieData.ReleaseYear)
	if res.Title == movieData.Title && res.ReleaseYear == movieData.ReleaseYear {
		return m.Movie{}, m.ErrMovieHasBeenMadeBefore
	}
	movie, error := s.repo.CreateMovie(movieData)
	if error != nil {
		return m.Movie{}, fmt.Errorf("Something happended during movie creation %s", err)
	}

	return movie, nil
}

func (s *MovieService) MoviePatcher(movieData m.MoviePatchDto, movieID int) (m.Movie, error) {
	fields := make(map[string]any)
	_, err := movieData.Validate()
	if err != nil {
		return m.Movie{}, err
	}

	if movieData.Title != nil {
		fields["title"] = *movieData.Title
	}

	if movieData.Duration != nil {
		fields["duration"] = *movieData.Duration
	}

	if movieData.ReleaseYear != nil {
		fields["release_year"] = *movieData.ReleaseYear
	}

	if len(fields) == 0 {
		return m.Movie{}, fmt.Errorf("No fields to update.")
	}

	movie, err := s.repo.ReplaceFieldsInMovie(movieID, fields)
	return movie, nil
}

func (s *MovieService) DeleteMovie(movieID int) (bool, error) {
	if movieID <= 0 {
		return false, fmt.Errorf("Invalid movie id.")
	}

	ok, err := s.repo.DeleteMovieByID(movieID)
	if err != nil {
		return false, err
	}

	return ok, nil
}

func (s *MovieService) FindActorsInMovie(movieID int) ([]m.ActorInFilmDto, error) {
	movies, err := s.repo.FindAllActorsInMovie(movieID)
	if err != nil {
		return []m.ActorInFilmDto{}, err
	}

	return movies, nil
}

func (s *MovieService) AddActorToMovie(movieID, actorID int) (m.MovieDto, error) {
	if movieID <= 0 {
		return m.MovieDto{}, fmt.Errorf("Invalid movie id.")
	}

	if actorID <= 0 {
		return m.MovieDto{}, fmt.Errorf("Invalid actor id.")
	}

	err := s.repo.AddActorToMovie(movieID, actorID)
	if err != nil {
		return m.MovieDto{}, err
	}

	movie, err := s.repo.FindMovieByID(movieID)
	if err != nil {
		return m.MovieDto{}, err
	}

	return movie, nil
}

func (s *MovieService) DeleteActorFromMovie(movieID, actorID int) (m.MovieDto, error) {
	if movieID <= 0 {
		return m.MovieDto{}, fmt.Errorf("Invalid movie id.")
	}

	if actorID <= 0 {
		return m.MovieDto{}, fmt.Errorf("Invalid actor id.")
	}

	err := s.repo.RemoveActorFromMovie(movieID, actorID)
	if err != nil {
		return m.MovieDto{}, err
	}

	movie, err := s.repo.FindMovieByID(movieID)
	if err != nil {
		return m.MovieDto{}, err
	}

	return movie, nil
}

func sortMovies(movies []m.MovieDto) {
	slices.SortFunc(movies, func(a, b m.MovieDto) int {
		if a.Id > b.Id {
			return 1
		} else {
			return -1
		}
	})
}
