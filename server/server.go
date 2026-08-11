package server

import (
	"context"
	"log"
	"movies_api/handler"
	"movies_api/middleware"
	"movies_api/repository"
	"movies_api/service"
	"net"
	"net/http"
)

type Config struct {
	Host string
	Port string
}

func Server(ctx context.Context) *http.Server {
	mux := http.NewServeMux()

	cfg := Config{
		Host: "localhost",
		Port: "8080",
	}

	dependencyWiring(mux)

	srv := &http.Server{
		Addr:    cfg.Host + ":" + cfg.Port,
		Handler: middleware.Recover(mux),
		BaseContext: func(_ net.Listener) context.Context {
			return ctx
		},
	}

	log.Println("Launching server at", srv.Addr)
	return srv
}

func dependencyWiring(mux *http.ServeMux) *http.ServeMux {
	repo := repository.NewRepository()
	movieService := service.NewMovieService(repo.MovieRepo)
	actorService := service.NewActorService(repo.ActorRepo)
	genreService := service.NewGenreService(repo.GenreRepo)

	movieHandler := handler.NewMovieHandler(movieService)
	actorHandler := handler.NewActorHandler(actorService)
	genreHandler := handler.NewGenreHandler(genreService)

	RegisterRoutes(mux, movieHandler, actorHandler, genreHandler)

	return mux
}
