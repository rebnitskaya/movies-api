package service

import (
	"fmt"
	r "movies_api/repository"
	"time"
)

type ActorService struct {
	repo r.ActorRepository
}

func (s *ActorService) GetAllActors() ([]r.Actor, error) {
	return []r.Actor{}, nil
}

func (s *ActorService) CreateActor(actorData r.Actor) (bool, error) {
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

	actor, err := s.repo.FindActorByNameAndBirthDate(actorData.Name, actorData.BirthDate)
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
	s.repo.DeleteActorByID(id)
	return false, nil
}

func (s *ActorService) GetActor(id int) (r.Actor, error) {
	s.repo.DeleteActorByID(id)
	return r.Actor{}, nil
}

func (s *ActorService) PatchActor(actorData string) (r.Actor, error) {
	actorId := 0
	data := make(map[string]string)
	s.repo.ReplaceFieldsInActor(actorId, data)
	return r.Actor{}, nil
}

func (s *ActorService) GetActorsWithName(name string) ([]r.Actor, error) {
	s.repo.FindActorsByName(name)
	return []r.Actor{}, nil
}
