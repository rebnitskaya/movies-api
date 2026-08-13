package server

import (
	h "movies_api/handler"
	"net/http"

	_ "movies_api/docs"

	httpSwagger "github.com/swaggo/http-swagger"
)

func RegisterRoutes(mux *http.ServeMux, movieHandler *h.MovieHandler, actorHandler *h.ActorHandler, genreHandler *h.GenreHandler) {
	mux.HandleFunc("/", h.HandleRoot)
	mux.HandleFunc("GET /actors", actorHandler.GetAllActors)
	mux.HandleFunc("POST /actors", actorHandler.PostActor)
	mux.HandleFunc("DELETE /actors/{id}", actorHandler.DeleteActor)
	mux.HandleFunc("/movies", movieHandler.GetAllMovies)
	mux.HandleFunc("/genres", genreHandler.GetAllGenres)

	mux.Handle(
		"/swagger/",
		httpSwagger.Handler(
			httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
		),
	)
}
