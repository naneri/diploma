package main

import (
	"diploma/cmd/config"
	"diploma/cmd/middleware"
	"github.com/caarlos0/env/v6"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

var cfg config.Config

func main() {
	configErr := env.Parse(&cfg)

	if configErr != nil {
		log.Fatalf("error parsing config: %v", configErr)
	}

	r := mainHandler()

	log.Println("Server started at port " + cfg.ServerAddress)
	log.Fatal(http.ListenAndServe(cfg.ServerAddress, r))
}

func mainHandler() *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.GzipMiddleware)
	return r
}
