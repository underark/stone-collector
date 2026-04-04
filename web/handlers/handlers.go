// Package handlers defines http handlers
package handlers

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/underark/stone-collector/internal/game"
	"github.com/underark/stone-collector/internal/models/state"
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

			t, err := template.ParseFiles("./web/templates/base.tmpl", "./web/templates/index.tmpl")
			if err != nil {
				fmt.Printf("Error rendering template: %s", err.Error())
				return
			}
			t.ExecuteTemplate(w, "base", state)
		}
	}
}

func TickHandler(userID int) func(w http.ResponseWriter, r *http.Request) {
	// TODO: simplify this
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			defer conn.Close(context.Background())

			ticks, err := updateUserTicks(userID, conn)
			if ticks <= 0 {
				return
			}

			fmt.Printf("Processing %d ticks\n", ticks)

			drops := game.GetDrops(ticks)
			if err != nil {
				fmt.Printf("Error generating drops: %s", err.Error())
				return
			}

			for _, d := range drops {
				result, err := conn.Exec(context.Background(), "UPDATE stones SET amount = amount + $3 WHERE owner_id = $1 AND material = $2;", userID, d.Material, d.Amount)
				if err != nil {
					fmt.Printf("Error updating database: %s", err.Error())
					return
				}

				if result.RowsAffected() == 0 {
					conn.Exec(context.Background(), "INSERT INTO stones (owner_id, material, amount) VALUES ($1, $2, $3);", userID, d.Material, d.Amount)
				}
			}
		}
		http.Redirect(w, r, "/home", http.StatusTemporaryRedirect)
	}
}

func StaticHandler(w http.ResponseWriter, r *http.Request) {
	data, err := os.ReadFile(r.URL.Path[1:])
	if err != nil {
		fmt.Printf("Error loading static files: %s\n", err.Error())
	}

	w.Header().Set("Content-Type", "text/css")

	w.Write(data)
}

// TODO: Handle error cases
func updateUserTicks(id int, dbConn *pgx.Conn) (int, error) {
	tx, err := dbConn.Begin(context.Background())
	defer tx.Rollback(context.Background())
	if err != nil {
		return 0, err
	}

	rows, err := tx.Query(context.Background(), "SELECT id, name, last_tick::text FROM users WHERE id = $1 FOR UPDATE;", id)
	if err != nil {
		return 0, err
	}

	u, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[user.User])
	if err != nil {
		return 0, err
	}

	ticks, err := game.TicksSince(u)
	if err != nil {
		return 0, err
	}

	newTicks, err := u.ConsumeTicks(ticks)
	if err != nil {
		return 0, err
	}
	_, err = tx.Exec(context.Background(), "UPDATE users SET last_tick = $2 WHERE id = $1", id, newTicks)
	if err != nil {
		return 0, err
	}

	tx.Commit(context.Background())
	return ticks, nil
}
