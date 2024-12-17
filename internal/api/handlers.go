package api

import (
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
	//TODO: handle errors
	handlers.tmpl.Render(w, "index", nil)
}

func (handlers *handlers) handleAuthPage(w http.ResponseWriter, r *http.Request) {
	//TODO: handle errors
	handlers.tmpl.Render(w, "auth", struct {
		Name string
	}{
		Name: "john doe",
	})
}
