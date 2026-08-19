package models

import (
	"fmt"
)

type MovieDto struct {
	Title       string  `json:"title"`
	ReleaseYear int     `json:"releaseYear"`
	Duration    int     `json:"duration"`
	Genres      []Genre `json:"genres"`
	Actors      []Actor `json:"actors"`
}

func (m MovieDto) Validate() (bool, error) {
	if m.Title == "" {
		return false, fmt.Errorf("Title can't be empty.")
	}

	if m.ReleaseYear < 1885 || m.ReleaseYear > 2050 || m.ReleaseYear == 0 {
		return false, fmt.Errorf("Invalid release year.")
	}

	if m.Duration <= 0 || m.Duration > 720 {
		return false, fmt.Errorf("Invalid movie duration year.")
	}

	return true, nil
}
