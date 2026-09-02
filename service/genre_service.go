package service

import (
	"context"
	"errors"
	"fmt"
	m "movies_api/models"
	r "movies_api/repository"
	"sort"
)

type GenreService struct {
	repo r.GenreRepository
}

func (s *GenreService) GetAllGenres(page, limit int, ctx context.Context) (m.GenresPaginated, error) {
	offset := page * limit

	genres, err := s.repo.FindAllGenres(limit, offset, ctx)
	if err != nil {
		return m.GenresPaginated{}, err
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

	pageSize := len(result)
	totalGenres, err := s.repo.CountGenres(ctx)
	if err != nil {
		return m.GenresPaginated{}, err
	}

	totalPages := (totalGenres + limit - 1) / limit
	if page > totalPages && totalPages > 0 {
		return m.GenresPaginated{
			Genres: []m.GenreDto{},
		}, nil
	}

	return m.GenresPaginated{
		Genres:     result,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

func (s *GenreService) GetGenre(id int, ctx context.Context) (m.GenreDto, error) {
	if id <= 0 {
		return m.GenreDto{}, fmt.Errorf("%w: invalid genre id.", m.ErrBadRequest)
	}

	genre, err := s.repo.FindGenreByID(id, ctx)
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

func (s *GenreService) CreateGenre(name string, ctx context.Context) (m.Genre, error) {
	if name == "" {
		return m.Genre{}, fmt.Errorf("%w: genre name cannot be empty", m.ErrBadRequest)
	}

	_, err := s.repo.FindGenreByName(name, ctx)
	if err == nil {
		return m.Genre{}, fmt.Errorf("%w: genre already exists", m.ErrBadRequest)
	}

	genre := m.Genre{Name: name}
	createdGenre, err := s.repo.CreateGenre(genre, ctx)
	if err != nil {
		return m.Genre{}, err
	}

	return createdGenre, nil
}

func (s *GenreService) PatchGenre(id int, name string, ctx context.Context) (m.Genre, error) {
	if name == "" {
		return m.Genre{}, fmt.Errorf("%w: genre name cannot be empty", m.ErrBadRequest)
	}

	if id <= 0 {
		return m.Genre{}, fmt.Errorf("%w: invalid genre id.", m.ErrBadRequest)
	}

	genre, err := s.repo.ReplaceFieldsInGenre(id, name, ctx)
	if err != nil {
		return m.Genre{}, err
	}
	return genre, nil
}

func (s *GenreService) DeleteGenre(id int, force bool, ctx context.Context) (bool, error) {
	if id <= 0 {
		return false, fmt.Errorf("%w: invalid genre id.", m.ErrBadRequest)
	}

	genre, err := s.repo.FindGenreByID(id, ctx)
	if err != nil {
		return false, err
	}

	if !force && len(genre.Movies) > 0 {
		return false, fmt.Errorf("%w: cannot delete genre '%s' because it has %d associated movies",
			m.ErrBadRequest,
			genre.Name,
			len(genre.Movies),
		)
	}

	_, err = s.repo.DeleteGenreByID(id, ctx)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (s *MovieService) AddGenreToMovie(movieID, genreID int, ctx context.Context) (m.MovieDto, error) {
	if movieID <= 0 {
		return m.MovieDto{}, fmt.Errorf("%w: invalid movie id.", m.ErrBadRequest)
	}

	if genreID <= 0 {
		return m.MovieDto{}, fmt.Errorf("%w: invalid genre id.", m.ErrBadRequest)
	}

	err := s.repo.AddGenreToMovie(movieID, genreID, ctx)
	if err != nil {
		return m.MovieDto{}, err
	}

	movie, err := s.repo.FindMovieByID(movieID, ctx)
	if err != nil {
		return m.MovieDto{}, err
	}

	genres, err := s.repo.FindGenresInMovie(movieID, ctx)
	if err != nil {
		return m.MovieDto{}, err
	}

	movie.Genres = genres

	return movie, nil
}

func (s *MovieService) DeleteGenreFromMovie(movieID, genreID int, ctx context.Context) (m.MovieDto, error) {
	if movieID <= 0 {
		return m.MovieDto{}, fmt.Errorf("%w invalid movie id.", m.ErrBadRequest)
	}

	if genreID <= 0 {
		return m.MovieDto{}, fmt.Errorf("%w: invalid genre id.", m.ErrBadRequest)
	}

	movie, err := s.repo.FindMovieByID(movieID, ctx)
	if err != nil {
		return m.MovieDto{}, err
	}

	err = s.repo.RemoveGenreFromMovie(movieID, genreID, ctx)
	if err != nil {
		return m.MovieDto{}, err
	}

	genres, err := s.repo.FindGenresInMovie(movieID, ctx)
	if err != nil {
		if errors.Is(err, m.ErrGenreNotFound) {
			movie.Genres = genres
		} else {
			return m.MovieDto{}, err
		}
	}

	return movie, nil
}
