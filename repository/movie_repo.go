package repository

import (
	"fmt"
	"movies_api/models"
)

func (r movieRepository) FindAllMovies() ([]models.Movie, error) {
	return []models.Movie{}, nil
}

func (r movieRepository) CreateMovie(movieData models.MovieDto) (models.Movie, error) {
	query := `
		INSERT INTO movies (title, release_year, duration)
		VALUES (?,?,?)
		RETURNING id, title, release_year, duration
	`

	var movie models.Movie
	err := r.db.QueryRow(
		query, movieData.Title, movieData.ReleaseYear, movieData.Duration,
	).Scan(
		&movie.Id,
		&movie.Title,
		&movie.ReleaseYear,
		&movie.Duration,
	)

	if err != nil {
		return models.Movie{}, fmt.Errorf("Failed to create movie: %w", err)
	}

	return movie, nil
}

func (r movieRepository) FindMovieByID(movieId int) (models.Movie, bool) {
	return models.Movie{}, false
}
func (r movieRepository) ReplaceFieldsInMovie(movieId int, filedsToUpdate map[string]string) (models.Movie, bool) {
	return models.Movie{}, false
}
func (r movieRepository) DeleteMovieByID(movieId int) (bool, error) {
	return false, nil
}
func (r movieRepository) FindMoviesByGenre(genreId int) ([]models.Movie, error) {
	return []models.Movie{}, nil
}
func (r movieRepository) FindMoviesByYear(year int) ([]models.Movie, error) {
	return []models.Movie{}, nil
}
func (r movieRepository) FindMoviesWithActor(actorId int) ([]models.Movie, error) {
	return []models.Movie{}, nil
}
func (r movieRepository) FindAllActorsInMovie(movieId int) ([]models.Actor, error) {
	return []models.Actor{}, nil
}
