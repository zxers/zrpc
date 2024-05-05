package zrpc

import (
	"encoding/binary"
	"net"
	"zrpc/message"
)


func ReadMsg(conn net.Conn) ([]byte, error) {
	// 读取消息头长度
	header := make([]byte, message.MsgHeaderLength)
	_, err := conn.Read(header)
	if err != nil {
		return nil, err
	}
	headerLen := binary.BigEndian.Uint32(header)

	// 读取消息体长度
	body := make([]byte, message.MsgBodyLength)
	_, err = conn.Read(body)
	if err != nil {
		return nil, err
	}
	bodyLen := binary.BigEndian.Uint32(body)

	// 将消息头长度与消息体长度信息放入结果中，并将剩余的信息读取出来
	data := make([]byte, headerLen+bodyLen)
	copy(data, header)
	copy(data[message.MsgHeaderLength:], body)
	_, err = conn.Read(data[message.MsgHeaderLength+message.MsgBodyLength:])
	if err != nil {
		return nil, err
	}
	return data, nil
}