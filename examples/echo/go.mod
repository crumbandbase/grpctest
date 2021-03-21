module github.com/crumbandbase/grpctest/echo

go 1.16

require (
	github.com/crumbandbase/grpctest v0.0.0
	github.com/golang/protobuf v1.4.2
	google.golang.org/grpc v1.36.0
	google.golang.org/protobuf v1.25.0
)

replace github.com/crumbandbase/grpctest => ../../
