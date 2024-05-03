package zrpc

import (
	"context"
	"test/zrpc/message"
)

type Proxy interface {
	Invoke(context.Context, *message.Request) (*message.Response, error)
}

