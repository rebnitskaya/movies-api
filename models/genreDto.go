package models

import (
	"fmt"
)

type GenreDto struct {
	Id     int               `json:"id"`
	Name   string            `json:"name"`
	Movies []MovieSummaryDto `json:"movies"`
}

type GenreWithoutMovies struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type MovieSummaryDto struct {
	Id   int    `json:"movieId"`
	Name string `json:"movieName"`
}

type GenresPaginated struct {
	Genres     []GenreDto `json:"genres"`
	Page       int        `json:"page"`
	PageSize   int        `json:"pageSize"`
	TotalPages int        `json:"totalPages"`
}

func (m GenreDto) Validate() (bool, error) {
	fmt.Errorf("Some error")
	return true, nil
}
