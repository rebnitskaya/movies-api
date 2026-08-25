package service

import (
	"database/sql"
	"errors"
	"fmt"
	m "movies_api/models"
	r "movies_api/repository"
	"sort"
)

type GenreService struct {
	repo r.GenreRepository
}

func (s *GenreService) GetAllGenres() ([]m.GenreDto, error) {
	genres, err := s.repo.FindAllGenres()
	if err != nil {
		return nil, err
	}
	sort.Slice(genres, func(i, j int) bool {
		return genres[i].Id < genres[j].Id
	})

	result := make([]m.GenreDto, 0, len(genres))
	for _, genre := range genres {
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
		result = append(result, genreDTO)
	}
	return result, nil
}

func (s *GenreService) GetGenre(id int) (m.GenreDto, error) {
	if id <= 0 {
		return m.GenreDto{}, fmt.Errorf("Invalid genre id")
	}

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

	if name == "" {
		return m.Genre{}, fmt.Errorf("Genre name cannot be empty")
	}

	if name == "" {
		return m.Genre{}, fmt.Errorf("Genre name cannot be empty")
	}

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

	if name == "" {
		return m.Genre{}, fmt.Errorf("Genre name cannot be empty")
	}

	if id <= 0 {
		return m.Genre{}, fmt.Errorf("Invalid genre id")
	}

	genre, err := s.repo.ReplaceFieldsInGenre(id, name)
	if err != nil {
		return m.Genre{}, err
	}
	return genre, nil
}

func (s *GenreService) DeleteGenre(id int, force bool) (bool, error) {
	if id <= 0 {
		return false, fmt.Errorf("Invalid genre id.")
	}

	genre, err := s.repo.FindGenreByID(id)
	if err != nil {
		return false, err
	}

	if !force && len(genre.Movies) > 0 {
		return false, fmt.Errorf("Cannot delete genre '%s' because it has %d associated movies",
			genre.Name,
			len(genre.Movies),
		)
	}

	if force {
		err := s.repo.RemoveGenreRelationships(id)
		if err != nil {
			return false, err
		}
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

func (s *MovieService) AddGenreToMovie(movieID, genreID int) (m.MovieDto, error) {
	if movieID <= 0 {
		return m.MovieDto{}, fmt.Errorf("Invalid movie id.")
	}

	if genreID <= 0 {
		return m.MovieDto{}, fmt.Errorf("Invalid genre id.")
	}

	err := s.repo.AddGenreToMovie(movieID, genreID)
	if err != nil {
		return m.MovieDto{}, err
	}

	movie, err := s.repo.FindMovieByID(movieID)
	if err != nil {
		return m.MovieDto{}, err
	}

	genres, err := s.repo.FindGenresInMovie(movieID)
	if err != nil {
		return m.MovieDto{}, err
	}

	movie.Genres = genres

	return movie, nil
}

func (s *MovieService) DeleteGenreFromMovie(movieID, genreID int) (m.MovieDto, error) {
	if movieID <= 0 {
		return m.MovieDto{}, fmt.Errorf("Invalid movie id.")
	}

	if genreID <= 0 {
		return m.MovieDto{}, fmt.Errorf("Invalid genre id.")
	}

	err := s.repo.RemoveGenreFromMovie(movieID, genreID)
	if err != nil {
		return m.MovieDto{}, err
	}

	movie, err := s.repo.FindMovieByID(movieID)
	if err != nil {
		return m.MovieDto{}, err
	}

	genres, err := s.repo.FindGenresInMovie(movieID)
	if err != nil {
		return m.MovieDto{}, err
	}

	movie.Genres = genres

	return movie, nil
}
