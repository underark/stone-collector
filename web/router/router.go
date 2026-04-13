package router

import (
	"log"
	"net/http"

	"github.com/underark/stone-collector/internal/service/game"
	"github.com/underark/stone-collector/internal/service/store"
	"github.com/underark/stone-collector/web/handlers"
	"github.com/underark/stone-collector/web/middleware"
)

type Router struct {
	mux            *http.ServeMux
	g              *game.GameService
	authMiddleware func(http.Handler) http.Handler
}

func New() *Router {
	store, err := store.New()
	if err != nil {
		log.Fatalf("Cannot initialize store: %s\n", err.Error())
	}

	r := &Router{
		mux:            http.NewServeMux(),
		g:              game.New(store),
		authMiddleware: middleware.NewAuthMiddleware(store),
	}

	attachHandlers(r)
	return r
}

func (r *Router) Serve() error {
	return http.ListenAndServe(":8080", r.authMiddleware(r.mux))
}

func attachHandlers(r *Router) {
	r.mux.Handle("GET /tick", handlers.TickHandler(r.g))
	r.mux.Handle("GET /start", http.HandlerFunc(handlers.StartHandler))
}
