package main

import (
	"github.com/caarlos0/env/v6"
	"github.com/go-chi/chi/v5"
	"github.com/naneri/diploma/cmd/config"
	"github.com/naneri/diploma/cmd/middleware"
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

	//r.Post("/api/user/register", )
	return r
}
