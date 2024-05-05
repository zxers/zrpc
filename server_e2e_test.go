package zrpc

import (
	"context"
	"errors"
	"fmt"
	"zrpc/serialize/json"
	"zrpc/serialize/proto"
	"testing"
)

func TestServer(t *testing.T) {
	server := NewServer()
	server.RegisterService(&UserServiceServer{})
	server.RegisterSerializer(&json.Serializer{})
	server.RegisterSerializer(&proto.Serializer{})
	server.Start(":8088")
}

type UserServiceServer struct {
	
}

func (u *UserServiceServer) Name() string {
	return "UserService-Client"
}

func (u *UserServiceServer) GetById(ctx context.Context, req *AnyRequest) (*AnyResponse, error) {
	fmt.Println("GetById, 执行了响应")
	return &AnyResponse{
		Msg: "server的响应",
	}, nil
}

func (u *UserServiceServer) GetByIdError(ctx context.Context, req *AnyRequest) (*AnyResponse, error) {
	fmt.Println("GetById, 执行了响应")
	return &AnyResponse{
		Msg: "server的响应",
	}, errors.New("测试error")
}