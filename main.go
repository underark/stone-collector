package main

import (
	"log"
	"net/http"

	"github.com/underark/stone-collector/internal/models/stones"
	"github.com/underark/stone-collector/web/handlers"
)

var storage = make([]stones.Stone, 0)

func main() {
	http.HandleFunc("/get", handlers.GetHandler(storage))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
