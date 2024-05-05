package zrpc

import (
	"context"
	"zrpc/message"
)

type Proxy interface {
	Invoke(context.Context, *message.Request) (*message.Response, error)
}

