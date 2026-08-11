package service

import (
	"movies_api/models"
	"movies_api/repository"
)

type GenreService struct {
	repo repository.GenreRepository
}

// to check
var genres = []models.Genre{
	{
		Id:   1,
		Name: "Programming horror",
	},
	{
		Id:   2,
		Name: "Funny videos",
	},
}

func (s *GenreService) GetAllGenres() ([]models.Genre, error) {
	return genres, nil
}

func NewGenreService(repo repository.GenreRepository) *GenreService {
	return &GenreService{
		repo: repo,
	}
}
