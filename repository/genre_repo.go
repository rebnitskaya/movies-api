package repository

import (
	"database/sql"
	"fmt"
	m "movies_api/models"
)

func (r genreRepository) FindAllGenres() ([]m.Genre, error) {
	return []m.Genre{}, nil
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
		return m.Genre{}, fmt.Errorf("Something happened during query execution: %w", err)
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
		return m.Genre{}, fmt.Errorf("Something happened during query execution: %w", err)
	}
	defer rows.Close()

	var genre m.Genre
	found := false

	for rows.Next() {
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
			return m.Genre{}, err
		}

		if !found {
			genre.Movies = []m.Movie{}
			found = true
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
		return m.Genre{}, err
	}

	if !found {
		return m.Genre{}, sql.ErrNoRows
	}

	return genre, nil

}

func (r genreRepository) ReplaceFieldsInGenre(id int, name string) (m.Genre, bool) {
	return m.Genre{}, false
}

func (r genreRepository) DeleteGenreByID(id int) (bool, error) {
	query := `
		DELETE FROM genres
		WHERE id = ?
	`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return false, err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}
	if rowsAffected == 0 {
		return false, sql.ErrNoRows
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
		return m.Genre{}, err
	}

	return genre, nil
}
