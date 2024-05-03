package zrpc

import (
	"context"
	// "fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient(t *testing.T) {
	client, err := NewClient(":8088")
	assert.NoError(t, err)
	usc := &UserServiceClient{}
	err = client.InitService(usc)
	if err != nil {
		t.Error(err)
	}
	resp, err := usc.GetById(context.Background(), &AnyRequest{Msg: "zxzx"})
	t.Log(resp)
	t.Log(err)
	// 单向调用
	usc.GetById(CtxWithOneWay(context.Background()), &AnyRequest{Msg: "zxzx"})
}