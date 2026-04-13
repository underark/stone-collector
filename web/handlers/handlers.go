// Package handlers defines http handlers
package handlers

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"os"

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
		userID := inject.GetUserID(r.Context())
		if userID == nil {
			http.Redirect(w, r, "/start", http.StatusFound)
			return
		}

		tradeID := r.URL.Query().Get("tradeID")
		err := g.ExecuteTrade(userID.(int), tradeID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/home", http.StatusFound)
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
