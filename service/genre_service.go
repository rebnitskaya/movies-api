package service

import (
	m "movies_api/models"
	r "movies_api/repository"
)

type GenreService struct {
	repo r.GenreRepository
}

func (s *GenreService) GetAllGenres() ([]m.Genre, error) {
	s.repo.FindAllGenres()
	return []m.Genre{}, nil
}

func (s *GenreService) GetGenre(id int) (m.Genre, error) {
	s.repo.FindGenreByID(id)
	return m.Genre{}, nil
}

func (s *GenreService) CreateGenre(name string) (m.Genre, error) {
	genre := m.Genre{}
	s.repo.CreateGenre(genre)
	return m.Genre{}, nil
}

func (s *GenreService) PatchGenre(id int, name string) (m.Genre, error) {
	s.repo.ReplaceFieldsInGenre(id, name)
	return m.Genre{}, nil
}

func (s *GenreService) DeleteGenre(id int) (m.Genre, error) {
	s.repo.DeleteGenreByID(id)
	return m.Genre{}, nil
}
