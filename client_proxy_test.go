package zrpc

import (
	"context"
	"zrpc/message"
	"zrpc/serialize/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitClientProxy(t *testing.T) {
	
	tests := []struct {
		name    string
		service *UserServiceClient
		wantErr error
	}{
		// TODO: Add test cases.
		{
			name: "UserServiceClient",
			service: &UserServiceClient{},
			wantErr: nil,
		},
	}
	mockProxy := &MockProxy{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := InitClientProxy(tt.service, mockProxy, &json.Serializer{});
			assert.Equal(t, tt.wantErr, err)
			_, _ = tt.service.GetById(context.Background(), &AnyRequest{})
		})
	}
}

type UserServiceClient struct {
	GetById func(ctx context.Context, req *AnyRequest) (*AnyResponse, error)
}

func (u *UserServiceClient) Name() string {
	return "UserService-Client"
}

type AnyRequest struct {
	Msg string `json:"msg"`
}

type AnyResponse struct {
	Msg string `json:"msg"`
}

type MockProxy struct {

}

func (m *MockProxy) Invoke(context.Context, *message.Request) (*message.Response, error) {
	return nil, nil
}

