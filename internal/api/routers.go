package api

import (
	"net/http"
)

func New() *http.ServeMux {
	r := http.NewServeMux()

	r.HandleFunc("GET /", handleLandingPage)

	return r
}
