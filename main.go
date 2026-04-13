package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/underark/stone-collector/web/router"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading environment variables: %v", err)
	}

	r := router.New()
	err := r.Serve()
	if err != nil {
		os.Exit(1)
	}
}
