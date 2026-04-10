package main

import (
	"log"
	"net/http"

	"github.com/joho/godotenv"
	h "github.com/underark/stone-collector/web/handlers"
	m "github.com/underark/stone-collector/web/middleware"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading environment variables: %v", err)
	}

	mux := http.NewServeMux()
	mux.Handle("GET /tick", m.CheckCookie(h.TickHandler())
	mux.Handle("GET /start", http.HandlerFunc(h.StartHandler))
	http.HandleFunc("/home", h.HomeHandler(1))
	http.HandleFunc("/trade", h.TradeHandler(4, 1))
	http.HandleFunc("/web/static/", h.StaticHandler)
	log.Fatal(http.ListenAndServe(":8080", mux))
}
