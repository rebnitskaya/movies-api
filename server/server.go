package server

import (
	"context"
	"log"
	"movies_api/middleware"
	"net"
	"net/http"
)

type Config struct {
	Host string
	Port string
}

func Server(ctx context.Context) *http.Server {
	mux := http.NewServeMux()

	RegisterRoutes(mux)

	cfg := Config{
		Host: "localhost",
		Port: "8080",
	}

	srv := &http.Server{
		Addr:    cfg.Host + ":" + cfg.Port,
		Handler: middleware.Recover(mux),
		BaseContext: func(_ net.Listener) context.Context {
			return ctx
		},
	}

	log.Println("Launching server at", srv.Addr)
	return srv
}
