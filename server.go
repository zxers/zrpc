package zrpc

import (
	"context"
	"log"
	"net"
	"reflect"
	"strconv"
	"zrpc/message"
	"zrpc/serialize"
	"zrpc/serialize/json"
	"time"
)

type Server struct {
	services map[string]Service
	serializers []serialize.Serializer
}

func NewServer() *Server {
	s := &Server{
		services: map[string]Service{},
		serializers: make([]serialize.Serializer, 256),
	}
	s.RegisterSerializer(&json.Serializer{})
	return s 
}

func (s *Server) RegisterService(service Service) {
	s.services[service.Name()] = service
}

func (s *Server) RegisterSerializer(serializer serialize.Serializer) {
	s.serializers[serializer.Id()] = serializer
}

func (s *Server) Start(addr string) error {
	listen, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Printf("accept connection error: %v", err)
		}
		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	for {
		// 读取请求
		data, err := ReadMsg(conn)
		if err != nil {
			return
		}
		
		req := message.DecodeReq(data)
		log.Println(req)

		// 本地执行服务调用
		service, ok := s.services[req.ServiceName]
		if !ok {
			return
		}

		method, ok := reflect.TypeOf(service).MethodByName(req.MethodName)
		if !ok {
			return
		}

		ctx := context.Background()
		var cancel context.CancelFunc
		for k, v := range req.Meta {
			if k == "deadline" {
				timeout, _ := strconv.ParseInt(v, 10, 64)
				ctx, cancel = context.WithDeadline(ctx, time.UnixMilli(timeout))
				defer cancel()
				continue
			}
			ctx = context.WithValue(ctx, k, v)
		}
		// 获取第零个方法调用参数，即结构体的指针本身
		param0 := reflect.New(method.Type.In(0))
		// 获取第一个方法调用参数，即ctx
		param1 := reflect.New(reflect.TypeOf(ctx))
		// 获取第二个方法调用参数，即req
		param2 := reflect.New(method.Type.In(2))
		
		if s.serializers[req.Serializer] == nil {
			return
		}
		err = s.serializers[req.Serializer].Decode(req.Data, &param2)
		if err != nil {
			return
		}
		
		res := method.Func.Call([]reflect.Value{param0.Elem(), param1.Elem(), param2.Elem()})
 
		if isOneWay(ctx) {
			continue
		}

		resData, err := s.serializers[req.Serializer].Encode(res[0].Interface())
		if err != nil {
			return
		}
		var errData []byte
		if !res[1].IsNil() {
			errData = []byte(res[1].Interface().(error).Error()) 
		}
	
		log.Println(string(errData))
		response := &message.Response{
			Error: errData,
			Data: resData,
		}
		response.SetHeadLength()
		response.BodyLength = uint32(len(response.Data))
		// 返回响应
		respData := message.EncodeResp(response)
		_, err = conn.Write(respData)
		if err != nil {
			return
		}
	}
}