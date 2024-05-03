package main

import (
	"context"
	"errors"
	"fmt"
)

type UserServiceServer struct{}

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

type AnyRequest struct {
	Msg string `json:"msg"`
}

type AnyResponse struct {
	Msg string `json:"msg"`
}
