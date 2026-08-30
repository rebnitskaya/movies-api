package main

import (
	"context"
	"fmt"
	"log"
	"movies_api/cli"
	_ "movies_api/docs"
	"movies_api/server"
	"net/http"
	"os"
)

// @title Movie Database API
// @version 1.0
// @description REST API for managing movies, genres and actors.
// @host localhost:8080
// @BasePath /api
func main() {
	ctx := context.Background()

	initiate, err := cli.FlagHandling()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}

	srv, err := server.Server(ctx, initiate)
	if err != nil {
		log.Fatal(err)
	}
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		errMsg := fmt.Sprintf("Server error: %v", err.Error())
		log.Fatal(errMsg)
	}
}
