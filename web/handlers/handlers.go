// Package handlers defines http handlers
package handlers

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/underark/stone-collector/internal/game"
	"github.com/underark/stone-collector/internal/models/stones"
	"github.com/underark/stone-collector/internal/models/user"
)

func HomeHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			defer conn.Close(context.Background())
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
				stone := stones.New()
				result, err := conn.Exec(context.Background(), "UPDATE stones SET amount = amount + 1 WHERE owner_id = $1 AND material = $2;", userID, stone.Material)
				if err != nil {
					fmt.Printf("Error updating database: %s", err.Error())
					return
				}

				if result.RowsAffected() == 0 {
					conn.Exec(context.Background(), "INSERT INTO stones (owner_id, material, amount) VALUES ($1, $2, 1);", userID, stone.Material)
				}
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
