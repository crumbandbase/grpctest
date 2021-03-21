package grpctest

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/test/bufconn"

	"github.com/crumbandbase/grpctest/internal"
)

const bufferSize = 1 << 20

type bufDialer func(context.Context, string) (net.Conn, error)

// Closer is a convenience type to be used in test suites.
type Closer func()

// Server is a gRPC server listening on a buffered stream, for use in
// end-to-end gRPC tests.
type Server struct {
	*grpc.Server

	Listener    *bufconn.Listener
	TLS         *tls.Config
	certificate *x509.Certificate
}

// NewServer returns a new insecure server. The caller should start the server
// by calling Serve on the returned type, and call Close when finished, to shut
// it down.
func NewServer(opts ...grpc.ServerOption) *Server {
	return &Server{
		Listener: bufconn.Listen(bufferSize),
		Server:   grpc.NewServer(opts...),
	}
}

// NewTLSServer returns a new secure server. The caller should start the
// server by calling Serve on the returned type, and call Close when finished,
// to shut it down.
func NewTLSServer(opts ...grpc.ServerOption) *Server {
	cert, err := tls.X509KeyPair(internal.LocalhostCert, internal.LocalhostKey)
	if err != nil {
		panic(fmt.Sprintf("grpctest: NewSecureServer: %v", err))
	}

	certificate, err := x509.ParseCertificate(cert.Certificate[0])
	if err != nil {
		panic(fmt.Sprintf("grpctest: NewSecureServer: %v", err))
	}

	config := &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.NoClientCert,
	}

	tlsCredentials := credentials.NewTLS(config)

	return &Server{
		Listener:    bufconn.Listen(bufferSize),
		TLS:         config,
		Server:      grpc.NewServer(append(opts, grpc.Creds(tlsCredentials))...),
		certificate: certificate,
	}
}

// Start starts the server.
func (s *Server) Serve() {
	go func() { s.Server.Serve(s.Listener) }()
}

// Close shuts down the server and blocks until all outstanding requests on
// this server have completed.
func (s *Server) Close() {
	s.Listener.Close()
	s.Server.GracefulStop()
}

// Certificate returns the certificate used by the server, or nil if the server
// doesn't use TLS.
func (s *Server) Certificate() *x509.Certificate {
	return s.certificate
}

// ClientConn returns a gRPC client connection configured for making requests
// to the server. It is configured to trust the server's TLS test certificate,
// if present, and will close its idle connections on Close.
func (s *Server) ClientConn(opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	return s.ClientConnContext(context.Background(), opts...)
}

// ClientConnContext returns a gRPC client connection configured for making
// requests to the server. It is configured to trust the server's TLS test
// certificate, if present, and will close its idle connections on Close.
func (s *Server) ClientConnContext(ctx context.Context, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	opts = append(opts, grpc.WithContextDialer(getBufdialer(s.Listener)))

	if s.certificate != nil {
		certPool := x509.NewCertPool()
		certPool.AddCert(s.certificate)

		tlsCredentials := credentials.NewTLS(&tls.Config{
			RootCAs: certPool,
		})

		opts = append(opts, grpc.WithTransportCredentials(tlsCredentials))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}

	return grpc.DialContext(ctx, "bufconn", opts...)
}

func getBufdialer(listener *bufconn.Listener) bufDialer {
	return func(context.Context, string) (net.Conn, error) {
		return listener.Dial()
	}
}
