package service

import (
	"fmt"
	m "movies_api/models"
	r "movies_api/repository"
)

type MovieService struct {
	repo r.MovieRepository
}

func (s *MovieService) GetAllMovies() ([]m.Movie, error) {
	res, err := s.repo.FindAllMovies()
	if err != nil {
		return []m.Movie{}, err
	}
	return res, nil
}

func (s *MovieService) FindOneMovie(id int) (m.Movie, error) {
	movie, err := s.repo.FindMovieByID(id)
	if err != nil {
		return m.Movie{}, err
	}

	return movie, nil
}

func (s *MovieService) GetAllMoviesWithGenre(genreID int) ([]m.Movie, error) {
	s.repo.FindMoviesByGenre(genreID)
	return []m.Movie{}, nil
}

func (s *MovieService) GetAllMoviesWithYear(year int) ([]m.Movie, error) {
	s.repo.FindMoviesByYear(year)
	return []m.Movie{}, nil
}

func (s *MovieService) GetAllMoviesWithActor(actorID int) ([]m.Movie, error) {
	s.repo.FindMoviesWithActor(actorID)
	return []m.Movie{}, nil
}

func (s *MovieService) GetActorsInMovie(movieId int) ([]m.Actor, error) {
	s.repo.FindAllActorsInMovie(movieId)
	return []m.Actor{}, nil
}

func (s *MovieService) MovieMaker(movieData m.MovieDto) (m.Movie, error) {
	_, err := movieData.Validate()
	if err != nil {
		return m.Movie{}, err
	}

	res, err := s.repo.FindMovieByTitleAndYear(movieData.Title, movieData.ReleaseYear)
	if res.Title == movieData.Title && res.ReleaseYear == movieData.ReleaseYear {
		return m.Movie{}, fmt.Errorf("This movie already has been made before.")
	}
	movie, error := s.repo.CreateMovie(movieData)
	if error != nil {
		return m.Movie{}, fmt.Errorf("Something happended during movie creation %s", err)
	}

	return movie, nil
}

func (s *MovieService) MoviePatcher(movieData m.Movie) (m.Movie, error) {
	fields := make(map[string]string)
	s.repo.ReplaceFieldsInMovie(movieData.Id, fields)
	return m.Movie{}, nil
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

func (s *MovieService) FindActorsInMovie(movieID int) ([]m.Actor, error) {
	s.repo.FindAllActorsInMovie(movieID)
	return []m.Actor{}, nil
}
