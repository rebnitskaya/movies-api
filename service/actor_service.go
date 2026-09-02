package service

import (
	"context"
	"errors"
	"fmt"
	m "movies_api/models"
	r "movies_api/repository"
	"sort"
)

type ActorService struct {
	repo r.ActorRepository
}

func (s *ActorService) GetAllActors(page, limit int, ctx context.Context) (m.ActorsPaginated, error) {
	offset := page * limit

	actors, err := s.repo.FindAllActors(limit, offset, ctx)
	if err != nil {
		return m.ActorsPaginated{}, err
	}

	sort.Slice(actors, func(i, j int) bool {
		return actors[i].Id < actors[j].Id
	})

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

	pageSize := len(result)
	totalActors, err := s.repo.CountActors(ctx)
	if err != nil {
		return m.ActorsPaginated{}, err
	}

	totalPages := (totalActors + limit - 1) / limit

	if page > totalPages && totalPages > 0 {
		return m.ActorsPaginated{
			Actors: []m.ActorWithoutMoviesDto{},
		}, nil
	}

	return m.ActorsPaginated{
		Actors:     result,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

func (s *ActorService) CreateActor(actorData m.ActorDto, ctx context.Context) (m.Actor, error) {
	_, err := actorData.Validate()
	if err != nil {
		return m.Actor{}, err
	}

	_, err = s.repo.FindActorByNameAndBirthDate(actorData.Name, actorData.BirthDate, ctx)

	if err == nil {
		return m.Actor{}, fmt.Errorf("%w: this actor already has been made before", m.ErrBadRequest)
	}

	if !errors.Is(err, m.ErrActorNotFound) {
		return m.Actor{}, err
	}

	actor := m.Actor{
		Name:      actorData.Name,
		BirthDate: actorData.BirthDate,
	}

	createdActor, err := s.repo.CreateActor(actor, ctx)
	if err != nil {
		return m.Actor{}, err
	}

	return createdActor, nil
}

func (s *ActorService) DeleteActor(id int, force bool, ctx context.Context) (bool, error) {
	if id <= 0 {
		return false, fmt.Errorf("%w: invalid actor id.", m.ErrBadRequest)
	}

	actor, err := s.repo.FindActorByID(id, ctx)
	if err != nil {
		return false, err
	}

	if !force && len(actor.Movies) > 0 {
		return false, fmt.Errorf("%w: cannot delete actor '%s' because it has %d associated movies",
			m.ErrBadRequest,
			actor.Name,
			len(actor.Movies),
		)
	}

	_, err = s.repo.DeleteActorByID(id, ctx)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (s *ActorService) GetActor(id int, ctx context.Context) (m.ActorWithoutMoviesDto, error) {
	actor, err := s.repo.FindActorByID(id, ctx)
	if err != nil {
		return m.ActorWithoutMoviesDto{}, err
	}

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
	return actorDTO, nil
}

func (s *ActorService) PatchActor(actorId int, data map[string]string, ctx context.Context) (m.Actor, error) {
	err := m.ValidatePatchActor(data)
	if err != nil {
		return m.Actor{}, err
	}

	actor, err := s.repo.ReplaceFieldsInActor(actorId, data, ctx)
	if err != nil {
		return m.Actor{}, err
	}
	return actor, nil
}

func (s *ActorService) GetActorsWithName(name string, ctx context.Context) ([]m.ActorWithoutMoviesDto, error) {
	actors, err := s.repo.FindActorsByName(name, ctx)
	if err != nil {
		return nil, err
	}

	result := make([]m.ActorWithoutMoviesDto, 0, len(actors))

	for _, actor := range actors {
		actorDto := m.ActorWithoutMoviesDto{
			Id:        actor.Id,
			Name:      actor.Name,
			BirthDate: actor.BirthDate,
			Movies:    []m.MovieWithoutActorsDto{},
		}

		for _, movie := range actor.Movies {
			actorDto.Movies = append(actor.Movies, m.MovieWithoutActorsDto{
				Id:          movie.Id,
				Title:       movie.Title,
				ReleaseYear: movie.ReleaseYear,
				Duration:    movie.Duration,
			})
		}
		result = append(result, actorDto)
	}
	return result, nil
}
