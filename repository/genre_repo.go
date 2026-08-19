package repository

import "movies_api/models"

func (r genreRepository) FindAllGenres() ([]models.Genre, error) {
	return []models.Genre{}, nil
}

func (r genreRepository) CreateGenre(genre models.Genre) (bool, error) {
	return false, nil
}

func (r genreRepository) FindGenreByID(id int) (models.Genre, bool) {
	return models.Genre{}, false
}

func (r genreRepository) ReplaceFieldsInGenre(id int, name string) (models.Genre, bool) {
	return models.Genre{}, false
}

func (r genreRepository) DeleteGenreByID(id int) (bool, error) {
	return false, nil
}
