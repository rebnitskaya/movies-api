package models

import (
	"fmt"
	"time"
)

type ActorDto struct {
	Name      string  `json:"name"`
	BirthDate string  `json:"birthDate"`
	Movies    []Movie `json:"movies"`
}

type ActorInFilmDto struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	BirthDate string `json:"birthDate"`
}

func (m ActorDto) Validate() (bool, error) {
	minDate := time.Date(
		1895, 12, 28,
		0, 0, 0, 0,
		time.UTC,
	)

	birthDate, err := time.Parse("2006-01-02", m.BirthDate)
	if err != nil {
		return false, fmt.Errorf("Actors birth date must be in YYYY-MM-DD format: %w", err)
	}

	if birthDate.Before(minDate) {
		return false, fmt.Errorf("There was no actors back then.")
	}

	if birthDate.After(time.Now()) {
		return false, fmt.Errorf("Actors birth date can't be in the future.")
	}
	return true, nil
}
