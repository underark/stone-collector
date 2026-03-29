// Package handlers defines http handlers
package handlers

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/underark/stone-collector/internal/game"
	"github.com/underark/stone-collector/internal/models/user"
)

func GetHandler(userID int) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			defer conn.Close(context.Background())

			rows, err := conn.Query(context.Background(), "SELECT id, name, last_tick::text FROM users WHERE id = $1", userID)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			user, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[user.User])
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			ticks, err := game.TicksSince(user)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

			fmt.Printf("Date in read user is %s\n", user.LastTick)
			fmt.Printf("Ticks are %d\n", ticks)
		}
	}
}

func CheckDatabase() {
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer conn.Close(context.Background())

	rows, _ := conn.Query(context.Background(), "SELECT * FROM users WHERE id = 1;")
	user, _ := pgx.CollectOneRow(rows, pgx.RowToStructByName[user.User])
	fmt.Println(user)
}
