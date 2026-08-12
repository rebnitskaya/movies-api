package repository

import (
	"database/sql"
	"movies_api/models"
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
	FindAllActors() ([]models.Actor, error)
	CreateActor(models.Actor) (bool, error)
	FindActorByNameAndBirthDate(models.Actor) (models.Actor, error)
}

type GenreRepositoryInterface interface {
	FindAllGenres() ([]models.Genre, error)
	CreateGEnre(models.Genre) (bool, error)
}

type MovieRepositoryInterface interface {
	FindAllMovies() ([]models.Movie, error)
	CreateActor(models.Movie) (bool, error)
}
