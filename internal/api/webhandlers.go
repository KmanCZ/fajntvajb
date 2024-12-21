package api

import (
	"fajntvajb/internal/files"
	"fajntvajb/internal/logger"
	"net/http"
)

func (handlers *handlers) handleLandingPage(w http.ResponseWriter, r *http.Request) {
	err := handlers.tmpl.Render(w, "index", nil)
	if err != nil {
		handleWebError(w, err)
	}
}

func (handlers *handlers) handleAuthPage(w http.ResponseWriter, r *http.Request) {

	rows, err := handlers.db.GetRows()
	if err != nil {
		handleWebError(w, err)
		return
	}
	err = handlers.tmpl.Render(w, "auth", struct {
		Name string
		Auth bool
		Rows []string
	}{
		Name: "john doe",
		Auth: true,
		Rows: rows,
	})

	if err != nil {
		handleWebError(w, err)
	}
}

func handleWebError(w http.ResponseWriter, err error) {
	log := logger.Get()
	log.Error().Err(err).Msg("Failed to render page")
	file, err := files.Files.ReadFile("templates/pages/error.html")
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusInternalServerError)
	_, err = w.Write(file)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
	}
}
