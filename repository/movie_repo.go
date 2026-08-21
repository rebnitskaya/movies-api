package repository

import (
	"database/sql"
	"fmt"
	"movies_api/models"
	m "movies_api/models"
	"strings"
)

func (r movieRepository) FindAllMovies() ([]models.MovieDto, error) {
	query := `
		SELECT m.id, m.title, m.release_year, m.duration, a.id, a.name, a.birth_date
		FROM movies m
		LEFT JOIN movie_actors ma ON ma.movie_id = m.id
		LEFT JOIN actors a ON a.id = ma.actor_id
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return []models.MovieDto{}, err
	}

	moviesMap := make(map[int]*models.MovieDto)

	for rows.Next() {
		var movieID int
		var title string
		var releaseYear int
		var duration int

		var actorID sql.NullInt64
		var actorName sql.NullString
		var actorBirthDate sql.NullString

		err := rows.Scan(
			&movieID,
			&title,
			&releaseYear,
			&duration,
			&actorID,
			&actorName,
			&actorBirthDate,
		)

		if err != nil {
			return []models.MovieDto{}, fmt.Errorf("Something happened during query execution: %w", err)
		}

		movie, exists := moviesMap[movieID]

		if !exists {
			movie = &models.MovieDto{
				Id:          movieID,
				Title:       title,
				ReleaseYear: releaseYear,
				Duration:    duration,
			}

			moviesMap[movieID] = movie
		}

		if actorID.Valid {
			movie.Actors = append(movie.Actors, m.ActorInFilmDto{
				Id:        int(actorID.Int64),
				Name:      actorName.String,
				BirthDate: actorBirthDate.String,
			})
		}
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Error while iterating movies: %w", err)
	}

	movies := make([]m.MovieDto, 0, len(moviesMap))
	for _, v := range moviesMap {
		movies = append(movies, *v)
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

func (r movieRepository) FindMovieByID(movieID int) (models.MovieDto, error) {
	query := `
		SELECT m.id, m.title, m.release_year, m.duration, a.id, a.name, a.birth_date
		FROM movies m
		LEFT JOIN movie_actors ma ON ma.movie_id = m.id
		LEFT JOIN actors a ON a.id = ma.actor_id
		WHERE m.id = ?
	`

	var movie m.MovieDto
	found := false

	rows, err := r.db.Query(query, movieID)
	if err != nil {
		return m.MovieDto{}, fmt.Errorf("Failed to find movie: %w", err)
	}

	for rows.Next() {
		found = true

		//in case if there is no actor in film
		var actorID sql.NullInt64
		var actorName sql.NullString
		var actorBirthDate sql.NullString

		err := rows.Scan(
			&movie.Id,
			&movie.Title,
			&movie.ReleaseYear,
			&movie.Duration,
			&actorID,
			&actorName,
			&actorBirthDate,
		)
		if err != nil {
			return m.MovieDto{}, fmt.Errorf("Something happened during query execution: %w", err)
		}

		if actorID.Valid {
			movie.Actors = append(movie.Actors, m.ActorInFilmDto{
				Id:        int(actorID.Int64),
				Name:      actorName.String,
				BirthDate: actorBirthDate.String,
			})
		}
	}

	if err := rows.Err(); err != nil {
		return m.MovieDto{}, fmt.Errorf("Something happened during query execution: %w", err)
	}

	if !found {
		return m.MovieDto{}, sql.ErrNoRows
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

		//since map can return stuf in random order
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

func (r movieRepository) DeleteMovieByID(movieID int) (bool, error) {
	query := `
		DELETE FROM movies
		WHERE id = ?
	`

	res, err := r.db.Exec(query, movieID)
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

func (r movieRepository) FindMoviesByGenre(genreID int) ([]models.MovieDto, error) {
	query := `
		SELECT m.id, m.title, m.release_year, m.duration, g.id, g.name
		FROM movies m
		LEFT JOIN genres_movies ma ON gm.movie_id = m.id
		LEFT JOIN genres a ON g.id = gm.genre_id
		WHERE g.id = ?
	`

	rows, err := r.db.Query(query, genreID)
	if err != nil {
		return []models.MovieDto{}, err
	}

	defer rows.Close()

	moviesMap := make(map[int]*models.MovieDto)

	for rows.Next() {
		var movieID int
		var title string
		var releaseYear int
		var duration int

		var genreID sql.NullInt64
		var genreName sql.NullString

		err := rows.Scan(
			&movieID,
			&title,
			&releaseYear,
			&duration,
			&genreID,
			&genreName,
		)

		if err != nil {
			return []models.MovieDto{}, fmt.Errorf("Something happened during query execution: %w", err)
		}

		movie, exists := moviesMap[movieID]

		if !exists {
			movie = &models.MovieDto{
				Id:          movieID,
				Title:       title,
				ReleaseYear: releaseYear,
				Duration:    duration,
			}

			moviesMap[movieID] = movie
		}

		if genreID.Valid {
			movie.Actors = append(movie.Actors, m.ActorInFilmDto{
				Id:   int(genreID.Int64),
				Name: genreName.String,
			})
		}
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Error while iterating movies: %w", err)
	}

	movies := make([]m.MovieDto, 0, len(moviesMap))
	for _, v := range moviesMap {
		movies = append(movies, *v)
	}

	return movies, nil
}

func (r movieRepository) FindMoviesByYear(year int) ([]models.Movie, error) {
	query := `
		SELECT id, title, release_year, duration
		FROM movies
		WHERE release_year = ?
	`

	movies := []m.Movie{}
	rows, err := r.db.Query(query, year)
	if err != nil {
		return movies, err
	}
	defer rows.Close()

	for rows.Next() {
		movie := models.Movie{}
		err := rows.Scan(
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

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Error while iterating movies: %w", err)
	}

	return movies, nil
}

func (r movieRepository) FindMoviesWithActor(actorID int) ([]models.MovieDto, error) {
	query := `
		SELECT m.id, m.title, m.release_year, m.duration, a.id, a.name, a.birth_date
		FROM movies m
		LEFT JOIN movie_actors ma ON ma.movie_id = m.id
		LEFT JOIN actors a ON a.id = ma.actor_id
		WHERE a.id = ?
	`

	rows, err := r.db.Query(query, actorID)
	if err != nil {
		return []models.MovieDto{}, err
	}

	defer rows.Close()

	moviesMap := make(map[int]*models.MovieDto)

	for rows.Next() {
		var movieID int
		var title string
		var releaseYear int
		var duration int

		var actorID sql.NullInt64
		var actorName sql.NullString
		var actorBirthDate sql.NullString

		err := rows.Scan(
			&movieID,
			&title,
			&releaseYear,
			&duration,
			&actorID,
			&actorName,
			&actorBirthDate,
		)

		if err != nil {
			return []models.MovieDto{}, fmt.Errorf("Something happened during query execution: %w", err)
		}

		movie, exists := moviesMap[movieID]

		if !exists {
			movie = &models.MovieDto{
				Id:          movieID,
				Title:       title,
				ReleaseYear: releaseYear,
				Duration:    duration,
			}

			moviesMap[movieID] = movie
		}

		if actorID.Valid {
			movie.Actors = append(movie.Actors, m.ActorInFilmDto{
				Id:        int(actorID.Int64),
				Name:      actorName.String,
				BirthDate: actorBirthDate.String,
			})
		}
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Error while iterating movies: %w", err)
	}

	movies := make([]m.MovieDto, 0, len(moviesMap))
	for _, v := range moviesMap {
		movies = append(movies, *v)
	}

	return movies, nil
}

func (r movieRepository) FindAllActorsInMovie(movieID int) ([]models.ActorInFilmDto, error) {
	query := `
		SELECT a.id, a.name, a.birth_date
		FROM movie_actors ma
		JOIN actors a ON a.id = ma.actor_id
		WHERE ma.movie_id = ?
	`

	rows, err := r.db.Query(query, movieID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var actors []models.ActorInFilmDto
	found := false
	for rows.Next() {
		var actor m.ActorInFilmDto
		found = true
		err := rows.Scan(
			&actor.Id,
			&actor.Name,
			&actor.BirthDate,
		)

		if err != nil {
			return nil, fmt.Errorf("Something happened during query execution: %w", err)
		}

		actors = append(actors, actor)
	}

	if !found {
		return nil, sql.ErrNoRows
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Error while iterating movies: %w", err)
	}

	return actors, nil
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

func (r movieRepository) AddActorToMovie(movieID, actorID int) error {
	query := `
		INSERT INTO movie_actors (movie_id, actor_id)
		VALUES (?, ?)
	`

	_, err := r.db.Exec(query, movieID, actorID)
	if err != nil {
		return fmt.Errorf("Failed to add actor to movie: %w", err)
	}

	return nil
}

func (r movieRepository) RemoveActorFromMovie(movieID, actorID int) error {
	query := `
		DELETE FROM movie_actors
		WHERE movie_id = ? AND actor_id = ?
	`

	_, err := r.db.Exec(query, movieID, actorID)
	if err != nil {
		return fmt.Errorf("Failed to remove actor from movie: %w", err)
	}

	return nil
}
