package server

import (
	"movies_api/handler"
	h "movies_api/handler"
	"net/http"

	_ "movies_api/docs"

	httpSwagger "github.com/swaggo/http-swagger"
)

func RegisterRoutes(mux *http.ServeMux, movieHandler *h.MovieHandler, actorHandler *h.ActorHandler, genreHandler *h.GenreHandler) {
	mux.HandleFunc("/", handler.HandleRoot)
	mux.HandleFunc("/movies", movieHandler.GetAllMovies)
	mux.HandleFunc("/actors", actorHandler.ActorHandlerRouter)
	mux.HandleFunc("/genres", genreHandler.GetAllGenres)

	mux.Handle(
		"/swagger/",
		httpSwagger.Handler(
			httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
		),
	)
}
