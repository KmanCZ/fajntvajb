package main

import (
	"fajntvajb/internal/api"
	"fmt"
	"net/http"
)

func main() {
	router, err := api.New()
	if err != nil {
		fmt.Printf("Something went wrong: %s", err.Error())
	}

	server := http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	fmt.Println("Listening on port 8080")
	fmt.Println("Server can be accesed on http://localhost:8080")
	err = server.ListenAndServe()
	if err != nil {
		fmt.Printf("Something went wrong: %s", err.Error())
	}
}
