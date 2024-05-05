package main

import (
	"context"
	"zrpc"
)



func main() {
	
	client, err := zrpc.NewClient(":8088")
	if err != nil {
		return
	}
	usc := &UserServiceClient{}
	err = client.InitService(usc)
	if err != nil {
		return
	}
	usc.GetById(context.Background(), &AnyRequest{Msg: "zxzx"})
	// 单向调用
	usc.GetById(zrpc.CtxWithOneWay(context.Background()), &AnyRequest{Msg: "zxzx"})
}

type UserServiceClient struct {
	GetById func(ctx context.Context, req *AnyRequest) (*AnyResponse, error)
}

func (u *UserServiceClient) Name() string {
	return "UserService-Client"
}

type AnyRequest struct {
	Msg string `json:"msg"`
}

type AnyResponse struct {
	Msg string `json:"msg"`
}
