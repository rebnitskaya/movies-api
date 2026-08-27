package server

import (
	h "movies_api/handler"
	m "movies_api/middleware"
	"net/http"
)

func RegisterRoutes(mux *http.ServeMux, movieHandler *h.MovieHandler, actorHandler *h.ActorHandler, genreHandler *h.GenreHandler) {
	mux.HandleFunc("/api", h.HandleRoot)

	//actors
	mux.HandleFunc("POST /api/actors", m.GlobalErrorHandler(actorHandler.PostActor))
	mux.HandleFunc("GET /api/actors", m.GlobalErrorHandler(actorHandler.GetAllActors))
	mux.HandleFunc("GET /api/actors/{id}", m.GlobalErrorHandler(actorHandler.GetActor))
	mux.HandleFunc("PATCH /api/actors/{id}", m.GlobalErrorHandler(actorHandler.PatchActor))
	mux.HandleFunc("DELETE /api/actors/{id}", m.GlobalErrorHandler(actorHandler.DeleteActor))

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
	mux.HandleFunc("POST /api/genres", m.GlobalErrorHandler(genreHandler.PostGenre))
	mux.HandleFunc("GET /api/genres", m.GlobalErrorHandler(genreHandler.GetAllGenres))
	mux.HandleFunc("GET /api/genres/{id}", m.GlobalErrorHandler(genreHandler.GetGenre))
	mux.HandleFunc("PATCH /api/genres", m.GlobalErrorHandler(genreHandler.PatchGenre))
	mux.HandleFunc("DELETE /api/genres/{id}", m.GlobalErrorHandler(genreHandler.DeleteGenre))
	mux.HandleFunc("POST /api/movies/{movieID}/genres/{genreID}", m.GlobalErrorHandler(movieHandler.PostGenreToMovie))
	mux.HandleFunc("DELETE /api/movies/{movieID}/genres/{genreID}", m.GlobalErrorHandler(movieHandler.DeleteGenreFromMovie))

	// planning to use it for automatic documentation for endpoints
	// mux.Handle(
	// 	"/swagger/",
	// 	httpSwagger.Handler(
	// 		httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
	// 	),
	// )
}
