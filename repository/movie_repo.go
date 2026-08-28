package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"movies_api/models"
	m "movies_api/models"
	"strings"
)

func (r movieRepository) FindAllMovies(isSearch bool, title string, limit, offset int) ([]models.MovieDto, error) {
	movies, err := r.findMovies(isSearch, title, limit, offset)
	if err != nil {
		return nil, err
	}

	movieIDs := make([]int, 0, len(movies))

	for id := range movies {
		movieIDs = append(movieIDs, id)
	}

	actors, err := r.findActorsForMovies(movieIDs)
	if err != nil {
		return nil, err
	}

	genres, err := r.findGenresForMovies(movieIDs)
	if err != nil {
		return nil, err
	}

	for id, movie := range movies {
		movie.Actors = actors[id]
		movie.Genres = genres[id]
	}

	result := make([]models.MovieDto, 0, len(movies))

	for _, movie := range movies {
		result = append(result, *movie)
	}

	return result, nil
}

func (r movieRepository) findMovies(isSearch bool, title string, limit, offset int) (map[int]*models.MovieDto, error) {
	query := `
		SELECT id, title, release_year, duration
		FROM movies
	`

	var args []any

	if isSearch {
		query += ` WHERE title LIKE ? COLLATE NOCASE`
		args = append(args, "%"+title+"%")
	}

	query += `
		ORDER BY id
		LIMIT ? OFFSET ?
	`

	args = append(args, limit, offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("%w: something happened during query execution: %w", m.ErrInternalIssue, err)
	}
	defer rows.Close()

	moviesMap := make(map[int]*models.MovieDto)

	for rows.Next() {
		var movie models.MovieDto

		err := rows.Scan(
			&movie.Id,
			&movie.Title,
			&movie.ReleaseYear,
			&movie.Duration,
		)
		if err != nil {
			return nil, fmt.Errorf("%w: something happened during query execution: %w", m.ErrInternalIssue, err)
		}

		moviesMap[movie.Id] = &movie
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%w: something happened during query execution: %w", m.ErrInternalIssue, err)
	}

	return moviesMap, nil
}

func (r movieRepository) findActorsForMovies(movieIDs []int) (map[int][]models.ActorInFilmDto, error) {
	if len(movieIDs) == 0 {
		return nil, nil
	}

	placeholders := make([]string, len(movieIDs))
	args := make([]any, len(movieIDs))

	for i, movieID := range movieIDs {
		placeholders[i] = "?"
		args[i] = movieID
	}

	query := fmt.Sprintf(`
		SELECT
			ma.movie_id,
			a.id,
			a.name,
			a.birth_date
		FROM movie_actors ma
		JOIN actors a ON a.id = ma.actor_id
		WHERE ma.movie_id IN (%s)
		ORDER BY ma.movie_id, a.id
	`, strings.Join(placeholders, ","))

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("%w: something happened during query execution: %w", m.ErrInternalIssue, err)
	}
	defer rows.Close()

	actorsMap := make(map[int][]models.ActorInFilmDto)

	for rows.Next() {
		var movieID int
		var actor models.ActorInFilmDto

		err := rows.Scan(
			&movieID,
			&actor.Id,
			&actor.Name,
			&actor.BirthDate,
		)
		if err != nil {
			return nil, fmt.Errorf("%w: something happened during query execution: %w", m.ErrInternalIssue, err)
		}

		actorsMap[movieID] = append(
			actorsMap[movieID],
			actor,
		)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%w: something happened during query execution: %w", m.ErrInternalIssue, err)
	}

	return actorsMap, nil
}

