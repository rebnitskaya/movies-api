package main

import (
	"context"
	"fmt"
	"log"
	"movies_api/server"
	"net/http"
)

// @title Movie Database API
// @version 1.0
// @description REST API for managing movies, genres and actors.
// @host localhost:8080
// @BasePath /
func main() {
	ctx := context.Background() //just for the start, will add some later

	srv := server.Server(ctx)

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		errMsg := fmt.Sprintf("Server error: %v", err.Error())
		log.Fatal(errMsg)
	}
}
