package service

import "movies_api/repository"

func NewActorService(repo repository.ActorRepository) *ActorService {
	return &ActorService{
		repo: repo,
	}
}

func NewMovieService(repo repository.MovieRepository) *MovieService {
	return &MovieService{
		repo: repo,
	}
}

func NewGenreService(repo repository.GenreRepository) *GenreService {
	return &GenreService{
		repo: repo,
	}
}
