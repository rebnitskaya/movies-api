package service

import (
	"context"
	"fmt"
	m "movies_api/models"
	r "movies_api/repository"
	"slices"
	"strings"
)

type MovieService struct {
	repo r.MovieRepository
}

func (s *MovieService) GetAllMovies(page, limit int, ctx context.Context) (m.MoviesPaginated, error) {
	offset := page * limit
	movies, err := s.repo.FindAllMovies(false, "", limit, offset, ctx)
	if err != nil {
		return m.MoviesPaginated{}, err
	}

	sortMovies(movies)

	pageSize := len(movies)
	totalMovies, err := s.repo.CountMovies(ctx)
	if err != nil {
		return m.MoviesPaginated{}, err
	}

	totalPages := (totalMovies + limit - 1) / limit
	if page > totalPages && totalPages > 0 {
		return m.MoviesPaginated{
			Movies: []m.MovieDto{},
		}, nil
	}

	return m.MoviesPaginated{
		Movies:     movies,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil

}

func (s *MovieService) FindOneMovie(id int, ctx context.Context) (m.MovieDto, error) {
	if id <= 0 {
		return m.MovieDto{}, fmt.Errorf("%w: invalid id", m.ErrInvalidInput)
	}

	movie, err := s.repo.FindMovieByID(id, ctx)
	if err != nil {
		return m.MovieDto{}, err
	}

	return movie, nil
}

func (s *MovieService) GetAllMoviesWithGenre(genreID, page, limit int, ctx context.Context) (m.MoviesPaginated, error) {
	if genreID <= 0 {
		return m.MoviesPaginated{}, fmt.Errorf("%w: invalid id", m.ErrInvalidInput)
	}

	offset := page * limit
	movies, err := s.repo.FindMoviesByGenre(genreID, limit, offset, ctx)
	if err != nil {
		return m.MoviesPaginated{}, err
	}

	sortMovies(movies)

	totalMovies, err := s.repo.CountMoviesByGenre(genreID, ctx)

	if err != nil {
		return m.MoviesPaginated{}, err
	}
	totalPages := (totalMovies + limit - 1) / limit
	if page > totalPages && totalPages > 0 {
		return m.MoviesPaginated{
			Movies: []m.MovieDto{},
		}, nil
	}
	return m.MoviesPaginated{
		Movies:     movies,
		Page:       page,
		PageSize:   len(movies),
		TotalPages: totalPages,
	}, nil

}

func (s *MovieService) GetAllMoviesWithYear(year, page, limit int, ctx context.Context) (m.MoviesPaginated, error) {
	if year < 1885 || year > 2050 || year == 0 {
		return m.MoviesPaginated{}, fmt.Errorf("%w: invalid release year.", m.ErrInvalidInput)
	}

	offset := page * limit
	movies, err := s.repo.FindMoviesByYear(year, limit, offset, ctx)
	if err != nil {
		return m.MoviesPaginated{}, err
	}

	sortMovies(movies)
	totalMovies, err := s.repo.CountMoviesByYear(year, ctx)

	if err != nil {
		return m.MoviesPaginated{}, err
	}
	totalPages := (totalMovies + limit - 1) / limit
	if page > totalPages && totalPages > 0 {
		return m.MoviesPaginated{
			Movies: []m.MovieDto{},
		}, nil
	}
	return m.MoviesPaginated{
		Movies:     movies,
		Page:       page,
		PageSize:   len(movies),
		TotalPages: totalPages,
	}, nil

}

func (s *MovieService) GetAllMoviesWithActor(actorID, page, limit int, ctx context.Context) (m.MoviesPaginated, error) {
	if actorID <= 0 {
		return m.MoviesPaginated{}, fmt.Errorf("%w: invalid id", m.ErrInvalidInput)
	}

	offset := page * limit
	movies, err := s.repo.FindMoviesWithActor(actorID, limit, offset, ctx)
	if err != nil {
		return m.MoviesPaginated{}, err
	}

	sortMovies(movies)

	totalMovies, err := s.repo.CountMoviesByActor(actorID, ctx)

	if err != nil {
		return m.MoviesPaginated{}, err
	}
	totalPages := (totalMovies + limit - 1) / limit
	if page > totalPages && totalPages > 0 {
		return m.MoviesPaginated{
			Movies: []m.MovieDto{},
		}, nil
	}
	return m.MoviesPaginated{
		Movies:     movies,
		Page:       page,
		PageSize:   len(movies),
		TotalPages: totalPages,
	}, nil

}

func (s *MovieService) GetAllMoviesWithTitle(title string, page, limit int, ctx context.Context) ([]m.MovieDto, error) {
	title = strings.ToLower(title)

	offset := page * limit

	movies, err := s.repo.FindAllMovies(true, title, limit, offset, ctx)
	if err != nil {
		return []m.MovieDto{}, err
	}

	sortMovies(movies)

	return movies, nil
}

func (s *MovieService) MovieMaker(movieData m.CreateMovieDto, ctx context.Context) (m.MovieDto, error) {
	_, err := movieData.Validate()
	if err != nil {
		return m.MovieDto{}, err
	}

	res, err := s.repo.FindMovieByTitleAndYear(movieData.Title, movieData.ReleaseYear, ctx)
	if res.Title == movieData.Title && res.ReleaseYear == movieData.ReleaseYear {
		return m.MovieDto{}, m.ErrMovieHasBeenMadeBefore
	}

	movie, error := s.repo.CreateMovie(movieData, ctx)
	if error != nil {
		return m.MovieDto{}, error
	}

	movieDto, error := s.repo.FindMovieByID(movie.Id, ctx)
	if error != nil {
		return m.MovieDto{}, error
	}

	return movieDto, nil
}

func (s *MovieService) MoviePatcher(movieData m.MoviePatchDto, movieID int, ctx context.Context) (m.Movie, error) {
	if movieID <= 0 {
		return m.Movie{}, fmt.Errorf("%w: invalid movie id", m.ErrInvalidInput)
	}

	fields := make(map[string]any)
	_, err := movieData.Validate()
	if err != nil {
		return m.Movie{}, err
	}

	if movieData.Title != nil {
		fields["title"] = *movieData.Title
	}
	if movieData.Duration != nil {
		fields["duration"] = *movieData.Duration
	}
	if movieData.ReleaseYear != nil {
		fields["release_year"] = *movieData.ReleaseYear
	}
	if len(fields) == 0 {
		return m.Movie{}, fmt.Errorf("%w: no fields to update.", m.ErrBadRequest)
	}

	movie, err := s.repo.ReplaceFieldsInMovie(movieID, fields, ctx)
	if err != nil {
		return m.Movie{}, err
	}

	return movie, nil
}

func (s *MovieService) DeleteMovie(movieID int, force bool, ctx context.Context) (bool, error) {
	if movieID <= 0 {
		return false, fmt.Errorf("%w: invalid movie id.", m.ErrBadRequest)
	}

	if !force {
		return false, fmt.Errorf("%w: cannot delete movie, use \"force\" to perform deletion.", m.ErrBadRequest)
	}

	ok, err := s.repo.DeleteMovieByID(movieID, ctx)
	if err != nil {
		return false, err
	}

	return ok, nil
}

func (s *MovieService) FindActorsInMovie(movieID int, ctx context.Context) ([]m.ActorInFilmDto, error) {
	if movieID <= 0 {
		return []m.ActorInFilmDto{}, fmt.Errorf("%w: invalid movie id.", m.ErrBadRequest)
	}

	movies, err := s.repo.FindAllActorsInMovie(movieID, ctx)
	if err != nil {
		return []m.ActorInFilmDto{}, err
	}

	return movies, nil
}

func (s *MovieService) AddActorToMovie(movieID, actorID int, ctx context.Context) (m.MovieDto, error) {
	if movieID <= 0 {
		return m.MovieDto{}, fmt.Errorf("%w: invalid movie id.", m.ErrBadRequest)
	}

	if actorID <= 0 {
		return m.MovieDto{}, fmt.Errorf("%w: invalid actor id.", m.ErrBadRequest)
	}

	err := s.repo.AddActorToMovie(movieID, actorID, ctx)
	if err != nil {
		return m.MovieDto{}, err
	}

	movie, err := s.repo.FindMovieByID(movieID, ctx)
	if err != nil {
		return m.MovieDto{}, err
	}

	return movie, nil
}

func (s *MovieService) DeleteActorFromMovie(movieID, actorID int, ctx context.Context) (m.MovieDto, error) {
	if movieID <= 0 {
		return m.MovieDto{}, fmt.Errorf("%w: invalid movie id.", m.ErrBadRequest)
	}

	if actorID <= 0 {
		return m.MovieDto{}, fmt.Errorf("%w: invalid actor id.", m.ErrBadRequest)
	}

	movie, err := s.repo.FindMovieByID(movieID, ctx)
	if err != nil {
		return m.MovieDto{}, err
	}

	err = s.repo.RemoveActorFromMovie(movieID, actorID, ctx)
	if err != nil {
		return m.MovieDto{}, err
	}

	return movie, nil
}

func sortMovies(movies []m.MovieDto) {
	slices.SortFunc(movies, func(a, b m.MovieDto) int {
		if a.Id > b.Id {
			return 1
		} else {
			return -1
		}
	})
}
