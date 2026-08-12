package server

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"movies_api/db"
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

func Server(ctx context.Context) (*http.Server, error) {
	mux := http.NewServeMux()

	cfg := Config{
		Host: "localhost",
		Port: "8080",
	}

	dataBase, err := db.OpenDB("movies_api.db")
	if err != nil {
		return nil, fmt.Errorf("Can't run the db: %w", err)
	}

	if err := db.CreateTables(dataBase); err != nil {
		return nil, fmt.Errorf("Can't initiate tables: %w", err)
	}

	dependencyWiring(mux, dataBase)

	srv := &http.Server{
		Addr:    cfg.Host + ":" + cfg.Port,
		Handler: middleware.Recover(mux),
		BaseContext: func(_ net.Listener) context.Context {
			return ctx
		},
	}

	log.Println("Launching server at", srv.Addr)
	return srv, nil
}

func dependencyWiring(mux *http.ServeMux, db *sql.DB) *http.ServeMux {
	repo := repository.NewRepository(db)

	movieService := service.NewMovieService(repo.MovieRepo)
	actorService := service.NewActorService(repo.ActorRepo)
	genreService := service.NewGenreService(repo.GenreRepo)

	movieHandler := handler.NewMovieHandler(movieService)
	actorHandler := handler.NewActorHandler(actorService)
	genreHandler := handler.NewGenreHandler(genreService)

	RegisterRoutes(mux, movieHandler, actorHandler, genreHandler)

	return mux
}
