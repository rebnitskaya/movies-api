package repository

import "fmt"

func (r movieRepository) FindAllMovies() ([]Movie, error) {
	return []Movie{}, nil
}

func (r movieRepository) CreateMovie(movieData Movie) (Movie, error) {
	query := `
		INSERT INTO movies (title, release_year, duration)
		VALUES (?,?,?)
		RETURNING id, title, release_year, duration
	`

	var movie Movie
	err := r.db.QueryRow(
		query, movieData.Title, movieData.ReleaseYear, movieData.Duration,
	).Scan(
		&movie.Id,
		&movie.Title,
		&movie.ReleaseYear,
		&movie.Duration,
	)

	if err != nil {
		return Movie{}, fmt.Errorf("Failed to create movie: %w", err)
	}

	return movie, nil
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
