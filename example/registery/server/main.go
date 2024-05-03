package main

import (
	"test/zrpc"
	"test/zrpc/serialize/json"
	"test/zrpc/serialize/proto"
)

func main() {
	server := zrpc.NewServer()
	server.RegisterService(&UserServiceServer{})
	server.RegisterSerializer(&json.Serializer{})
	server.RegisterSerializer(&proto.Serializer{})
	server.Start(":8088")
}