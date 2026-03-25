// Package handlers defines http handlers
package handlers

import (
	"io"
	"net/http"

	"github.com/underark/stone-collector/internal/models"
)

func GetHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		stone := stones.New()
		io.WriteString(w, stone.Material)
	}
}
