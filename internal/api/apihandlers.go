package api

import (
	"fajntvajb/internal/logger"
	"net/http"
)

func (handlers *handlers) handleHTMXTest(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("test"))
	if err != nil {
		handleAPIError(w, err)
	}
}

func handleAPIError(w http.ResponseWriter, err error) {
	log := logger.Get()
	log.Error().Err(err).Msg("Failed to render page")
	http.Error(w, "Something went wrong", http.StatusInternalServerError)
}
