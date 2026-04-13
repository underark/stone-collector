// Package middleware defines middleware to be used in the web app
package middleware

import (
	"net/http"

	"github.com/underark/stone-collector/web/inject"
)

type authDB interface {
	GetUserFromSession(session string) (int, error)
}

func NewAuthMiddleware(db authDB) func(http.Handler) http.Handler {
	return func(route http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c, err := r.Cookie("stone-game-user")
			if err != nil {
				route.ServeHTTP(w, r)
				return
			}

			id, err := db.GetUserFromSession(c.Value)
			if err != nil {
				route.ServeHTTP(w, r)
				return
			}

			r = inject.UserID(r, id)
			route.ServeHTTP(w, r)
		})
	}
}
