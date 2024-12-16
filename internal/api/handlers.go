package api

import (
	"fmt"
	"net/http"
)

func handleLandingPage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello to fajntvajb")
}
