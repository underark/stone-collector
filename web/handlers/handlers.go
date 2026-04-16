// Package handlers defines http handlers
package handlers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/sessions"
	"github.com/underark/stone-collector/internal/service/game"
	"github.com/underark/stone-collector/web/inject"
)

func HomeHandler(g *game.GameService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		state, ok := inject.GetState(r.Context())
		if !ok {
			http.Redirect(w, r, "/start", http.StatusFound)
			return
		}

		t, err := template.ParseFiles("./web/templates/base.tmpl", "./web/templates/index.tmpl")
		if err != nil {
			log.Printf("Error parsing /home template: %s\n", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		t.ExecuteTemplate(w, "base", state)
	}
}

func InventoryHandler(g *game.GameService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		state, ok := inject.GetState(r.Context())
		if !ok {
			http.Redirect(w, r, "/start", http.StatusFound)
			return
		}

		t, err := template.ParseFiles("./web/templates/base.tmpl", "./web/templates/inventory.tmpl")
		if err != nil {
			log.Printf("Error parsing /inventory template: %s\n", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		t.ExecuteTemplate(w, "base", state)
	}
}

func StartHandler(g *game.GameService, c *sessions.CookieStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, ok := inject.GetState(r.Context())
		if ok {
			http.Redirect(w, r, "/home", http.StatusFound)
			return
		}

		session, err := c.Get(r, "stone-collector")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		sessionID, err := g.InsertNewUser(session.Options.MaxAge)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		session.Values["session_id"] = sessionID
		session.Save(r, w)

		http.Redirect(w, r, "/home", http.StatusFound)
	}
}

func TickHandler(g *game.GameService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		state, ok := inject.GetState(r.Context())
		if !ok {
			http.Redirect(w, r, "/start", http.StatusFound)
			return
		}

		err := g.ProcessTicks(state.ID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func CreateTradeMenuHandler(g *game.GameService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		state, ok := inject.GetState(r.Context())
		if !ok {
			http.Redirect(w, r, "/start", http.StatusFound)
			return
		}

		t, err := template.ParseFiles("./web/templates/base.tmpl", "./web/templates/create.tmpl")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = t.ExecuteTemplate(w, "base", state)
		if err != nil {
			fmt.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func CreateTradeHandler(g *game.GameService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		state, ok := inject.GetState(r.Context())
		if !ok {
			http.Redirect(w, r, "/start", http.StatusFound)
			return
		}

		offerMaterial := r.FormValue("material")
		offerAmount := r.FormValue("amount")
		offerAmountInt, err := strconv.Atoi(offerAmount)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		requestMaterial := r.FormValue("material-req")
		requestAmount := r.FormValue("amount-req")
		requestAmountInt, err := strconv.Atoi(requestAmount)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = g.CreateTrade(state.ID, offerMaterial, offerAmountInt, requestMaterial, requestAmountInt)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/trade", http.StatusFound)
	}
}

func TradeHandler(g *game.GameService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		state, ok := inject.GetState(r.Context())
		if !ok {
			http.Redirect(w, r, "/start", http.StatusFound)
			return
		}

		formVal := r.FormValue("id")
		id, err := strconv.Atoi(formVal)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = g.ExecuteTrade(state.ID, id)
		if err != nil {
			http.Redirect(w, r, "/trade", http.StatusFound)
			return
		}

		http.Redirect(w, r, "/home", http.StatusFound)
	}
}

func TradeMenuHandler(g *game.GameService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		state, ok := inject.GetState(r.Context())
		if !ok {
			http.Redirect(w, r, "/start", http.StatusFound)
			return
		}

		trades, err := g.GetTrades(state.ID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		state.Trades = trades

		t, err := template.ParseFiles("./web/templates/base.tmpl", "./web/templates/trades.tmpl")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		t.ExecuteTemplate(w, "base", state)
	}
}

func StaticHandler(w http.ResponseWriter, r *http.Request) {
	data, err := os.ReadFile(r.URL.Path[1:])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/css")

	w.Write(data)
}
