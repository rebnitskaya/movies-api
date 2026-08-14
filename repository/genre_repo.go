package repository

func (r genreRepository) FindAllGenres() ([]Genre, error) {
	return []Genre{}, nil
}

func (r genreRepository) CreateGenre(genre Genre) (bool, error) {
	return false, nil
}

func (r genreRepository) FindGenreByID(id int) (Genre, bool) {
	return Genre{}, false
}

func (r genreRepository) ReplaceFieldsInGenre(id int, name string) (Genre, bool) {
	return Genre{}, false
}

func (r genreRepository) DeleteGenreByID(id int) (bool, error) {
	return false, nil
}
