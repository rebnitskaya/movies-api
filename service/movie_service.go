package service

import (
	"movies_api/models"
	"movies_api/repository"
)

type MovieService struct {
	repo repository.MovieRepository
}

// to check
var movies = []models.Movie{
	{
		Id:          1,
		Title:       "Titanic",
		ReleaseYear: 1997,
		Duration:    194,
		Genres:      []models.Genre{},
		Actors:      []models.Actor{},
	},
	{
		Id:          2,
		Title:       "Interstellar",
		ReleaseYear: 2014,
		Duration:    169,
		Genres:      []models.Genre{},
		Actors:      []models.Actor{},
	},
}

func (s *MovieService) GetAllMovies() ([]models.Movie, error) {
	return movies, nil
}
