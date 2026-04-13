// Package handlers defines http handlers
package handlers

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/underark/stone-collector/internal/service/game"
	"github.com/underark/stone-collector/web/inject"
)

func HomeHandler(g *game.GameService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := inject.GetUserID(r.Context())
		if userID == nil {
			http.Redirect(w, r, "/start", http.StatusFound)
			return
		}

		state, err := g.GetUserState(userID.(int))
		t, err := template.ParseFiles("./web/templates/base.tmpl", "./web/templates/index.tmpl")
		if err != nil {
			fmt.Printf("Error rendering template: %s", err.Error())
			return
		}
		t.ExecuteTemplate(w, "base", state)
	}
}

func StartHandler(g *game.GameService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := g.InsertNewUser()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		c := http.Cookie{
			Name:     "stone-game-user",
			Value:    session,
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
		}
		http.SetCookie(w, &c)
		http.Redirect(w, r, "/home", http.StatusFound)
	}
}

func TickHandler(g *game.GameService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userID := inject.GetUserID(ctx)
		if userID == nil {
			http.Redirect(w, r, "/start", http.StatusFound)
			return
		}

		err := g.ProcessTicks(userID.(int))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func TradeHandler(g *game.GameService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		r, err := tx.Exec(context.Background(), "UPDATE stones SET amount = amount - $1 WHERE material = $2 AND owner_id = $3 AND amount >= $1;", trade.Amount, trade.Material, trade.OwnerID)
		if err != nil {
			fmt.Printf("Error updating trade owner stone amount: %s\n", err.Error())
			return
		} else if r.RowsAffected() == 0 {
			fmt.Println("Not enough stones: owner")
			return
		}

		_, err = tx.Exec(context.Background(), "UPDATE stones SET amount = amount + $1 WHERE material = $2 AND owner_id = $3;", trade.Amount, trade.Material, userID)
		if err != nil {
			fmt.Printf("Error updating trade responder stone amount: %s\n", err.Error())
			return
		}

		r, err = tx.Exec(context.Background(), "UPDATE stones SET amount = amount - $1 WHERE material = $2 AND owner_id = $3 AND amount >= $1;", trade.AmountReq, trade.MaterialReq, userID)
		if err != nil {
			fmt.Printf("Error updating trade responder stone amount 2: %s\n", err.Error())
			return
		} else if r.RowsAffected() == 0 {
			fmt.Println("Not enough stones: responder")
			return
		}

		_, err = tx.Exec(context.Background(), "UPDATE stones SET amount = amount + $1 WHERE material = $2 AND owner_id = $3;", trade.AmountReq, trade.MaterialReq, trade.OwnerID)
		if err != nil {
			fmt.Printf("Error updating trade owner stone amount 2: %s\n", err.Error())
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

func StaticHandler(w http.ResponseWriter, r *http.Request) {
	data, err := os.ReadFile(r.URL.Path[1:])
	if err != nil {
		fmt.Printf("Error loading static files: %s\n", err.Error())
	}

	w.Header().Set("Content-Type", "text/css")

	w.Write(data)
}
