package cookies

import "context"

type k string

func GetUserID(ctx context.Context) any {
	return ctx.Value(k("userID"))
}
