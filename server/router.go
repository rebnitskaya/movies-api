package server

import (
	h "movies_api/handler"
	hf "movies_api/helper"
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger"
)

func RegisterRoutes(mux *http.ServeMux, movieHandler *h.MovieHandler, actorHandler *h.ActorHandler, genreHandler *h.GenreHandler) {
	mux.HandleFunc("/api", h.HandleRoot)

	//actors
	mux.HandleFunc("POST /api/actors", hf.GlobalErrorHandler(actorHandler.PostActor))
	mux.HandleFunc("GET /api/actors", hf.GlobalErrorHandler(actorHandler.GetAllActors))
	mux.HandleFunc("GET /api/actors/{id}", hf.GlobalErrorHandler(actorHandler.GetActor))
	mux.HandleFunc("PATCH /api/actors/{id}", hf.GlobalErrorHandler(actorHandler.PatchActor))
	mux.HandleFunc("DELETE /api/actors/{id}", hf.GlobalErrorHandler(actorHandler.DeleteActor))

	//movies
	mux.HandleFunc("POST /api/movies", hf.GlobalErrorHandler(movieHandler.PostMovie))
	mux.HandleFunc("GET /api/movies/", hf.GlobalErrorHandler(movieHandler.GetMovies))
	mux.HandleFunc("GET /api/movies/{id}", hf.GlobalErrorHandler(movieHandler.GetMovie))
	mux.HandleFunc("PATCH /api/movies/{id}", hf.GlobalErrorHandler(movieHandler.PatchMovie))
	mux.HandleFunc("DELETE /api/movies/{id}", hf.GlobalErrorHandler(movieHandler.DeleteMovie))
	mux.HandleFunc("GET /api/movies/{id}/actors", hf.GlobalErrorHandler(movieHandler.GetActorsInMovie))
	mux.HandleFunc("POST /api/movies/{movieID}/actors/{actorID}", hf.GlobalErrorHandler(movieHandler.PostActorToMovie))
	mux.HandleFunc("DELETE /api/movies/{movieID}/actors/{actorID}", hf.GlobalErrorHandler(movieHandler.DeleteActorFromMovie))
	mux.HandleFunc("GET /api/movies/search", hf.GlobalErrorHandler(movieHandler.SearchByTitle))

	//genres
	mux.HandleFunc("POST /api/genres", hf.GlobalErrorHandler(genreHandler.PostGenre))
	mux.HandleFunc("GET /api/genres", hf.GlobalErrorHandler(genreHandler.GetAllGenres))
	mux.HandleFunc("GET /api/genres/{id}", hf.GlobalErrorHandler(genreHandler.GetGenre))
	mux.HandleFunc("PATCH /api/genres/{id}", hf.GlobalErrorHandler(genreHandler.PatchGenre))
	mux.HandleFunc("DELETE /api/genres/{id}", hf.GlobalErrorHandler(genreHandler.DeleteGenre))
	mux.HandleFunc("POST /api/movies/{movieID}/genres/{genreID}", hf.GlobalErrorHandler(movieHandler.PostGenreToMovie))
	mux.HandleFunc("DELETE /api/movies/{movieID}/genres/{genreID}", hf.GlobalErrorHandler(movieHandler.DeleteGenreFromMovie))

	// automatic documentation for endpoints
	mux.Handle(
		"/swagger/",
		httpSwagger.Handler(
			httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
		),
	)
}
