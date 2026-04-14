// Package handlers defines http handlers
package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"os"

	"github.com/gorilla/sessions"
	"github.com/underark/stone-collector/internal/models"
	"github.com/underark/stone-collector/internal/service/game"
	"github.com/underark/stone-collector/web/inject"
)

type data struct {
	route []string
	state models.State
}

func HomeHandler(g *game.GameService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := inject.GetUserID(r.Context())
		if userID == nil {
			http.Redirect(w, r, "/start", http.StatusFound)
			return
		}

		state, err := g.GetUserState(userID.(int))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		fmt.Println(state)

		t, err := template.ParseFiles("./web/templates/base.tmpl", "./web/templates/index.tmpl")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		t.ExecuteTemplate(w, "base", state)
	}
}

func InventoryHandler(g *game.GameService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := inject.GetUserID(r.Context())
		if userID == nil {
			http.Redirect(w, r, "/start", http.StatusFound)
			return
		}

		state, err := g.GetUserState(userID.(int))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		t, err := template.ParseFiles("./web/templates/base.tmpl", "./web/templates/inventory.tmpl")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		t.ExecuteTemplate(w, "base", state)
	}
}

func StartHandler(g *game.GameService, c *sessions.CookieStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := inject.GetUserID(r.Context())
		if userID != nil {
			http.Redirect(w, r, "/home", http.StatusFound)
			return
		}

		sessionID, err := g.InsertNewUser()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		session, err := c.Get(r, "stone-collector")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		session.Values["user_id"] = sessionID
		session.Save(r, w)

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
		userID := inject.GetUserID(r.Context())
		if userID == nil {
			http.Redirect(w, r, "/start", http.StatusFound)
			return
		}

		tradeID := r.FormValue("id")
		fmt.Printf("Trade id is: %s\n", tradeID)
		err := g.ExecuteTrade(userID.(int), tradeID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/home", http.StatusFound)
	}
}

func TradeMenuHandler(g *game.GameService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := inject.GetUserID(r.Context())
		if userID == nil {
			http.Redirect(w, r, "/start", http.StatusFound)
			return
		}

		state, err := g.GetUserState(userID.(int))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		trades, err := g.GetTrades()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		state.Trades = trades

		t, err := template.ParseFiles("./web/templates/base.tmpl", "./web/templates/trades.tmpl")
		if err != nil {
			fmt.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		t.ExecuteTemplate(w, "base", state)
	}
}

func StaticHandler(w http.ResponseWriter, r *http.Request) {
	data, err := os.ReadFile(r.URL.Path[1:])
	fmt.Println(r.URL.Path[1:])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/css")

	w.Write(data)
}
