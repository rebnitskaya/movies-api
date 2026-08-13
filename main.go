package main

import (
	"context"
	"fmt"
	"log"
	"movies_api/server"
	"net/http"
)

func main() {
	ctx := context.Background() //just for the start, will add some later

	srv, err := server.Server(ctx)
	if err != nil {
		log.Fatal(err)
	}
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		errMsg := fmt.Sprintf("Server error: %v", err.Error())
		log.Fatal(errMsg)
	}
}