func (r movieRepository) findGenresForMovies(movieIDs []int) (map[int][]models.GenreWithoutMovies, error) {
	if len(movieIDs) == 0 {
		return make(map[int][]models.GenreWithoutMovies), nil
	}

	placeholders := make([]string, len(movieIDs))
	args := make([]any, len(movieIDs))

	for i, movieID := range movieIDs {
		placeholders[i] = "?"
		args[i] = movieID
	}

	query := fmt.Sprintf(`
		SELECT
			gm.movie_id,
			g.id,
			g.name
		FROM genres_movies gm
		JOIN genres g ON g.id = gm.genre_id
		WHERE gm.movie_id IN (%s)
		ORDER BY gm.movie_id, g.id
	`, strings.Join(placeholders, ","))

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("%w: something happened during query execution: %w", m.ErrInternalIssue, err)
	}

	genresMap := make(map[int][]models.GenreWithoutMovies)

	for rows.Next() {
		var movieID int
		var genre models.GenreWithoutMovies

		err := rows.Scan(
			&movieID,
			&genre.Id,
			&genre.Name,
		)

		if err != nil {
			return nil, fmt.Errorf("%w: something happened during query execution: %w", m.ErrInternalIssue, err)
		}

		genresMap[movieID] = append(genresMap[movieID], genre)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%w: something happened during query execution: %w", m.ErrInternalIssue, err)
	}

	return genresMap, nil
}

func (r movieRepository) CreateMovie(movieData models.CreateMovieDto) (models.Movie, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return models.Movie{}, fmt.Errorf("%w: failed to start transaction: %w", m.ErrInternalIssue, err)
	}
	defer tx.Rollback()

	//creating a movie
	query := `
		INSERT INTO movies (title, release_year, duration)
		VALUES (?,?,?)
		RETURNING id, title, release_year, duration
	`

	var movie models.Movie
	err = tx.QueryRow(
		query, movieData.Title, movieData.ReleaseYear, movieData.Duration,
	).Scan(
		&movie.Id,
		&movie.Title,
		&movie.ReleaseYear,
		&movie.Duration,
	)

	if err != nil {
		return models.Movie{}, fmt.Errorf("%w: something happened during query execution: %w", m.ErrInternalIssue, err)
	}

	//adding a connection to actors if able
	actorQuery := `
		INSERT INTO movie_actors (movie_id, actor_id)
		VALUES (?, ?)
	`

	if movieData.Actors != nil {
		for _, actorID := range movieData.Actors {
			if err := r.actorExists(tx, actorID); err != nil {
				return models.Movie{}, err
			}

			_, err = tx.Exec(actorQuery, movie.Id, actorID)
			if err != nil {
				return models.Movie{}, fmt.Errorf("%w: failed to add actor %d to movie: %w", m.ErrInternalIssue, actorID, err)
			}
		}
	}

	//adding a connection to actors if able
	genreQuery := `
		INSERT INTO genres_movies (movie_id, genre_id)
		VALUES (?, ?)
	`
	if movieData.Genres != nil {
		for _, genreID := range movieData.Genres {
			if err := r.genreExists(tx, genreID); err != nil {
				return models.Movie{}, err
			}

			_, err = tx.Exec(genreQuery, movie.Id, genreID)
			if err != nil {
				return models.Movie{}, fmt.Errorf("%w: failed to add genre %d to movie: %w", m.ErrInternalIssue, genreID, err)
			}
		}
	}

	//commiting
	if err = tx.Commit(); err != nil {
		return models.Movie{}, fmt.Errorf("%w: failed to execute movie creation: %w", m.ErrInternalIssue, err)
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
		return m.MovieDto{}, fmt.Errorf("%w: something happened during query execution: %w", m.ErrInternalIssue, err)
	}
	defer rows.Close()

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
			return m.MovieDto{}, fmt.Errorf("%w: something happened during query execution: %w", m.ErrInternalIssue, err)
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
		return m.MovieDto{}, fmt.Errorf("%w: something happened during query execution: %w", m.ErrInternalIssue, err)
	}

	if !found {
		return m.MovieDto{}, m.ErrMovieNotFound
	}

	moviesGenresMap, err := r.findGenresForMovies([]int{movie.Id})
	if err != nil {
		return m.MovieDto{}, fmt.Errorf("%w: %w", m.ErrInternalIssue, err)
	}

	movieInMap, ok := moviesGenresMap[movie.Id]
	if !ok {
		return movie, nil
	}

	movie.Genres = movieInMap

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
			return models.Movie{}, fmt.Errorf("%w: field %q cannot be changed", m.ErrBadRequest, field)
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
		if errors.Is(err, sql.ErrNoRows) {
			return m.Movie{}, m.ErrMovieNotFound
		}

		return m.Movie{}, fmt.Errorf("%w: something happened during query execution: %w", m.ErrInternalIssue, err)
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
		return false, fmt.Errorf("%w: something happened during movie deletion: %w", m.ErrInternalIssue, err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("%w: something happened during movie deletion: %w", m.ErrInternalIssue, err)
	}

	if rowsAffected == 0 {
		return false, m.ErrMovieNotFound
	}

	return true, nil
}

