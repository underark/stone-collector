package inject

import (
	"context"
	"net/http"

	"github.com/underark/stone-collector/internal/models"
)

type k string

func GetUserID(ctx context.Context) any {
	return ctx.Value(k("userID"))
}

func UserID(r *http.Request, id int) *http.Request {
	ctx := context.WithValue(r.Context(), k("userID"), id)
	return r.WithContext(ctx)
}

func State(r *http.Request, state models.State) *http.Request {
	ctx := context.WithValue(r.Context(), k("state"), state)
	return r.WithContext(ctx)
}

func GetState(ctx context.Context) (models.State, bool) {
	val := ctx.Value(k("state"))
	state, ok := val.(models.State)
	return state, ok
}
