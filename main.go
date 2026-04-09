package main

import (
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/underark/stone-collector/web/handlers"
	"github.com/underark/stone-collector/web/middleware"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading environment variables: %v", err)
	}

	mux := http.NewServeMux()
	mux.Handle("GET /tick", middleware.CheckCookie(http.HandlerFunc(handlers.TickHandler)))
	mux.Handle("GET /start", http.HandlerFunc(handlers.StartHandler))
	http.HandleFunc("/home", handlers.HomeHandler(1))
	http.HandleFunc("/trade", handlers.TradeHandler(4, 1))
	http.HandleFunc("/web/static/", handlers.StaticHandler)
	log.Fatal(http.ListenAndServe(":8080", mux))
}
