// Package middleware defines middleware to be used in the web app
package middleware

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5"
)

func CheckCookie(route http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := pgx.Connect(r.Context(), os.Getenv("DATABASE_URL"))
		if err != nil {
			fmt.Printf("Error connecting to database in CheckCookie: %s\n", err.Error())
			return
		}
		defer conn.Close(r.Context())

		c, err := r.Cookie("stone-game-user")
		if err != nil {
			fmt.Printf("Error getting cookie info: %s\n", err.Error())
			return
		}

		row, err := conn.Query(r.Context(), "SELECT id FROM users WHERE session_id = $1;", c.Value)
		if err != nil {
			fmt.Printf("Error getting user info from database: %s\n", err.Error())
			return
		}
		user, err := pgx.CollectOneRow(row, pgx.RowToStructByNameLax[user.User])
		if err != nil {
			fmt.Printf("Error scanning user info to struct: %s\n", err.Error())
			return
		}

		ctx := context.WithValue(r.Context(), k("userID"), user.ID)
		r = r.WithContext(ctx)
		route.ServeHTTP(w, r)
	})
}
