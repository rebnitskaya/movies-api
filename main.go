package main

import (
	"context"
	"fmt"
	"log"
	"movies_api/cli"
	"movies_api/server"
	"net/http"
	"os"
)

func main() {
	ctx := context.Background() //just for the start, will add some later

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
