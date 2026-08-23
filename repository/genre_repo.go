package repository

import (
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

func (r genreRepository) FindGenreByID(id int) (m.Genre, bool) {
	return m.Genre{}, false
}

func (r genreRepository) ReplaceFieldsInGenre(id int, name string) (m.Genre, bool) {
	return m.Genre{}, false
}

func (r genreRepository) DeleteGenreByID(id int) (bool, error) {
	return false, nil
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
