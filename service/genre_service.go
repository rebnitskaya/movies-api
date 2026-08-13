package service

import (
	r "movies_api/repository"
)

type GenreService struct {
	repo r.GenreRepository
}

// to check
var genres = []r.Genre{
	{
		Id:   1,
		Name: "Programming horror",
	},
	{
		Id:   2,
		Name: "Funny videos",
	},
}

func (s *GenreService) GetAllGenres() ([]r.Genre, error) {
	return genres, nil
}
