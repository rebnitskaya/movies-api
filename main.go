package main

import (
	"log"
	"net/http"

	"movies-api/handlers"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/movies", handlers.MoviesHandler)

	log.Println("Server started on localhost:8080")

	log.Fatal(http.ListenAndServe(":8080", mux))
}
