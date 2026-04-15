// Package middleware defines middleware to be used in the web app
package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/sessions"
	"github.com/underark/stone-collector/internal/models"
	"github.com/underark/stone-collector/web/inject"
)

type authDB interface {
	GetUserFromSession(session string) (models.Session, error)
}

func NewAuthMiddleware(db authDB, c *sessions.CookieStore) func(http.Handler) http.Handler {
	return func(route http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session, err := c.Get(r, "stone-collector")
			if err != nil {
				route.ServeHTTP(w, r)
				return
			}

			sessionID, ok := session.Values["session_id"]
			if !ok {
				route.ServeHTTP(w, r)
				return
			}

			s, err := db.GetUserFromSession(sessionID.(string))
			if err != nil {
				return
			}

			fmt.Printf("User id is %d\n", s.ID)

			r = inject.UserID(r, s.ID)
			route.ServeHTTP(w, r)
		})
	}
}

func isExpired(expiry string) (bool, error) {
	e, err := time.Parse(time.DateTime, expiry)
	if err != nil {
		return true, err
	}

	return time.Now().UTC().After(e), nil
}
