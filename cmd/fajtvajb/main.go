package main

import (
	"fajntvajb/internal/api"
	"fajntvajb/internal/logger"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	log := logger.Get()
	if err != nil {
		log.Error().Err(err).Msg("Failed to load .env file")
	}

	router, err := api.New()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create router")
		return
	}

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	server := http.Server{
		Addr:    ":" + port,
		Handler: logger.RequestLogger(router),
	}

	log.Info().Str("port", port).Msg("Listening on port " + port)
	log.Info().Msg("Server can be accesed on http://localhost:" + port)

	err = server.ListenAndServe()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to start server")
	}
}
