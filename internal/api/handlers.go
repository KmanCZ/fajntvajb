package api

import (
	"fajntvajb/internal/files"
	"html/template"
	"net/http"
)

type handlers struct {
	tmpl *template.Template
}

func NewHandlers() (*handlers, error) {
	tmpl, err := template.New("layout.html").ParseFS(files.Files, "templates/*/*.html")
	if err != nil {
		return nil, err
	}

	res := handlers{
		tmpl: tmpl,
	}
	return &res, nil
}

func (handlers *handlers) handleLandingPage(w http.ResponseWriter, r *http.Request) {
	//TODO: handle errors
	handlers.tmpl.ExecuteTemplate(w, "index.html", nil)
}

func (handlers *handlers) handleAuthPage(w http.ResponseWriter, r *http.Request) {
	//TODO: handle errors
	handlers.tmpl.ExecuteTemplate(w, "profile.html", struct {
		Name string
	}{
		Name: "john doe",
	})

}
