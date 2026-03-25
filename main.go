package main

import (
	"log"
	"net/http"

	"github.com/underark/stone-collector/web/handlers"
)

func main() {
	http.HandleFunc("/get", handlers.GetHandler())
	log.Fatal(http.ListenAndServe(":8080", nil))
}
