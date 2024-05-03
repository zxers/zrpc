package zrpc

import (
	"context"
	"errors"
	"net"
	"test/zrpc/message"
	"test/zrpc/serialize/json"
	"time"

	"github.com/silenceper/pool"
)

type Client struct {
	connPool pool.Pool

}

func NewClient(address string, opts ...ClientOption) (*Client, error) {
	// if address == "///" {
	// 	address = solve(address)
	// }
	poolConfig := &pool.Config{
		InitialCap: 10,
		MaxCap: 100,
		MaxIdle: 50,
		Factory: func() (interface{}, error) {
			return net.Dial("tcp", address)
		},
		IdleTimeout: time.Minute,
		Close: func(i interface{}) error {
			return i.(net.Conn).Close()
		},
	}
	connPool, err := pool.NewChannelPool(poolConfig)

	if err != nil {
		return nil, err
	}

	client := &Client{
		connPool: connPool,
	}
	for _, opt := range opts {
		opt(client)
	}
	return client, nil
}

func (c *Client) InitService(service Service) error {
	return InitClientProxy(service, c, &json.Serializer{})
}

func (c *Client) Invoke(ctx context.Context, req *message.Request) (*message.Response, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	done := make(chan struct{})
	var (
		resp *message.Response
		err error
	)
	go func ()  {
		resp, err = c.doInvoke(ctx, req)	
		close(done)
	}()
	select {
	case <-done:
		if err != nil {
			return nil, err
		}
		return resp, nil
	case <-ctx.Done():
		return nil, errors.New("invoke timeout")
	}
}

func (c *Client) doInvoke(ctx context.Context, req *message.Request) (*message.Response, error) {
	// 发送请求
	conn, err := c.connPool.Get()
	if err != nil {
		return nil, err
	}
	netConn := conn.(net.Conn)
	data := message.EncodeReq(req)
	_, err = netConn.Write(data)
	if err != nil {
		return nil, err
	}

	// 读取响应
	data, err = ReadMsg(netConn)
	if err != nil {
		return nil, err
	}
	resp := message.DecodeResp(data)
	return resp, nil
}