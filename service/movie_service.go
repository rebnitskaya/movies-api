package service

import (
	r "movies_api/repository"
)

type MovieService struct {
	repo r.MovieRepository
}

func (s *MovieService) GetAllMovies() ([]r.Movie, error) {
	s.repo.FindAllMovies()
	return []r.Movie{}, nil
}

func (s *MovieService) FindOneMovie(id int) ([]r.Movie, error) {
	s.repo.FindMovieByID(id)
	return []r.Movie{}, nil
}

func (s *MovieService) GetAllMoviesWithGenre(genreID int) ([]r.Movie, error) {
	s.repo.FindMoviesByGenre(genreID)
	return []r.Movie{}, nil
}

func (s *MovieService) GetAllMoviesWithYear(year int) ([]r.Movie, error) {
	s.repo.FindMoviesByYear(year)
	return []r.Movie{}, nil
}

func (s *MovieService) GetAllMoviesWithActor(actorID int) ([]r.Movie, error) {
	s.repo.FindMoviesWithActor(actorID)
	return []r.Movie{}, nil
}

func (s *MovieService) GetActorsInMovie(movieId int) ([]r.Actor, error) {
	s.repo.FindAllActorsInMovie(movieId)
	return []r.Actor{}, nil
}

func (s *MovieService) MovieMaker(movieData r.Movie) (r.Movie, error) {
	s.repo.CreateMovie(movieData)
	return r.Movie{}, nil
}

func (s *MovieService) MoviePatcher(movieData r.Movie) (r.Movie, error) {
	fields := make(map[string]string)
	s.repo.ReplaceFieldsInMovie(movieData.Id, fields)
	return r.Movie{}, nil
}

func (s *MovieService) DeleteMovie(movieID int) (bool, error) {
	s.repo.DeleteMovieByID(movieID)
	return false, nil
}

func (s *MovieService) FindActorsInMovie(movieID int) ([]r.Actor, error) {
	s.repo.FindAllActorsInMovie(movieID)
	return []r.Actor{}, nil
}
