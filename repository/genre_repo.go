package repository

import (
	"database/sql"
	"errors"
	"fmt"
	m "movies_api/models"
)

func (r genreRepository) FindAllGenres() ([]m.Genre, error) {
	query := `
		SELECT g.id, g.name, m.id, m.title, m.release_year, m.duration
		FROM genres g
		LEFT JOIN genres_movies mg ON g.id = mg.genre_id
		LEFT JOIN movies m ON mg.movie_id = m.id
		ORDER by g.id
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("%w: something happened during query execution: %w", m.ErrInternalIssue, err)
	}
	defer rows.Close()

	genresMap := make(map[int]*m.Genre)

	found := false
	for rows.Next() {
		found = true
		var genre m.Genre

		var movieID sql.NullInt64
		var movieTitle sql.NullString
		var releaseYear sql.NullInt64
		var duration sql.NullInt64

		err := rows.Scan(
			&genre.Id,
			&genre.Name,
			&movieID,
			&movieTitle,
			&releaseYear,
			&duration,
		)

		if err != nil {
			return nil, fmt.Errorf("%w: something happened during query execution: %w", m.ErrInternalIssue, err)
		}

		existingGenre, exists := genresMap[genre.Id]

		if !exists {
			genre.Movies = []m.Movie{}
			genresMap[genre.Id] = &genre
			existingGenre = &genre

		}

		if movieID.Valid {
			existingGenre.Movies = append(existingGenre.Movies, m.Movie{
				Id:          int(movieID.Int64),
				Title:       movieTitle.String,
				ReleaseYear: int(releaseYear.Int64),
				Duration:    int(duration.Int64),
			})
		}
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%w: something happened during query execution: %w", m.ErrInternalIssue, err)
	}

	genres := make([]m.Genre, 0, len(genresMap))

	for _, genre := range genresMap {
		genres = append(genres, *genre)
	}

	if !found {
		return nil, m.ErrGenreNotFound
	}

	return genres, nil
}

func (r genreRepository) CreateGenre(genre m.Genre) (m.Genre, error) {
	query := `
		INSERT INTO genres (name)
		VALUES (?)
		RETURNING id, name
	`
	err := r.db.QueryRow(query, genre.Name).Scan(
		&genre.Id,
		&genre.Name,
	)

	if err != nil {
		return m.Genre{}, fmt.Errorf("%w: something happened during query execution: %w", m.ErrInternalIssue, err)
	}

	return genre, nil

}

func (r genreRepository) FindGenreByID(id int) (m.Genre, error) {
	query := `
	SELECT g.id, g.name, m.id, m.title, m.release_year, m.duration
		FROM genres g
		LEFT JOIN genres_movies mg ON g.id = mg.genre_id
		LEFT JOIN movies m ON mg.movie_id = m.id
		WHERE g.id = ?
	`
	rows, err := r.db.Query(query, id)
	if err != nil {
		return m.Genre{}, fmt.Errorf("%w: something happened during query execution: %w", m.ErrInternalIssue, err)
	}
	defer rows.Close()

	var genre m.Genre
	found := false

	for rows.Next() {
		found = true
		var movieID sql.NullInt64
		var movieTitle sql.NullString
		var releaseYear sql.NullInt64
		var duration sql.NullInt64

		err := rows.Scan(
			&genre.Id,
			&genre.Name,
			&movieID,
			&movieTitle,
			&releaseYear,
			&duration,
		)
		if err != nil {
			return m.Genre{}, fmt.Errorf("%w: something happened during query execution: %w", m.ErrInternalIssue, err)
		}

		if movieID.Valid {
			genre.Movies = append(genre.Movies, m.Movie{
				Id:          int(movieID.Int64),
				Title:       movieTitle.String,
				ReleaseYear: int(releaseYear.Int64),
				Duration:    int(duration.Int64),
			})
		}
	}

	if err := rows.Err(); err != nil {
		return m.Genre{}, fmt.Errorf("%w: something happened during query execution: %w", m.ErrInternalIssue, err)
	}

	if !found {
		return m.Genre{}, m.ErrGenreNotFound
	}

	return genre, nil

}

func (r genreRepository) ReplaceFieldsInGenre(id int, name string) (m.Genre, error) {
	query := `
		UPDATE genres
		SET name = ?
		WHERE id = ?
		RETURNING id, name
	`
	var genre m.Genre

	err := r.db.QueryRow(query, name, id).Scan(
		&genre.Id,
		&genre.Name,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return m.Genre{}, m.ErrGenreNotFound
		}

		return m.Genre{}, fmt.Errorf("%w: something happened during query execution: %w", m.ErrInternalIssue, err)
	}

	return genre, nil
}

func (r genreRepository) DeleteGenreByID(id int) (bool, error) {
	query := `
		DELETE FROM genres
		WHERE id = ?
	`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return false, fmt.Errorf("%w: something happened during query execution: %w", m.ErrInternalIssue, err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("%w: something happened during query execution: %w", m.ErrInternalIssue, err)
	}
	if rowsAffected == 0 {
		return false, m.ErrGenreNotFound
	}

	return true, nil
}

func (r genreRepository) FindGenreByName(name string) (m.Genre, error) {
	query := `
		SELECT id, name
		FROM genres
		WHERE name = ?
	`

	var genre m.Genre

	err := r.db.QueryRow(query, name).Scan(
		&genre.Id,
		&genre.Name,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return m.Genre{}, m.ErrMovieNotFound
		}

		return m.Genre{}, fmt.Errorf("%w: something happened during query execution: %w", m.ErrInternalIssue, err)
	}

	return genre, nil
}

// func (r genreRepository) RemoveGenreRelationships(id int) error {
// 	query := `
// 		DELETE FROM genres_movies
// 		WHERE genre_id = ?
// 	`

// 	res, err := r.db.Exec(query, id)
// 	if err != nil {
// 		return fmt.Errorf("%w: something happened during query execution: %w", m.ErrInternalIssue, err)
// 	}

// 	rowsAffected, err := res.RowsAffected()
// 	if err != nil {
// 		return fmt.Errorf("%w: something happened during query execution: %w", m.ErrInternalIssue, err)
// 	}

// 	if rowsAffected == 0 {
// 		return m.ErrGenreNotFound
// 	}

// 	return nil
// }
