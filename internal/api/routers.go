package api

import (
	"net/http"
)

func New() *http.ServeMux {
	r := http.NewServeMux()
	handlers := NewHandlers()

	r.HandleFunc("GET /", handlers.handleLandingPage)

	return r
}
