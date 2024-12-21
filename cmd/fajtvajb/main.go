package main

import (
	"fajntvajb/internal/api"
	"fajntvajb/internal/logger"
	"net/http"
	"os"
)

func main() {
	log := logger.Get()
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
