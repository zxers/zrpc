package message

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncodeDecodeReq(t *testing.T) {

	tests := []struct {
		name string
		req *Request
		want []byte
	}{
		// TODO: Add test cases.
		{
			name:"test1",
			req: &Request{
				MessageId: 12,
				ServiceName: string([]byte{57,58}),
				MethodName: string([]byte{59,60}),
				Meta: map[string]string{
					"12": "sas",
				},
				Data: []byte{12, 2},
			},
		},

	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.req.SetHeadLength()
			tt.req.BodyLength = uint32(len(tt.req.Data))
			data := EncodeReq(tt.req)
			fmt.Println(data)
			res := DecodeReq(data)
			fmt.Println(res)
			assert.Equal(t, tt.req, res)
		})
	}
}

func TestEncodeDecodeResp(t *testing.T) {

	tests := []struct {
		name string
		resp *Response
		want []byte
	}{
		// TODO: Add test cases.
		{
			name:"test1",
			resp: &Response{
				MessageId: 12,
				Error: []byte{},
				Data: []byte{12, 2},
			},
		},

	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.resp.SetHeadLength()
			tt.resp.BodyLength = uint32(len(tt.resp.Data))
			data := EncodeResp(tt.resp)
			fmt.Println(data)
			res := DecodeResp(data)
			fmt.Println(res)
			assert.Equal(t, tt.resp, res)
		})
	}
}