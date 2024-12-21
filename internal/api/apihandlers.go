package api

import (
	"fajntvajb/internal/logger"
	"net/http"
)

func (handlers *handlers) handleHTMXPostTest(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		handleAPIError(w, err)
	}
	value := r.FormValue("name")
	err = handlers.db.InsertRow(value)
	if err != nil {
		handleAPIError(w, err)
	}
	_, err = w.Write([]byte("<li>" + value + "</li>"))
	if err != nil {
		handleAPIError(w, err)
	}
}

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
