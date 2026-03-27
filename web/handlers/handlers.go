// Package handlers defines http handlers
package handlers

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/underark/stone-collector/internal/models/stones"
	"github.com/underark/stone-collector/internal/models/user"
)

func GetHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			stone := stones.New()
			conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			defer conn.Close(context.Background())

			_, err = conn.Exec(context.TODO(), "INSERT INTO stones (owner_id, material) VALUES (1, $1);", stone.Material)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
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
