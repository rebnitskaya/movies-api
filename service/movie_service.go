package service

import (
	r "movies_api/repository"
)

type MovieService struct {
	repo r.MovieRepository
}

// to check
var movies = []r.Movie{
	{
		Id:          1,
		Title:       "Titanic",
		ReleaseYear: 1997,
		Duration:    194,
		Genres:      []r.Genre{},
		Actors:      []r.Actor{},
	},
	{
		Id:          2,
		Title:       "Interstellar",
		ReleaseYear: 2014,
		Duration:    169,
		Genres:      []r.Genre{},
		Actors:      []r.Actor{},
	},
}

func (s *MovieService) GetAllMovies() ([]r.Movie, error) {
	return movies, nil
}
