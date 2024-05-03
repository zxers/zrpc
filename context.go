package zrpc

import "context"

func CtxWithOneWay(ctx context.Context) context.Context {
	return context.WithValue(ctx, "oneWay", "true")
}

func isOneWay(ctx context.Context) bool {
	return ctx.Value("oneWay") != nil 
}