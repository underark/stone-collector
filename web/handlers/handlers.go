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
			if err != nil {
				fmt.Printf("Error updating ticks: %s\n", err.Error())
				return
			} else if ticks == 0 {
				fmt.Println("No ticks to process")
				return
			}

			fmt.Printf("Processing %d ticks\n", ticks)

			drops := game.GetDrops(ticks)

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
			w.WriteHeader(http.StatusOK)
		}
	}
}

func TradeHandler(userID int, tradeID int) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			defer conn.Close(context.Background())

			tx, err := conn.Begin(context.Background())
			defer tx.Rollback(context.Background())
			if err != nil {
				fmt.Printf("Error creating transaction: %s\n", err.Error())
				return
			}

			rows, err := tx.Query(context.Background(), "SELECT * FROM trades WHERE id = $1 FOR UPDATE;", tradeID)
			if err != nil {
				fmt.Printf("Error reading trade info: %s\n", err.Error())
				return
			}
			trade, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[state.Trade])
			if err != nil {
				fmt.Printf("Error scanning trade info to struct: %s\n", err.Error())
				return
			}

			rows, err = tx.Query(context.Background(), "SELECT * FROM stones WHERE owner_id = $1 AND material = $2 FOR UPDATE;", trade.OwnerID, trade.Material)
			if err != nil {
				fmt.Printf("Error getting stones for trade owner: %s\n", err.Error())
				return
			}
			ownerInv, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[state.Inventory])
			if err != nil {
				fmt.Printf("Error scannning owner inventory: %s\n", err.Error())
				return
			}

			rows, err = tx.Query(context.Background(), "SELECT * FROM stones WHERE owner_id = $1 AND material = $2 FOR UPDATE;", userID, trade.MaterialReq)
			if err != nil {
				fmt.Printf("Error getting stones for trade responder: %s\n", err.Error())
				return
			}
			responderInv, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[state.Inventory])
			if err != nil {
				fmt.Printf("Error getting scanning responder inventory: %s\n", err.Error())
				return
			}

			if ownerInv.Material == trade.Material {
				fmt.Println("EEEEE")
			}
			if (ownerInv.Material == trade.Material) && (ownerInv.Amount >= trade.Amount) {
				_, err = tx.Exec(context.Background(), "UPDATE stones SET amount = amount - $1 WHERE material = $2 AND owner_id = $3;", trade.Amount, trade.Material, trade.OwnerID)
				if err != nil {
					fmt.Printf("Error updating trade owner stone amount: %s\n", err.Error())
					return
				}

				_, err = tx.Exec(context.Background(), "UPDATE stones SET amount = amount + $1 WHERE material = $2 AND owner_id = $3;", trade.Amount, trade.Material, userID)
				if err != nil {
					fmt.Printf("Error updating trade responder stone amount: %s\n", err.Error())
					return
				}
			} else {
				fmt.Printf("Error: incorrect materials/amount gathered from database for trade owner: need %d %s have %d %s\n", trade.Amount, trade.Material, ownerInv.Amount, ownerInv.Material)
				return
			}

			if responderInv.Material == trade.MaterialReq && responderInv.Amount >= trade.AmountReq {
				_, err = tx.Exec(context.Background(), "UPDATE stones SET amount = amount - $1 WHERE material = $2 AND owner_id = $3;", trade.AmountReq, trade.MaterialReq, userID)
				if err != nil {
					fmt.Printf("Error updating trade responder stone amount 2: %s\n", err.Error())
					return
				}

				_, err = tx.Exec(context.Background(), "UPDATE stones SET amount = amount + $1 WHERE material = $2 AND owner_id = $3;", trade.AmountReq, trade.MaterialReq, trade.OwnerID)
				if err != nil {
					fmt.Printf("Error updating trade owner stone amount 2: %s\n", err.Error())
					return
				}
			} else {
				fmt.Println("Error: incorrect materials amount gathered from database for trade responder")
				return
			}

			err = tx.Commit(context.Background())
			if err != nil {
				fmt.Printf("Error comitting trade transaction: %s\n", err.Error())
				return
			}

			fmt.Printf("Trade %d successfully commmitted!\n", trade.ID)
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
