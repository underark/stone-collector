package router

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/sessions"
	"github.com/underark/stone-collector/internal/service/game"
	"github.com/underark/stone-collector/internal/service/store"
	"github.com/underark/stone-collector/web/handlers"
	"github.com/underark/stone-collector/web/middleware"
)

type Router struct {
	mux            *http.ServeMux
	g              *game.GameService
	cookieStore    *sessions.CookieStore
	authMiddleware func(http.Handler) http.Handler
}

func New() *Router {
	store, err := store.New()
	if err != nil {
		log.Fatalf("Cannot initialize store: %s\n", err.Error())
	}

	cookieStore := sessions.NewCookieStore([]byte(os.Getenv("COOKIE_SECRET")))

	r := &Router{
		mux:            http.NewServeMux(),
		g:              game.New(store),
		cookieStore:    cookieStore,
		authMiddleware: middleware.NewAuthMiddleware(store, cookieStore),
	}

	attachHandlers(r)
	return r
}

func (r *Router) Serve() error {
	return http.ListenAndServe(":8080", r.authMiddleware(r.mux))
}

func attachHandlers(r *Router) {
	r.mux.Handle("/home", handlers.HomeHandler(r.g))
	r.mux.Handle("GET /start", handlers.StartHandler(r.g, r.cookieStore))
	r.mux.Handle("GET /tick", handlers.TickHandler(r.g))
	r.mux.Handle("GET /inventory", handlers.InventoryHandler(r.g))
	r.mux.Handle("GET /trade", handlers.TradeMenuHandler(r.g))
	r.mux.Handle("POST /trade", handlers.TradeHandler(r.g))
	r.mux.Handle("/web/static/", http.HandlerFunc(handlers.StaticHandler))
}