func (r movieRepository) FindMoviesByGenre(genreID int) ([]models.MovieDto, error) {
	query := `
		SELECT m.id, m.title, m.release_year, m.duration, g.id, g.name
		FROM movies m
		JOIN genres_movies gm ON gm.movie_id = m.id
		JOIN genres g ON g.id = gm.genre_id
		WHERE m.id IN (
		    SELECT movie_id
		    FROM genres_movies
		    WHERE genre_id = ?
		)
	`

	rows, err := r.db.Query(query, genreID)
	if err != nil {
		return []models.MovieDto{}, fmt.Errorf("%w: something happened during query execution: %w", m.ErrInternalIssue, err)
	}

	defer rows.Close()

	moviesMap := make(map[int]*models.MovieDto)
	moviesIDs := []int{}
	found := false

	for rows.Next() {
		found = true
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
			return []models.MovieDto{}, fmt.Errorf("%w: something happened during query execution: %w", m.ErrInternalIssue, err)
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
			moviesIDs = append(moviesIDs, movieID)
		}

		if genreID.Valid {
			movie.Genres = append(movie.Genres, m.GenreWithoutMovies{
				Id:   int(genreID.Int64),
				Name: genreName.String,
			})
		}
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%w: something happened during query execution: %w", m.ErrInternalIssue, err)
	}

	if !found {
		return nil, m.ErrMovieNotFound
	}

	moviesActorsMap, err := r.findActorsForMovies(moviesIDs)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", m.ErrInternalIssue, err)
	}

	movies := make([]m.MovieDto, 0, len(moviesMap))
	for i, v := range moviesMap {
		v.Actors = moviesActorsMap[i]
		movies = append(movies, *v)
	}

	return movies, nil
}

func (r movieRepository) FindMoviesByYear(year int) ([]models.MovieDto, error) {
	query := `
		SELECT m.id, m.title, m.release_year, m.duration, g.id, g.name
		FROM movies m
		JOIN genres_movies gm ON gm.movie_id = m.id
		JOIN genres g ON g.id = gm.genre_id
		WHERE release_year = ?
	`

	moviesMap := make(map[int]*models.MovieDto)
	movieIDs := []int{}
	rows, err := r.db.Query(query, year)
	if err != nil {
		return nil, fmt.Errorf("%w: something happened during query execution: %w", m.ErrInternalIssue, err)
	}
	defer rows.Close()

	found := false
	for rows.Next() {
		found = true
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
			return nil, fmt.Errorf("%w: something happened during query execution: %w", m.ErrInternalIssue, err)
		}

		movie, exists := moviesMap[movieID]
		movieIDs = append(movieIDs, movieID)

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
			movie.Genres = append(movie.Genres, m.GenreWithoutMovies{
				Id:   int(genreID.Int64),
				Name: genreName.String,
			})
		}
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%w: something happened during query execution: %w", m.ErrInternalIssue, err)
	}

	if !found {
		return nil, m.ErrMovieNotFound
	}

	moviesActorsMap, err := r.findActorsForMovies(movieIDs)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", m.ErrInternalIssue, err)
	}

	movies := make([]m.MovieDto, 0, len(moviesMap))
	for i, v := range moviesMap {
		v.Actors = moviesActorsMap[i]
		movies = append(movies, *v)
	}

	return movies, nil
}

