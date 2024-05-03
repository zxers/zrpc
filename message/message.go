package message

import (
	"encoding/binary"
	"strings"
)

const (
	spliter = '\n'

	MsgHeaderLength = 4
	MsgBodyLength = 4
)

type Request struct {
	// 头部字段
	HeadLength  uint32
	BodyLength  uint32
	MessageId   uint32
	Version     uint8
	Compresser  uint8
	Serializer  uint8
	ServiceName string
	MethodName  string

	Meta map[string]string
	// 协议体
	Data []byte
}

func (r *Request) SetHeadLength() {
	r.HeadLength += 15
	r.HeadLength += uint32(len(r.ServiceName) + 1 + len(r.MethodName) + 1)
	for k, v := range r.Meta {
		r.HeadLength += uint32(len(k) + 1 + len(v) + 1)
	}
}

func EncodeReq(req *Request) []byte {
	res := make([]byte, req.HeadLength+req.BodyLength)

	cur := res
	binary.BigEndian.PutUint32(cur, req.HeadLength)
	cur = cur[4:]

	binary.BigEndian.PutUint32(cur, req.BodyLength)
	cur = cur[4:]

	binary.BigEndian.PutUint32(cur, req.MessageId)
	cur = cur[4:]

	cur[0] = req.Version
	cur = cur[1:]

	cur[0] = req.Compresser
	cur = cur[1:]

	cur[0] = req.Serializer
	cur = cur[1:]

	copy(cur, []byte(req.ServiceName))
	cur = cur[len(req.ServiceName):]
	cur[0] = spliter
	cur = cur[1:]

	copy(cur, []byte(req.MethodName))
	cur = cur[len(req.MethodName):]
	cur[0] = spliter
	cur = cur[1:]	

	for k, v := range req.Meta {
		copy(cur, []byte(k))
		cur = cur[len(k):]
		cur[0] = spliter
		cur = cur[1:]
		
		copy(cur, []byte(v))
		cur = cur[len(v):]
		cur[0] = spliter
		cur = cur[1:]
	}
	copy(cur, req.Data)
	return res
}

func DecodeReq(data []byte) *Request {
	req := &Request{
		Meta: map[string]string{},
		
	}
	req.HeadLength = binary.BigEndian.Uint32(data[:4])

	req.BodyLength = binary.BigEndian.Uint32(data[4:8])

	req.MessageId = binary.BigEndian.Uint32(data[8:12])

	req.Version = data[12]

	req.Compresser = data[13]

	req.Serializer = data[14]

	meta := data[15:req.HeadLength]
	index := strings.IndexByte(string(meta), byte(spliter))
	req.ServiceName = string(meta[:index])
	meta = meta[index+1:]

	index = strings.IndexByte(string(meta), byte(spliter))
	req.MethodName = string(meta[:index])
	meta = meta[index+1:]

	for  {
		index = strings.IndexByte(string(meta), byte(spliter))
		if index == -1 {
			break
		}
		key := string(meta[:index])
		meta = meta[index+1:]

		index = strings.IndexByte(string(meta), byte(spliter))
		val := string(meta[:index])
		meta = meta[index+1:]

		req.Meta[key] = val
	}
	req.Data = data[req.HeadLength:]
	return req
}

type Response struct {
	// 头部字段
	HeadLength uint32
	BodyLength uint32
	MessageId  uint32
	Version    uint8
	Compresser uint8
	Serializer uint8

	Error []byte
	// 协议体
	Data []byte
}

func (resp *Response) SetHeadLength() {
	resp.HeadLength = uint32(15 + len(resp.Error))
}

// 这里处理 Resp 我直接复制粘贴，是因为我觉得复制粘贴会使可读性更高

func EncodeResp(resp *Response) []byte {
	bs := make([]byte, resp.HeadLength+resp.BodyLength)

	cur := bs
	// 1. 写入 HeadLength，四个字节
	binary.BigEndian.PutUint32(cur[:4], resp.HeadLength)
	cur = cur[4:]
	// 2. 写入 BodyLength 四个字节
	binary.BigEndian.PutUint32(cur[:4], resp.BodyLength)
	cur = cur[4:]

	// 3. 写入 message id, 四个字节
	binary.BigEndian.PutUint32(cur[:4], resp.MessageId)
	cur = cur[4:]

	// 4. 写入 version，因为本身就是一个字节，所以不用进行编码了
	cur[0] = resp.Version
	cur = cur[1:]

	// 5. 写入压缩算法
	cur[0] = resp.Compresser
	cur = cur[1:]

	// 6. 写入序列化协议
	cur[0] = resp.Serializer
	cur = cur[1:]
	// 7. 写入 error
	copy(cur, resp.Error)
	cur = cur[len(resp.Error):]

	// 剩下的数据
	copy(cur, resp.Data)
	return bs
}

// DecodeResp 解析 Response
func DecodeResp(bs []byte) *Response {
	resp := &Response{}
	// 按照 EncodeReq 写下来
	// 1. 读取 HeadLength
	resp.HeadLength = binary.BigEndian.Uint32(bs[:4])
	// 2. 读取 BodyLength
	resp.BodyLength = binary.BigEndian.Uint32(bs[4:8])
	// 3. 读取 message id
	resp.MessageId = binary.BigEndian.Uint32(bs[8:12])
	// 4. 读取 Version
	resp.Version = bs[12]
	// 5. 读取压缩算法
	resp.Compresser = bs[13]
	// 6. 读取序列化协议
	resp.Serializer = bs[14]
	// 7. error 信息
	resp.Error = bs[15:resp.HeadLength]

	// 剩下的就是数据了
	resp.Data = bs[resp.HeadLength:]
	return resp
}