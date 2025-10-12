package main

import (
	"languages-api/internal/config"
	"languages-api/internal/controller"
	"languages-api/internal/mariadb"
	"languages-api/internal/repo"
	"languages-api/internal/router"

	"net/http"

	"github.com/TV4/graceful"
	"github.com/rs/zerolog/log"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatal().Msgf("Error getting configurations: %v", err)
	}

	db, err := repo.New(cfg, mariadb.MariaConnector{})
	if err != nil {
		log.Fatal().Msgf("Error creating database client: %v", err)
	}

	defer func() {
		err := db.Close()
		if err != nil {
			log.Fatal().Msgf("Error closing the database client: %v", err)
		}
	}()

	ctrl := controller.New(cfg)

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router.CreateHandler(ctrl, db),
	}

	log.Info().Msgf("Listening on port %s", cfg.Port)
	graceful.LogListenAndServe(srv)
}
