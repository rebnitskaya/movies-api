package handler

import (
	"movies_api/service"
	"net/http"
)

type AppHandler func(w http.ResponseWriter, r *http.Request) error

type MovieHandler struct {
	service *service.MovieService
}

type ActorHandler struct {
	service *service.ActorService
}

type GenreHandler struct {
	service *service.GenreService
}

func NewMovieHandler(service *service.MovieService) *MovieHandler {
	return &MovieHandler{
		service: service,
	}
}

func NewActorHandler(service *service.ActorService) *ActorHandler {
	return &ActorHandler{
		service: service,
	}
}

func NewGenreHandler(service *service.GenreService) *GenreHandler {
	return &GenreHandler{
		service: service,
	}
}
