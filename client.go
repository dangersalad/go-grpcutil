package grpcutil

import (
	"fmt"

	"github.com/dangersalad/go-environment"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// GetInternalConnection will return a new gRPC connection with no
// credentials
func GetInternalConnection(host string, opt ...grpc.DialOption) (*grpc.ClientConn, error) {
	opts := append(opt, grpc.WithInsecure())
	conn, err := grpc.Dial(host, opts...)
	if err != nil {
		return nil, fmt.Errorf("setting up grpc client: %w", err)
	}
	return conn, nil
}

// GetSecureConnection will return a new gRPC connection with the
// credentials specified in SSL_CRT_FILE
func GetSecureConnection(host, serverOverride string, opt ...grpc.DialOption) (*grpc.ClientConn, error) {
	conf, err := environment.ReadOptions(environment.Options{
		EnvKeyCrtFile: "/ssl/tls.crt",
	})
	if err != nil {
		return nil, fmt.Errorf("reading environment: %w", err)
	}
	crt := conf[EnvKeyCrtFile]

	creds, err := credentials.NewClientTLSFromFile(crt, serverOverride)
	if err != nil {
		return nil, fmt.Errorf("creating credentials: %w", err)
	}

	opts := append(opt, grpc.WithTransportCredentials(creds))
	// Create a connection with the TLS credentials
	conn, err := grpc.Dial(host, opts...)
	if err != nil {
		return nil, fmt.Errorf("dialing service: %w", err)
	}

	return conn, nil
}
