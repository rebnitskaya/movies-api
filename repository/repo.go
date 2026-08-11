package repository

type Repository struct {
	MovieRepo MovieRepository
	ActorRepo ActorRepository
	GenreRepo GenreRepository
}

func NewRepository() Repository {
	return Repository{
		MovieRepo: MovieRepository{},
		ActorRepo: ActorRepository{},
		GenreRepo: GenreRepository{},
	}
}
