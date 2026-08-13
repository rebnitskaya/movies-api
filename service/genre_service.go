package service

import (
	r "movies_api/repository"
)

type GenreService struct {
	repo r.GenreRepository
}

func (s *GenreService) GetAllGenres() ([]r.Genre, error) {
	s.repo.FindAllGenres()
	return []r.Genre{}, nil
}

func (s *GenreService) GetGenre(id int) (r.Genre, error) {
	s.repo.FindGenreByID(id)
	return r.Genre{}, nil
}

func (s *GenreService) CreateGenre(name string) (r.Genre, error) {
	genre := r.Genre{}
	s.repo.CreateGenre(genre)
	return r.Genre{}, nil
}

func (s *GenreService) PatchGenre(id int, name string) (r.Genre, error) {
	s.repo.ReplaceFieldsInGenre(id, name)
	return r.Genre{}, nil
}

func (s *GenreService) DeleteGenre(id int) (r.Genre, error) {
	s.repo.DeleteGenreByID(id)
	return r.Genre{}, nil
}
