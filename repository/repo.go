package repository

import (
	"database/sql"
)

type Repository struct {
	MovieRepo MovieRepository
	ActorRepo ActorRepository
	GenreRepo GenreRepository
}

func NewRepository(db *sql.DB) Repository {
	return Repository{
		MovieRepo: MovieRepository{db},
		ActorRepo: ActorRepository{db},
		GenreRepo: GenreRepository{db},
	}
}

type ActorRepositoryInterface interface {
	FindAllActors() ([]Actor, error)
	CreateActor(Actor) (bool, error)
	FindActorByNameAndBirthDate(Actor) (Actor, error)
}

type GenreRepositoryInterface interface {
	FindAllGenres() ([]Genre, error)
	CreateGEnre(Genre) (bool, error)
}

type MovieRepositoryInterface interface {
	FindAllMovies() ([]Movie, error)
	CreateActor(Movie) (bool, error)
}
