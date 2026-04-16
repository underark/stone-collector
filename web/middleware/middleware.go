// Package middleware defines middleware to be used in the web app
package middleware

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/sessions"
	"github.com/underark/stone-collector/internal/models"
	"github.com/underark/stone-collector/internal/service/game"
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
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			if session.IsNew {
				route.ServeHTTP(w, r)
				return
			}

			sessionID, ok := session.Values["session_id"]
			if !ok {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			s, err := db.GetUserFromSession(sessionID.(string))
			if err != nil {
				log.Printf("Error getting user data from session for id %s: %s\n", sessionID, err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			expired, err := isExpired(s.SessionExpiry)
			if err != nil {
				log.Printf("Error evaluating expiry for session %s: %s\n", sessionID, err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			if expired {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			fmt.Printf("User id is %d\n", s.ID)

			r = inject.UserID(r, s.ID)
			route.ServeHTTP(w, r)
		})
	}
}

func NewStateMiddleware(g *game.GameService) func(http.Handler) http.Handler {
	return func(route http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID := inject.GetUserID(r.Context())
			if userID == nil {
				route.ServeHTTP(w, r)
				return
			}

			state, err := g.GetUserState(userID.(int))
			if err != nil {
				log.Printf("Error getting user state for user %d: %s\n", userID, err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			state.ID = userID.(int)

			r = inject.State(r, state)
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
