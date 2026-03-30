// Package handlers defines http handlers
package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/underark/stone-collector/internal/game"
	"github.com/underark/stone-collector/internal/models/locations"
	"github.com/underark/stone-collector/internal/models/state"
	"github.com/underark/stone-collector/internal/models/stones"
	"github.com/underark/stone-collector/internal/models/user"
)

func HomeHandler(userID int) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			defer conn.Close(context.Background())

			rows, err := conn.Query(context.Background(), "SELECT sum(amount) AS stones FROM stones WHERE owner_id = $1;", userID)
			if err != nil {
				fmt.Printf("Error collecting stone total: %s", err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			state, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[state.State])
			if err != nil {
				fmt.Printf("Error collecting stone total: %s", err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(state)
		}
	}
}

func TickHandler(userID int) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			defer conn.Close(context.Background())

			u, err := loadUser(userID, conn)
			if err != nil {
				fmt.Printf("Error loading user: %s", err.Error())
			}

			ticks, err := game.TicksSince(u)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

			for range ticks {
				stone := stones.New(locations.Park)
				result, err := conn.Exec(context.Background(), "UPDATE stones SET amount = amount + 1 WHERE owner_id = $1 AND material = $2;", userID, stone.Material)
				if err != nil {
					fmt.Printf("Error updating database: %s", err.Error())
					return
				}

				if result.RowsAffected() == 0 {
					conn.Exec(context.Background(), "INSERT INTO stones (owner_id, material, amount) VALUES ($1, $2, 1);", userID, stone.Material)
				}
			}

			newTicks, err := u.ConsumeTicks(ticks)
			if err != nil {
				fmt.Printf("Error generating new ticks: %s", err.Error())
				return
			}

			_, err = conn.Exec(context.Background(), "UPDATE users SET last_tick = $2 WHERE id = $1", userID, newTicks)
			if err != nil {
				fmt.Printf("Error updating database last_ticks: %s", err.Error())
				return
			}

			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		}
	}
}

func loadUser(userID int, dbConn *pgx.Conn) (user.User, error) {
	rows, err := dbConn.Query(context.Background(), "SELECT id, name, last_tick::text FROM users WHERE id = $1", userID)
	if err != nil {
		return user.User{}, err
	}

	u, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[user.User])
	if err != nil {
		return user.User{}, err
	}

	return u, nil
}
