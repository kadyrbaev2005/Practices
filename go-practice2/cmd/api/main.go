package main

import (
	"log"
	"net/http"

	"go-practice2/internal/handler"
	"go-practice2/internal/middleware"
)

func main() {
	mux := http.NewServeMux()

	// routes
	mux.Handle("/user", middleware.Auth(http.HandlerFunc(handler.UserHandler)))

	log.Println("Starting server on :8080")
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal(err)
	}
}