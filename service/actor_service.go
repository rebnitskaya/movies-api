package service

import (
	"movies_api/models"
	"movies_api/repository"
)

type ActorService struct {
	repo repository.ActorRepository
}

// to check
var actors = []models.Actor{
	{
		Id:        1,
		Name:      "Seva Bondar",
		BirthDate: "12.04.1994",
	},
	{
		Id:        2,
		Name:      "Alina Rebnitskaya",
		BirthDate: "12.04.2004",
	},
}

func (s *ActorService) GetAllActors() ([]models.Actor, error) {
	return actors, nil
}

func NewActorService(repo repository.ActorRepository) *ActorService {
	return &ActorService{
		repo: repo,
	}
}
