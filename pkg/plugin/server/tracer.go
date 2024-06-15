package server

import "context"

type key int

const (
	reqidKey key = 0
)

func GetReqidFromContext(ctx context.Context) string {
	// todo ? int to string?
	ret, _ := ctx.Value(reqidKey).(string)
	return ret
}
