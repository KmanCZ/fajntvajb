package api

import (
	"html/template"
	"net/http"
	"os"
)

func handleLandingPage(w http.ResponseWriter, r *http.Request) {
	//TODO: move up so it can be open only once
	files := os.DirFS(".")
	tmpl, err := template.New("").ParseFS(files, "templates/*")
	if err != nil {
		http.Error(w, "Something went wrong", 500)
	}

	tmpl.ExecuteTemplate(w, "index.html", struct {
		Name string
	}{
		Name: "John Doe",
	})
}
