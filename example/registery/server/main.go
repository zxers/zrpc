package main

import (
	"zrpc"
	"zrpc/serialize/json"
	"zrpc/serialize/proto"
)

func main() {
	server := zrpc.NewServer()
	server.RegisterService(&UserServiceServer{})
	server.RegisterSerializer(&json.Serializer{})
	server.RegisterSerializer(&proto.Serializer{})
	server.Start(":8088")
}