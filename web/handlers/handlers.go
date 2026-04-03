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

			// TODO: roll workers and user info into the same query with a table join + extending state struct
			ticks := updateUserTicks(userID, conn)

			workersInfo, err := conn.Query(context.Background(), "SELECT location_id FROM workers WHERE owner_id = $1;", userID)
			if err != nil {
				fmt.Printf("Error getting workers: %s", err.Error())
				return
			}

			workers, err := pgx.CollectRows(workersInfo, pgx.RowToStructByName[state.Worker])
			if err != nil {
				fmt.Printf("Error scanning workers: %s", err.Error())
				return
			}

			fmt.Printf("Processing %d ticks\n", ticks)

			for _, w := range workers {
				drops, err := game.DropsFromLocation(w.LocationID, ticks)
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
}

func StaticHandler(w http.ResponseWriter, r *http.Request) {
	data, err := os.ReadFile(r.URL.Path[1:])
	if err != nil {
		fmt.Printf("Error loading static files: %s\n", err.Error())
	}

	w.Header().Set("Content-Type", "text/css")

	w.Write(data)
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

func updateUserTicks(id int, dbConn *pgx.Conn) (ticks int) {
	tx, _ := dbConn.Begin(context.Background())
	rows, _ := tx.Query(context.Background(), "SELECT id, name, last_tick::text FROM users WHERE id = $1 FOR UPDATE NOWAIT", id)

	u, _ := pgx.CollectOneRow(rows, pgx.RowToStructByName[user.User])

	ticks, _ = game.TicksSince(u)

	newTicks, _ := u.ConsumeTicks(ticks)

	tx.Exec(context.Background(), "UPDATE users SET last_tick = $2 WHERE id = $1", id, newTicks)

	tx.Commit(context.Background())

	return
}
