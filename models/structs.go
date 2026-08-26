package models

import (
	"fmt"
	"time"
)

type Genre struct {
	Id     int     `json:"id"`
	Name   string  `json:"name"`
	Movies []Movie `json:"movies"`
}

type Actor struct {
	Id        int     `json:"id"`
	Name      string  `json:"name"`
	BirthDate string  `json:"birthDate"`
	Movies    []Movie `json:"movies"`
}

type Movie struct {
	Id          int     `json:"id"`
	Title       string  `json:"title"`
	ReleaseYear int     `json:"releaseYear"`
	Duration    int     `json:"duration"`
	Genres      []Genre `json:"genres"`
	Actors      []Actor `json:"actors"`
}

func ValidatePatchActor(fields map[string]string) error {
	if len(fields) == 0 {
		return fmt.Errorf("%w: no fields to update", ErrBadRequest)
	}

	for field, value := range fields {
		switch field {
		case "name":
			if value == "" {
				return fmt.Errorf("%w: actor name can't be empty", ErrBadRequest)
			}
		case "birthDate":
			if value == "" {
				return fmt.Errorf("%w: actor birth date can't be empty", ErrBadRequest)
			}
			minDate := time.Date(
				1895, 12, 28,
				0, 0, 0, 0,
				time.UTC,
			)

			birthDate, err := time.Parse("2006-01-02", value)
			if err != nil {
				return fmt.Errorf("%w: actors birth date must be in YYYY-MM-DD format: %w", ErrBadRequest, err)
			}

			if birthDate.Before(minDate) {
				return fmt.Errorf("%w: there was no actors back then", ErrBadRequest)
			}

			if birthDate.After(time.Now()) {
				return fmt.Errorf("%w: actors birth date can't be in the future", ErrBadRequest)
			}

		default:
			return fmt.Errorf("%w: invalid actor field: %s", ErrBadRequest, field)
		}
	}
	return nil
}
