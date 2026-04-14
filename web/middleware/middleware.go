// Package middleware defines middleware to be used in the web app
package middleware

import (
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/underark/stone-collector/web/inject"
)

type authDB interface {
	GetUserFromSession(session string) (int, error)
}

func NewAuthMiddleware(db authDB, c *sessions.CookieStore) func(http.Handler) http.Handler {
	return func(route http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session, err := c.Get(r, "stone-collector")
			if err != nil {
				route.ServeHTTP(w, r)
				return
			}

			userID, ok := session.Values["user_id"]
			if !ok {
				route.ServeHTTP(w, r)
				return
			}

			id, err := db.GetUserFromSession(userID.(string))
			if err != nil {
				route.ServeHTTP(w, r)
				return
			}

			r = inject.UserID(r, id)
			route.ServeHTTP(w, r)
		})
	}
}
