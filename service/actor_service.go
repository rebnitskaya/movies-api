package service

import (
	"fmt"
	m "movies_api/models"
	r "movies_api/repository"
)

type ActorService struct {
	repo r.ActorRepository
}

func (s *ActorService) GetAllActors() ([]m.ActorWithoutMoviesDto, error) {
	actors, err := s.repo.FindAllActors()
	if err != nil {
		return nil, err
	}
	result := make([]m.ActorWithoutMoviesDto, 0, len(actors))
	for _, actor := range actors {
		actorDTO := m.ActorWithoutMoviesDto{
			Id:        actor.Id,
			Name:      actor.Name,
			BirthDate: actor.BirthDate,
			Movies:    []m.MovieWithoutActorsDto{},
		}

		for _, movie := range actor.Movies {
			actorDTO.Movies = append(actorDTO.Movies, m.MovieWithoutActorsDto{
				Id:          movie.Id,
				Title:       movie.Title,
				ReleaseYear: movie.ReleaseYear,
				Duration:    movie.Duration,
			})

		}
		result = append(result, actorDTO)
	}
	return result, nil
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

func (s *ActorService) DeleteActor(id int) error {
	deleted, err := s.repo.DeleteActorByID(id)
	if err != nil {
		return err
	}
	if !deleted {
		return fmt.Errorf("Actor not found")
	}
	return nil
}

func (s *ActorService) GetActor(id int) (m.Actor, error) {
	actor, err := s.repo.FindActorByID(id)
	if err != nil {
		return m.Actor{}, err
	}
	return actor, nil
}

func (s *ActorService) PatchActor(actorId int, data map[string]string) (m.Actor, error) {
	err := m.ValidatePatchActor(data)
	if err != nil {
		return m.Actor{}, err
	}

	actor, err := s.repo.ReplaceFieldsInActor(actorId, data)
	if err != nil {
		return m.Actor{}, err
	}
	return actor, nil
}

func (s *ActorService) GetActorsWithName(name string) ([]m.Actor, error) {
	actors, err := s.repo.FindActorsByName(name)
	if err != nil {
		return nil, err
	}
	return actors, nil
}
