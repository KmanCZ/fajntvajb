package api

import (
	"fajntvajb/internal/files"
	"fajntvajb/internal/logger"
	"io/fs"
	"net/http"
)

func New() (*http.ServeMux, error) {
	log := logger.Get()

	r := http.NewServeMux()
	handlers, err := NewHandlers()
	if err != nil {
		log.Error().Err(err).Msg("Failed to create handlers")
		return nil, err
	}

	static, err := fs.Sub(files.Files, "static")
	if err != nil {
		log.Error().Err(err).Msg("Failed to create static file server")
		return nil, err
	}

	r.Handle("GET /static/", http.StripPrefix("/static", http.FileServerFS(static)))

	r.HandleFunc("GET /auth", handlers.handleAuthPage)
	r.HandleFunc("GET /", handlers.handleLandingPage)
	r.HandleFunc("GET /test", handlers.handleHTMXTest)

	return r, nil
}
