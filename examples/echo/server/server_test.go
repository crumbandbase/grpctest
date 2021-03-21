package server_test

import (
	"context"
	"testing"

	"github.com/crumbandbase/grpctest"
	echopb "github.com/crumbandbase/grpctest/echo/proto"
	"github.com/crumbandbase/grpctest/echo/server"
)

func setupRecoveryServer(t *testing.T) (grpctest.Closer, echopb.EchoClient) {
	s := grpctest.NewServer()

	conn, err := s.ClientConn()
	if err != nil {
		t.Fatal(err)
	}

	echopb.RegisterEchoServer(s, &server.EchoServer{})
	s.Serve()

	return s.Close, echopb.NewEchoClient(conn)
}

func TestEcho(t *testing.T) {
	t.Run("returns the same message sent to the server", func(t *testing.T) {
		closer, client := setupRecoveryServer(t)
		defer closer()

		message := "Hello, world"
		resp, err := client.Echo(context.Background(), &echopb.EchoRequest{Message: message})
		if err != nil {
			t.Errorf("expected: nil, got: %v", err)
		}

		if resp.Message != message {
			t.Errorf("expected: %s, got: %s", message, resp.Message)
		}
	})
}
