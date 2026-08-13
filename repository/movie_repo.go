package repository

func (r movieRepository) FindAllMovies() ([]Movie, error) {
	return []Movie{}, nil
}
func (r movieRepository) CreateMovie(movieData Movie) (bool, error) {
	return false, nil
}
func (r movieRepository) FindMovieByID(movieId int) (Movie, bool) {
	return Movie{}, false
}
func (r movieRepository) ReplaceFieldsInMovie(movieId int, filedsToUpdate map[string]string) (Movie, bool) {
	return Movie{}, false
}
func (r movieRepository) DeleteMovieByID(movieId int) (bool, error) {
	return false, nil
}
func (r movieRepository) FindMoviesByGenre(genreId int) ([]Movie, error) {
	return []Movie{}, nil
}
func (r movieRepository) FindMoviesByYear(year int) ([]Movie, error) {
	return []Movie{}, nil
}
func (r movieRepository) FindMoviesWithActor(actorId int) ([]Movie, error) {
	return []Movie{}, nil
}
func (r movieRepository) FindAllActorsInMovie(movieId int) ([]Actor, error) {
	return []Actor{}, nil
}
