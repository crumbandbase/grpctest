package main

import (
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"

	echopb "github.com/crumbandbase/grpctest/echo/proto"
	"github.com/crumbandbase/grpctest/echo/server"
)

func main() {
	s := grpc.NewServer()

	echopb.RegisterEchoServer(s, &server.EchoServer{})
	healthpb.RegisterHealthServer(s, &health.Server{})

	// Register reflection service on grpc server for debugging.
	reflection.Register(s)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Println("server started on [::]:50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
