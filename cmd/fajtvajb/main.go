package main

import (
	"fajntvajb/internal/api"
	"fajntvajb/internal/logger"
	"net/http"
)

func main() {
	log := logger.Get()
	router, err := api.New()
	if err != nil {
		log.Error().Err(err).Msg("Failed to create router")
		return
	}

	server := http.Server{
		Addr:    ":8080",
		Handler: logger.RequestLogger(router),
	}

	log.Info().Msg("Listening on port 8080")
	log.Info().Msg("Server can be accesed on http://localhost:8080")
	err = server.ListenAndServe()
	if err != nil {
		log.Error().Err(err).Msg("Failed to start server")
	}
}