func (r movieRepository) FindMoviesWithActor(actorID int) ([]models.MovieDto, error) {
	query := `
		SELECT m.id, m.title, m.release_year, m.duration, a.id, a.name, a.birth_date
		FROM movies m
		LEFT JOIN movie_actors ma ON ma.movie_id = m.id
		LEFT JOIN actors a ON a.id = ma.actor_id
		WHERE m.id IN (
			  SELECT movie_id
			  FROM movie_actors
			  WHERE actor_id = ?
		)
	`

	rows, err := r.db.Query(query, actorID)
	if err != nil {
		return []models.MovieDto{}, fmt.Errorf("%w: something happened during query execution: %w", m.ErrInternalIssue, err)
	}

	defer rows.Close()

	moviesMap := make(map[int]*models.MovieDto)
	movieIDs := []int{}
	found := false

	for rows.Next() {
		found = true
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
			return []models.MovieDto{}, fmt.Errorf("%w: something happened during query execution: %w", m.ErrInternalIssue, err)
		}

		movie, exists := moviesMap[movieID]
		movieIDs = append(movieIDs, movieID)
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
		return nil, fmt.Errorf("%w: something happened during query execution: %w", m.ErrInternalIssue, err)
	}

	if !found {
		return nil, m.ErrMovieNotFound
	}

	moviesGenresMap, err := r.findGenresForMovies(movieIDs)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", m.ErrInternalIssue, err)
	}

	movies := make([]m.MovieDto, 0, len(moviesMap))
	for i, v := range moviesMap {
		v.Genres = moviesGenresMap[i]
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
		return nil, fmt.Errorf("%w: something happened during query execution: %w", m.ErrInternalIssue, err)
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
			return nil, fmt.Errorf("%w: something happened during query execution: %w", m.ErrInternalIssue, err)
		}

		actors = append(actors, actor)
	}

	if !found {
		return nil, sql.ErrNoRows
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%w: something happened during query execution: %w", m.ErrInternalIssue, err)
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
		return m.Movie{}, fmt.Errorf("%w: something happened during query execution: %w", m.ErrInternalIssue, err)
	}

	defer rows.Close()

	var movie m.Movie
	found := false
	for rows.Next() {
		found = true
		err := rows.Scan(
			&movie.Id,
			&movie.Title,
			&movie.ReleaseYear,
			&movie.Duration,
		)

		if err != nil {
			return m.Movie{}, fmt.Errorf("%w: something happened during query execution: %w", m.ErrInternalIssue, err)
		}
	}

	err = rows.Err()
	if err != nil {
		return m.Movie{}, fmt.Errorf("%w: something happened during query execution: %w", m.ErrInternalIssue, err)
	}

	if !found {
		return m.Movie{}, m.ErrMovieNotFound
	}

	return movie, nil
}

func (r movieRepository) AddActorToMovie(movieID, actorID int) error {
	query := `
		INSERT INTO movie_actors (movie_id, actor_id)
		VALUES (?, ?)
	`

	checkActor := `
		SELECT id FROM actors WHERE id = ?
	`

	checkMovie := `
		SELECT id FROM movies WHERE id = ?
	`

	var id int

	err := r.db.QueryRow(checkActor, actorID).Scan(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("%w: actor with id %d doesn't exist, can't add it to movie.", m.ErrActorNotFound, actorID)
		}

		return fmt.Errorf("%w: failed to check genre: %w", m.ErrInternalIssue, err)
	}

	err = r.db.QueryRow(checkMovie, movieID).Scan(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("%w: movie with id %d doesn't exist, can't add a genre to it.", m.ErrMovieNotFound, movieID)
		}

		return fmt.Errorf("%w: failed to check movie: %w", m.ErrInternalIssue, err)
	}

	_, err = r.db.Exec(query, movieID, actorID)
	if err != nil {
		return fmt.Errorf("%w: something happened during query execution: %w", m.ErrInternalIssue, err)
	}

	return nil
}

