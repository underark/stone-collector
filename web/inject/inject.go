package inject

import (
	"context"
	"net/http"
)

type k string

func GetUserID(ctx context.Context) any {
	return ctx.Value(k("userID"))
}

func UserID(r *http.Request, id int) *http.Request {
	ctx := context.WithValue(r.Context(), k("userID"), id)
	return r.WithContext(ctx)
}
