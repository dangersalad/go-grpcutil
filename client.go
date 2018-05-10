package grpcutil

import (
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

// GetInternalConnection will return a new gRPC connection with no
// credentials
func GetInternalConnection(host string) (*grpc.ClientConn, error) {
	conn, err := grpc.Dial(host, grpc.WithInsecure())
	if err != nil {
		return nil, errors.Wrap(err, "setting up grpc client")
	}
	return conn, nil
}
