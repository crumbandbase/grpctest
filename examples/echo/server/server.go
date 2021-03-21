package server

import (
	"context"

	echopb "github.com/crumbandbase/grpctest/echo/proto"
)

type EchoServer struct {
	echopb.UnimplementedEchoServer
}

// grpcurl -plaintext -d '{"message": "Hello, world"}' :50051 echo.v1.Echo/Echo | jq -r .message
func (s *EchoServer) Echo(ctx context.Context, in *echopb.EchoRequest) (*echopb.EchoResponse, error) {
	return &echopb.EchoResponse{Message: in.Message}, nil
}
