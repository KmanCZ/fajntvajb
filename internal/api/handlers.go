package api

import (
	"html/template"
	"io/fs"
	"net/http"
	"os"
)

type handlers struct {
	files fs.FS
}

func NewHandlers() handlers {
	res := handlers{
		files: os.DirFS("."),
	}
	return res
}

func (handlers *handlers) handleLandingPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.New("").ParseFS(handlers.files, "templates/*")
	if err != nil {
		http.Error(w, "Something went wrong", 500)
	}

	tmpl.ExecuteTemplate(w, "index.html", struct {
		Name string
	}{
		Name: "John Doe",
	})
}
