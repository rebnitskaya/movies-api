package server

import (
	"movies_api/handlers"
	"net/http"

	_ "movies_api/docs"

	httpSwagger "github.com/swaggo/http-swagger"
)

func RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/", handlers.HandleIndex)
	mux.HandleFunc("GET /movies", handlers.GetMoviesHandler)
	mux.HandleFunc("POST /movies", handlers.CreateMovieHandler)

	mux.Handle(
		"/swagger/",
		httpSwagger.Handler(
			httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
		),
	)
}
