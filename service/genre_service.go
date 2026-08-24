package service

import (
	"database/sql"
	"errors"
	"fmt"
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

func (s *GenreService) GetGenre(id int) (m.GenreDto, error) {
	genre, err := s.repo.FindGenreByID(id)
	if err != nil {
		return m.GenreDto{}, err
	}
	genreDTO := m.GenreDto{
		Id:     genre.Id,
		Name:   genre.Name,
		Movies: []m.MovieSummaryDto{},
	}
	for _, movie := range genre.Movies {
		genreDTO.Movies = append(genreDTO.Movies, m.MovieSummaryDto{
			Id:   movie.Id,
			Name: movie.Title,
		})

	}
	return genreDTO, nil
}

func (s *GenreService) CreateGenre(name string) (m.Genre, error) {

	_, err := s.repo.FindGenreByName(name)
	if err == nil {
		return m.Genre{}, fmt.Errorf("Genre already exists")
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return m.Genre{}, err
	}

	genre := m.Genre{Name: name}
	createdGenre, err := s.repo.CreateGenre(genre)
	if err != nil {
		return m.Genre{}, err
	}

	return createdGenre, nil
}

func (s *GenreService) PatchGenre(id int, name string) (m.Genre, error) {
	s.repo.ReplaceFieldsInGenre(id, name)
	return m.Genre{}, nil
}

func (s *GenreService) DeleteGenre(id int) (bool, error) {
	if id <= 0 {
		return false, fmt.Errorf("Invalid movie id.")
	}

	deleted, err := s.repo.DeleteGenreByID(id)
	if err != nil {
		return false, err
	}
	if !deleted {
		return false, fmt.Errorf("Genre not found")
	}
	return true, nil
}