func (r movieRepository) AddGenreToMovie(movieID, genreID int) error {
	query := `
		INSERT INTO genres_movies (movie_id, genre_id)
		VALUES (?, ?)
	`

	checkGenre := `
		SELECT id FROM genres WHERE id = ?
	`

	checkMovie := `
		SELECT id FROM movies WHERE id = ?
	`

	var id int

	err := r.db.QueryRow(checkGenre, genreID).Scan(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("%w: genre with id %d doesn't exist, can't add it to movie.", m.ErrGenreNotFound, genreID)
		}

		return fmt.Errorf("%w: failed to check genre: %w", m.ErrInternalIssue, err)
	}

	err = r.db.QueryRow(checkMovie, movieID).Scan(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("%w: movie with id %d doesn't exist, can't add a genre to it.", m.ErrMovieNotFound, movieID)
		}

		return fmt.Errorf("%w: failed to check movie: %w", m.ErrInternalIssue, err)
	}

	_, err = r.db.Exec(query, movieID, genreID)
	if err != nil {
		return fmt.Errorf("%w: something happened during query execution: %w", m.ErrInternalIssue, err)
	}

	return nil
}

func (r movieRepository) FindGenresInMovie(MovieID int) ([]m.GenreWithoutMovies, error) {
	query := `
		SELECT g.id, g.name
	 	FROM genres g
		JOIN genres_movies gm ON g.id = gm.genre_id
		WHERE gm.movie_id = ?
	`

	rows, err := r.db.Query(query, MovieID)
	if err != nil {
		return nil, fmt.Errorf("%w: something happened during query execution: %w", m.ErrInternalIssue, err)
	}
	defer rows.Close()

	genres := []m.GenreWithoutMovies{}

	found := false
	for rows.Next() {
		found = true
		var genre m.GenreWithoutMovies

		err := rows.Scan(
			&genre.Id,
			&genre.Name,
		)
		if err != nil {
			return nil, fmt.Errorf("%w: something happened during query execution: %w", m.ErrInternalIssue, err)
		}

		genres = append(genres, genre)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%w: something happened during query execution: %w", m.ErrInternalIssue, err)
	}

	if !found {
		return nil, m.ErrGenreNotFound
	}

	return genres, nil
}

func (r movieRepository) RemoveActorFromMovie(movieID, actorID int) error {
	query := `
		DELETE FROM movie_actors
		WHERE movie_id = ? AND actor_id = ?
	`

	res, err := r.db.Exec(query, movieID, actorID)
	if err != nil {
		return fmt.Errorf("%w: something happened during query execution: %w", m.ErrInternalIssue, err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%w: something happened during movie deletion: %w", m.ErrInternalIssue, err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("%w: movie with id %d has no actor with id %d.", m.ErrBadRequest, movieID, actorID)
	}

	return nil
}

func (r movieRepository) RemoveGenreFromMovie(movieID, genreID int) error {
	query := `
		DELETE FROM genres_movies
		WHERE movie_id = ? AND genre_id = ?
	`

	res, err := r.db.Exec(query, movieID, genreID)
	if err != nil {
		return fmt.Errorf("%w: something happened during query execution: %w", m.ErrInternalIssue, err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%w: something happened during movie deletion: %w", m.ErrInternalIssue, err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("%w: movie with id %d has no genre with id %d.", m.ErrBadRequest, movieID, genreID)
	}

	return nil
}

func (r movieRepository) actorExists(tx *sql.Tx, actorID int) error {
	var id int

	err := tx.QueryRow(
		`SELECT id FROM actors WHERE id = ?`,
		actorID,
	).Scan(&id)

	if errors.Is(err, sql.ErrNoRows) {
		fmt.Println("errors.Is(err, sql.ErrNoRows")
		return fmt.Errorf("%w: actor %d doesn't exist", m.ErrActorNotFound, actorID)
	}

	if err != nil {
		return fmt.Errorf("%w: failed to check actor: %w", m.ErrInternalIssue, err)
	}

	return nil
}

func (r movieRepository) genreExists(tx *sql.Tx, genreID int) error {
	var id int

	err := tx.QueryRow(
		`SELECT id FROM genres WHERE id = ?`,
		genreID,
	).Scan(&id)

	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("%w: genre %d doesn't exist", m.ErrGenreNotFound, genreID)
	}

	if err != nil {
		return fmt.Errorf("%w: failed to check genre: %w", m.ErrInternalIssue, err)
	}

	return nil
}

func (r movieRepository) CountMovies() (int, error) {
	var count int

	query := `
		SELECT COUNT(*)
		FROM movies
	`

	err := r.db.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}
