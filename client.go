package grpcutil

import (
	"github.com/dangersalad/go-environment"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
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

// GetSecureConnection will return a new gRPC connection with the
// credentials specified in SSL_CRT_FILE
func GetSecureConnection(host, serverOverride string) (*grpc.ClientConn, error) {
	conf, err := environment.ReadOptions(environment.Options{
		EnvKeyCrtFile: "/ssl/tls.crt",
	})
	if err != nil {
		return nil, errors.Wrap(err, "reading environment")
	}
	crt := conf[EnvKeyCrtFile]

	creds, err := credentials.NewClientTLSFromFile(crt, serverOverride)
	if err != nil {
		return nil, errors.Wrap(err, "creating credentials")
	}

	// Create a connection with the TLS credentials
	conn, err := grpc.Dial(host, grpc.WithTransportCredentials(creds))
	if err != nil {
		return nil, errors.Wrap(err, "dialing service")
	}

	return conn, nil
}
