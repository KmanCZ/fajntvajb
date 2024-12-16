package api

import (
	"net/http"
)

func New() (*http.ServeMux, error) {
	r := http.NewServeMux()
	handlers, err := NewHandlers()
	if err != nil {
		return nil, err
	}

	r.HandleFunc("GET /auth", handlers.handleAuthPage)
	r.HandleFunc("GET /", handlers.handleLandingPage)

	return r, nil
}
