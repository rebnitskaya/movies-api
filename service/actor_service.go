package service

import (
	"fmt"
	"movies_api/repository"
	"time"
)

type ActorService struct {
	repo repository.ActorRepository
}

func (s *ActorService) GetAllActors() ([]repository.Actor, error) {
	actors, err := s.repo.FindAllActors()
	if err != nil {
		return nil, err
	}
	return actors, nil
}

func (s *ActorService) CreateActor(actorData repository.Actor) (bool, error) {
	minDate := time.Date(
		1895, 12, 28,
		0, 0, 0, 0,
		time.UTC,
	)

	birthDate, err := time.Parse("2006-01-02", actorData.BirthDate)
	if err != nil {
		return false, fmt.Errorf("Actors birth date must be in YYYY-MM-DD format: %w", err)
	}

	if birthDate.Before(minDate) {
		return false, fmt.Errorf("There was no actors back then.")
	}

	if birthDate.After(time.Now()) {
		return false, fmt.Errorf("Actors birth date can't be in the future.")
	}

	actor, err := s.repo.FindActorByNameAndBirthDate(actorData)
	if actor.BirthDate == actorData.BirthDate && actor.Name == actorData.Name {
		return false, fmt.Errorf("This actor already has been made before.")
	}

	ok, err := s.repo.CreateActor(actorData)
	if err != nil {
		return false, err
	}

	return ok, nil
}

func (s *ActorService) DeleteActor(id int) (bool, error) {
	ok, err := s.repo.DeletActorById(id)
	if err != nil {
		return false, err
	}

	if !ok {
		return false, fmt.Errorf("Failed to delete user for some reason.")
	}

	return true, nil
}
