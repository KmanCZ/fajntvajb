package api

import (
	"fajntvajb/internal/files"
	"fajntvajb/internal/logger"
	"net/http"
)

func (handlers *handlers) handleLandingPage(w http.ResponseWriter, r *http.Request) {
	path := "index"
	if r.URL.Path != "/" {
		w.WriteHeader(http.StatusNotFound)
		path = "404"
	}

	err := handlers.tmpl.Render(w, r, path, nil)
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
