package models

import (
	"fmt"
	"time"
)

type ActorDto struct {
	Name      string `json:"name"`
	BirthDate string `json:"birthDate"`
}

type ActorInFilmDto struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	BirthDate string `json:"birthDate"`
}

type ActorWithoutMoviesDto struct {
	Id        int                     `json:"id"`
	Name      string                  `json:"name"`
	BirthDate string                  `json:"birthDate"`
	Movies    []MovieWithoutActorsDto `json:"movies"`
}

type ActorsPaginated struct {
	Actors     []ActorWithoutMoviesDto `json:"actors"`
	Page       int                     `json:"page"`
	PageSize   int                     `json:"pageSize"`
	TotalPages int                     `json:"totalPages"`
}

func (m ActorDto) Validate() (bool, error) {
	minDate := time.Date(
		1895, 12, 28,
		0, 0, 0, 0,
		time.UTC,
	)

	if m.Name == "" {
		return false, fmt.Errorf("%w: actor name cannot be empty", ErrBadRequest)
	}

	birthDate, err := time.Parse("2006-01-02", m.BirthDate)
	if err != nil {
		return false, fmt.Errorf("%w: actors birth date must be in YYYY-MM-DD format.", ErrBadRequest)
	}

	if birthDate.Before(minDate) {
		return false, fmt.Errorf("%w: there was no actors back then.", ErrBadRequest)
	}

	if birthDate.After(time.Now()) {
		return false, fmt.Errorf("%w: actors birth date can't be in the future.", ErrBadRequest)
	}
	return true, nil
}
