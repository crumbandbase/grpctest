package grpctest_test

import (
	"testing"

	"github.com/crumbandbase/grpctest"
)

func TestNeServer(t *testing.T) {
	t.Run("returns a server with a buffered listener", func(t *testing.T) {
		if s := grpctest.NewServer(); s == nil {
			t.Error("expected: grpctest.Server, got: nil")
		}
	})
}

func TestNewTLSServer(t *testing.T) {
	t.Run("returns a TLS server with a buffered listener", func(t *testing.T) {
		if s := grpctest.NewTLSServer(); s == nil {
			t.Error("expected: grpctest.Server, got: nil")
		}
	})
}

func TestCertificate(t *testing.T) {
	t.Run("returns nil when the server does not use TLS", func(t *testing.T) {
		if s := grpctest.NewServer(); s.Certificate() != nil {
			t.Errorf("expected: nil, got: %v", s.Certificate())
		}
	})

	t.Run("returns a TLS certificate when the server uses TLS", func(t *testing.T) {
		if s := grpctest.NewTLSServer(); s == nil {
			t.Error("expected: x509.Certificate, got: nil")
		}
	})
}
