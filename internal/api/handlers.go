package api

import (
	"fajntvajb/internal/files"
	"fajntvajb/internal/files/templates"
	"net/http"
)

type handlers struct {
	tmpl *templates.Template
}

func NewHandlers() (*handlers, error) {
	templates, err := templates.New()
	if err != nil {
		return nil, err
	}

	res := handlers{
		tmpl: templates,
	}
	return &res, nil
}

func (handlers *handlers) handleLandingPage(w http.ResponseWriter, r *http.Request) {
	err := handlers.tmpl.Render(w, "index", nil)
	if err != nil {
		handleError(w)
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
		handleError(w)
	}
}

func handleError(w http.ResponseWriter) {
	file, err := files.Files.ReadFile("templates/pages/error.html")
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusInternalServerError)
	_, err = w.Write(file)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
}
