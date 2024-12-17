package api

import (
	"fajntvajb/internal/files"
	"io/fs"
	"net/http"
)

func New() (*http.ServeMux, error) {
	r := http.NewServeMux()
	handlers, err := NewHandlers()
	if err != nil {
		return nil, err
	}

	static, err := fs.Sub(files.Files, "static")
	if err != nil {
		return nil, err
	}

	r.Handle("GET /static/", http.StripPrefix("/static", http.FileServerFS(static)))

	r.HandleFunc("GET /auth", handlers.handleAuthPage)
	r.HandleFunc("GET /", handlers.handleLandingPage)


	return r, nil
}
