package api

import (
	"html/template"
	"io/fs"
	"net/http"
	"os"
)

type handlers struct {
	files fs.FS
	tmpl  *template.Template
}

func NewHandlers() (*handlers, error) {
	files := os.DirFS(".")
	tmpl, err := template.New("layout.html").ParseFS(files, "templates/*.html", "templates/*/*.html")
	if err != nil {
		return nil, err
	}

	res := handlers{
		files: files,
		tmpl:  tmpl,
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
