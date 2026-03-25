// Package handlers defines http handlers
package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/underark/stone-collector/internal/models/state"
	"github.com/underark/stone-collector/internal/models/stones"
)

func GetHandler(storage []stones.Stone) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			storage = append(storage, stones.New())
			s := state.State{Stones: len(storage)}
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(s)
		}
	}
}
