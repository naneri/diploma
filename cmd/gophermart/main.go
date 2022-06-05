package main

import (
	"flag"
	"github.com/caarlos0/env/v6"
	"github.com/go-chi/chi/v5"
	"github.com/naneri/diploma/cmd/gophermart/config"
	"github.com/naneri/diploma/cmd/gophermart/controllers"
	"github.com/naneri/diploma/cmd/gophermart/middleware"
	"github.com/naneri/diploma/internal/services"
	"github.com/naneri/diploma/internal/user"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
)

var cfg config.Config
var db *gorm.DB
var userRepo *user.DbRepository

func main() {

	configErr := env.Parse(&cfg)

	if configErr != nil {
		log.Fatalf("error parsing config: %v", configErr)
	}

	if flag.Lookup("a") == nil {
		flag.StringVar(&cfg.ServerAddress, "a", cfg.ServerAddress, "default server Port")
		flag.StringVar(&cfg.DatabaseAddress, "d", cfg.DatabaseAddress, "database DSN")
		flag.StringVar(&cfg.AccrualAddress, "f", cfg.AccrualAddress, "accrual system address")
	}

	flag.Parse()

	var dbErr error
	db, dbErr = gorm.Open(postgres.Open(cfg.DatabaseAddress), &gorm.Config{})
	if dbErr != nil {
		log.Fatalf("error connecting to database")
	}
	services.RunMigrations(db)
	userRepo = user.InitDatabaseRepository(db)

	r := mainHandler()

	log.Println("Server started at port " + cfg.ServerAddress)
	log.Fatal(http.ListenAndServe(cfg.ServerAddress, r))
}

func mainHandler() *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.GzipMiddleware)
	r.Use(middleware.DecompressGZIP)

	authController := controllers.AuthController{
		UserRepo: userRepo,
		Config:   &cfg,
	}

	r.Post("/api/user/register", authController.Register)

	r.Group(func(r chi.Router) {
		r.Use(middleware.IDMiddleware)
	})

	return r
}
