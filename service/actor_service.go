package service

import (
	"fmt"
	m "movies_api/models"
	r "movies_api/repository"
)

type ActorService struct {
	repo r.ActorRepository
}

func (s *ActorService) GetAllActors() ([]m.Actor, error) {
	actors, err := s.repo.FindAllActors()
	if err != nil {
		return nil, err
	}
	return actors, nil
}

func (s *ActorService) CreateActor(actorData m.ActorDto) (bool, error) {
	_, err := actorData.Validate()
	if err != nil {
		return false, err
	}

	actorExist, err := s.repo.FindActorByNameAndBirthDate(actorData.Name, actorData.BirthDate)
	if err != nil {
		return false, err
	}

	if actorExist.BirthDate == actorData.BirthDate && actorExist.Name == actorData.Name {
		return false, fmt.Errorf("This actor already has been made before.")
	}

	actor := m.Actor{
		Name:      actorData.Name,
		BirthDate: actorData.BirthDate,
		Movies:    actorData.Movies,
	}

	ok, err := s.repo.CreateActor(actor)
	if err != nil {
		return false, err
	}

	return ok, nil
}

func (s *ActorService) DeleteActor(id int) (bool, error) {
	s.repo.DeleteActorByID(id)
	return false, nil
}

func (s *ActorService) GetActor(id int) (m.Actor, error) {
	actor, found := s.repo.FindActorByID(id)
	if !found {
		return m.Actor{}, fmt.Errorf("Actor not found")
	}
	return actor, nil
}

func (s *ActorService) PatchActor(actorData string) (m.Actor, error) {
	actorId := 0
	data := make(map[string]string)
	s.repo.ReplaceFieldsInActor(actorId, data)
	return m.Actor{}, nil
}

func (s *ActorService) GetActorsWithName(name string) ([]m.Actor, error) {
	actors, err := s.repo.FindActorsByName(name)
	if err != nil {
		return nil, err
	}
	return actors, nil
}
