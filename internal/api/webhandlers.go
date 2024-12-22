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

func (handlers *handlers) handleRegisterPage(w http.ResponseWriter, _ *http.Request) {
	err := handlers.tmpl.Render(w, "register", nil)
	if err != nil {
		handleWebError(w, err)
	}
}

func (handlers *handlers) handleLoginPage(w http.ResponseWriter, _ *http.Request) {
	err := handlers.tmpl.Render(w, "login", nil)
	if err != nil {
		handleWebError(w, err)
	}
}

func (handlers *handlers) handleAuthPage(w http.ResponseWriter, r *http.Request) {
	err := handlers.tmpl.Render(w, "auth", struct {
		Name string
		Auth bool
	}{
		Name: "john doe",
		Auth: true,
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
