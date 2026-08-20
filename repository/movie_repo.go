package repository

import (
	"database/sql"
	"fmt"
	"movies_api/models"
	m "movies_api/models"
	"strings"
)

func (r movieRepository) FindAllMovies() ([]models.Movie, error) {
	query := `
		SELECT id, title, release_year, duration
		FROM movies
	`

	res, err := r.db.Query(query)
	if err != nil {
		return []models.Movie{}, err
	}

	defer res.Close()

	movies := []models.Movie{}
	for res.Next() {
		movie := models.Movie{}

		err := res.Scan(
			&movie.Id,
			&movie.Title,
			&movie.ReleaseYear,
			&movie.Duration,
		)

		if err != nil {
			return []models.Movie{}, fmt.Errorf("Something happened during query execution: %w", err)
		}

		movies = append(movies, movie)
	}

	if err := res.Err(); err != nil {
		return nil, fmt.Errorf("Error while iterating movies: %w", err)
	}

	return movies, nil
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

func (r movieRepository) FindMovieByID(movieId int) (models.Movie, error) {
	query := `
		SELECT id, title, release_year, duration
		FROM movies
		WHERE id = ?
	`

	var movie m.Movie

	err := r.db.QueryRow(query, movieId).Scan(
		&movie.Id,
		&movie.Title,
		&movie.ReleaseYear,
		&movie.Duration,
	)

	if err != nil {
		return m.Movie{}, fmt.Errorf("Something happened during query execution: %w", err)
	}

	return movie, nil
}

func (r movieRepository) ReplaceFieldsInMovie(movieID int, filedsToUpdate map[string]any) (models.Movie, error) {
	alowedToChange := map[string]bool{
		"title":        true,
		"release_year": true,
		"duration":     true,
	}

	setFields := make([]string, 0, len(filedsToUpdate))
	args := make([]any, 0, len(filedsToUpdate))

	for field, value := range filedsToUpdate {
		if !alowedToChange[field] {
			return m.Movie{}, fmt.Errorf("Not allowed to change the field: %s", field)
		}

		setFields = append(setFields, field+" = ?")
		args = append(args, value)
	}

	query := `
		UPDATE movies
		SET ` + strings.Join(setFields, ", ") + `
		WHERE id = ?
		RETURNING id, title, release_year, duration
	`

	args = append(args, movieID)

	var movie m.Movie
	err := r.db.QueryRow(query, args...).Scan(
		&movie.Id,
		&movie.Title,
		&movie.ReleaseYear,
		&movie.Duration,
	)

	if err != nil {
		return m.Movie{}, fmt.Errorf("Failed to update movie: %w", err)
	}

	return movie, nil
}

func (r movieRepository) DeleteMovieByID(movieId int) (bool, error) {
	query := `
		DELETE FROM movies
		WHERE id = ?
	`

	res, err := r.db.Exec(query, movieId)
	if err != nil {
		return false, fmt.Errorf("Failed to delete movie: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("Failed to check deleted movie: %w", err)
	}

	if rowsAffected == 0 {
		return false, sql.ErrNoRows
	}

	return true, nil
}

func (r movieRepository) FindMoviesByGenre(genreId int) ([]models.Movie, error) {
	return []models.Movie{}, nil
}

func (r movieRepository) FindMoviesByYear(year int) ([]models.Movie, error) {
	query := `
		SELECT id, title, release_year, duration
		FROM movies
		WHERE release_year = ?
	`

	movies := []m.Movie{}
	res, err := r.db.Query(query, year)
	if err != nil {
		return movies, err
	}
	defer res.Close()

	for res.Next() {
		movie := models.Movie{}
		err := res.Scan(
			&movie.Id,
			&movie.Title,
			&movie.ReleaseYear,
			&movie.Duration,
		)

		if err != nil {
			return []models.Movie{}, fmt.Errorf("Something happened during query execution: %w", err)
		}

		movies = append(movies, movie)
	}

	if err := res.Err(); err != nil {
		return nil, fmt.Errorf("Error while iterating movies: %w", err)
	}

	return movies, nil
}

func (r movieRepository) FindMoviesWithActor(actorId int) ([]models.Movie, error) {
	return []models.Movie{}, nil
}
func (r movieRepository) FindAllActorsInMovie(movieId int) ([]models.Actor, error) {
	return []models.Actor{}, nil
}

func (r movieRepository) FindMovieByTitleAndYear(title string, year int) (models.Movie, error) {
	query := `
		SELECT id, title, release_year, duration
		FROM movies
		WHERE title = ? AND release_year = ?
	`
	rows, err := r.db.Query(query, title, year)
	if err != nil {
		return m.Movie{}, fmt.Errorf("Something happened during query execution: %w", err)
	}

	defer rows.Close()

	var movie m.Movie

	for rows.Next() {
		err := rows.Scan(
			&movie.Id,
			&movie.Title,
			&movie.ReleaseYear,
			&movie.Duration,
		)

		if err != nil {
			return m.Movie{}, fmt.Errorf("Something happened during query execution: %w", err)
		}
	}

	err = rows.Err()
	if err != nil {
		return m.Movie{}, fmt.Errorf("Something happened during query execution: %w", err)
	}

	return movie, nil
}
