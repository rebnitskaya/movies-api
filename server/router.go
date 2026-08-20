package server

import (
	h "movies_api/handler"
	"net/http"
)

func RegisterRoutes(mux *http.ServeMux, movieHandler *h.MovieHandler, actorHandler *h.ActorHandler, genreHandler *h.GenreHandler) {
	mux.HandleFunc("/api", h.HandleRoot)

	//actors
	mux.HandleFunc("POST /api/actors", actorHandler.PostActor)
	mux.HandleFunc("GET /api/actors", actorHandler.GetAllActors)
	mux.HandleFunc("GET /api/actors/{id}", actorHandler.GetActor)
	mux.HandleFunc("PATCH /api/actors/{id}", actorHandler.PatchActor)
	mux.HandleFunc("DELETE /api/actors/{id}", actorHandler.DeleteActor)

	//movies
	mux.HandleFunc("POST /api/movies", movieHandler.PostMovie)
	mux.HandleFunc("GET /api/movies/", movieHandler.GetMovies)
	mux.HandleFunc("GET /api/movies/{id}", movieHandler.GetMovie)
	mux.HandleFunc("PATCH /api/movies/{id}", movieHandler.PatchMovie)
	mux.HandleFunc("DELETE /api/movies/{id}", movieHandler.DeleteMovie)
	mux.HandleFunc("GET /api/movies/{id}/actors", movieHandler.GetActorsInMovie)
	mux.HandleFunc("POST /api/movies/{movieID}/actors/{actorID}", movieHandler.PostActorToMovie)
	mux.HandleFunc("DELETE /api/movies/{movieID}/actors/{actorID}", movieHandler.DeleteActorFromMovie)

	//genres
	mux.HandleFunc("POST /api/genres", genreHandler.PostGenre)
	mux.HandleFunc("GET /api/genres", genreHandler.GetAllGenres)
	mux.HandleFunc("GET /api/genres/{id}", genreHandler.GetGenre)
	mux.HandleFunc("PATCH /api/genres", genreHandler.PatchGenre)
	mux.HandleFunc("DELETE /api/genres/{id}", genreHandler.DeleteGenre)

	// planning to use it for automatic documentation for endpoints
	// mux.Handle(
	// 	"/swagger/",
	// 	httpSwagger.Handler(
	// 		httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
	// 	),
	// )
}
