package server

import (
	h "movies_api/handler"
	m "movies_api/middleware"
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
	mux.HandleFunc("POST /api/movies", m.GlobalErrorHandler(movieHandler.PostMovie))
	mux.HandleFunc("GET /api/movies/", m.GlobalErrorHandler(movieHandler.GetMovies))
	mux.HandleFunc("GET /api/movies/{id}", m.GlobalErrorHandler(movieHandler.GetMovie))
	mux.HandleFunc("PATCH /api/movies/{id}", m.GlobalErrorHandler(movieHandler.PatchMovie))
	mux.HandleFunc("DELETE /api/movies/{id}", m.GlobalErrorHandler(movieHandler.DeleteMovie))
	mux.HandleFunc("GET /api/movies/{id}/actors", m.GlobalErrorHandler(movieHandler.GetActorsInMovie))
	mux.HandleFunc("POST /api/movies/{movieID}/actors/{actorID}", m.GlobalErrorHandler(movieHandler.PostActorToMovie))
	mux.HandleFunc("DELETE /api/movies/{movieID}/actors/{actorID}", m.GlobalErrorHandler(movieHandler.DeleteActorFromMovie))

	//genres
	mux.HandleFunc("POST /api/genres", genreHandler.PostGenre)
	mux.HandleFunc("GET /api/genres", genreHandler.GetAllGenres)
	mux.HandleFunc("GET /api/genres/{id}", genreHandler.GetGenre)
	mux.HandleFunc("PATCH /api/genres", genreHandler.PatchGenre)
	mux.HandleFunc("DELETE /api/genres/{id}", genreHandler.DeleteGenre)
	mux.HandleFunc("POST /api/movies/{movieID}/genres/{genreID}", movieHandler.PostGenreToMovie)
	mux.HandleFunc("DELETE /api/movies/{movieID}/genres/{genreID}", movieHandler.DeleteGenreFromMovie)

	// planning to use it for automatic documentation for endpoints
	// mux.Handle(
	// 	"/swagger/",
	// 	httpSwagger.Handler(
	// 		httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
	// 	),
	// )
}
