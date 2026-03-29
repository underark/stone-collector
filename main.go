package main

import (
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/underark/stone-collector/web/handlers"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading environment variables: %v", err)
	}

	http.HandleFunc("/get", handlers.GetHandler(1))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
