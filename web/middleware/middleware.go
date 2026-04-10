// Package middleware defines middleware to be used in the web app
package middleware

import (
	"net/http"

	"github.com/underark/stone-collector/internal/service/game"
	"github.com/underark/stone-collector/web/inject"
)

func CheckCookie(route http.Handler, g *game.GameService) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("stone-game-user")
		if err != nil {
			route.ServeHTTP(w, r)
			return
		}

		id, err := g.GetUserFromSession(c.Value)
		if err != nil {
			route.ServeHTTP(w, r)
			return
		}

		r = inject.UserID(r, id)
		route.ServeHTTP(w, r)
	})
}
