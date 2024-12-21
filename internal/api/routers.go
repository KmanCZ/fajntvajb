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

	// Serve static files from /static
	r.Handle("GET /static/", http.StripPrefix("/static", http.FileServerFS(static)))

	// Define page routes
	r.HandleFunc("GET /auth", handlers.handleAuthPage)
	r.HandleFunc("GET /", handlers.handleLandingPage)

	// Define API routes
	r.HandleFunc("POST /api/click", handlers.handleHTMXPostTest)
	r.HandleFunc("GET /api/test", handlers.handleHTMXTest)

	return r, nil
}
