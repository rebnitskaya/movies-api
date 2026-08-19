package models

import (
	"fmt"
)

type GenreDto struct {
	Id     int               `json:"id"`
	Name   string            `json:"name"`
	Movies []MovieSummaryDto `json:"movies"`
}

type MovieSummaryDto struct {
	Id   int    `json:"movieId"`
	Name string `json:"movieName"`
}

func (m GenreDto) Validate() (bool, error) {
	fmt.Errorf("Some error")
	return true, nil
}
