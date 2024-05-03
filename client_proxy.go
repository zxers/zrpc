package zrpc

import (
	"context"
	"errors"
	"log"
	"reflect"
	"strconv"
	"test/zrpc/message"
	"test/zrpc/serialize"
)

func InitClientProxy(service Service, p Proxy, serialize serialize.Serializer) error {
	if reflect.TypeOf(service).Kind() != reflect.Ptr {
		return errors.New("类型错误")
	}
	val := reflect.ValueOf(service).Elem()
	typ := reflect.TypeOf(service).Elem()
	numField := typ.NumField()
	for i := 0; i < numField; i++ {
		fieldTyp := typ.Field(i)
		fieldVal := val.Field(i)
		if !fieldVal.CanSet() {
			continue
		}
		makeVal := reflect.MakeFunc(fieldTyp.Type, func(args []reflect.Value) (results []reflect.Value) {
			ctx, ok := args[0].Interface().(context.Context)
			meta := make(map[string]string)		
			
			// 单向调用	
			oneWay := isOneWay(ctx)
			if oneWay {
				meta["oneWay"] = "true"
			}
			if !ok {
				panic("ctx param error")
			}

			// 超时控制
			deadline, ok := ctx.Deadline()
			if ok {
				meta["deadline"] = strconv.FormatInt(deadline.UnixMilli(), 10)
			}

			data, err := serialize.Encode(args[1].Interface())
			if err != nil {
				return
			}
			req := &message.Request{
				ServiceName: service.Name(),
				MethodName: fieldTyp.Name,
				Data: data,
				Serializer: serialize.Id(),
				Meta: meta,
			}
			req.SetHeadLength()
			req.BodyLength = uint32(len(req.Data)) 
			log.Println(req)
			// 通过网络发送请求
			resp, err := p.Invoke(ctx, req)
			log.Println(resp)
			log.Println(string(resp.Error))
			log.Println(string(resp.Data))
			// 解析响应数据
			respData := reflect.New(fieldTyp.Type.Out(0))
			err = serialize.Decode(resp.Data, respData.Interface())
			if err != nil {
				return
			}
			respData = respData.Elem()
			results = append(results, respData)
			// 解析响应error
			results = append(results, reflect.ValueOf(errors.New(string(resp.Error))))
			return 
		})
		fieldVal.Set(makeVal)
	}
	return nil
}

type Service interface {
	Name() string
}