package models

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

type CreateGenreDto struct {
	Name string `json:"name"`
}
