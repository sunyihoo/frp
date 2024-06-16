package server

import "context"

type key int

const (
	reqidKey key = 0
)

func NewReqidContext(ctx context.Context, reqid string) context.Context {
	return context.WithValue(ctx, reqidKey, reqid)
}

func GetReqidFromContext(ctx context.Context) string {
	// todo ? int to string?
	ret, _ := ctx.Value(reqidKey).(string)
	return ret
}
