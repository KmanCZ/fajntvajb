package main

import (
	"fajntvajb/internal/api"
	"fmt"
	"net/http"
)

func main() {
	router := api.New()

	server := http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	fmt.Println("Listening on port 8080")
	err := server.ListenAndServe()
	if err != nil {
		fmt.Printf("Something went wrong: %s", err.Error())
	}
}
