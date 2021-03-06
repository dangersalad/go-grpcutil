// Package grpcutil provides helper functions for boilerplate grpc setup
package grpcutil // import "github.com/dangersalad/go-grpcutil"

import (
	"fmt"
	"net"
	"strings"

	env "github.com/dangersalad/go-environment"
	"google.golang.org/grpc"
)

// EnvKeySecureServer is the environment key to trigger making a
// secure server. Any non empty value will do.
const EnvKeySecureServer = "GRPC_SECURE"

// ServerSet is the gRPC server and it's listener together
type ServerSet struct {
	server   *grpc.Server
	listener net.Listener
	secure   bool
}

// Server returns the ServerSet's gRPC server
func (s *ServerSet) Server() *grpc.Server {
	return s.server
}

// Listener returns the ServerSet's listener
func (s *ServerSet) Listener() net.Listener {
	return s.listener
}

// IsSecure returns true if the server uses SSL
func (s *ServerSet) IsSecure() bool {
	return s.secure
}

// Serve starts the server
func (s *ServerSet) Serve() error {
	return fmt.Errorf("serving grpc: %w", s.server.Serve(s.listener))
}

// IsSecure returns the value of the --secure flag or the env var GRPC_SECURE
func IsSecure() bool {
	conf := env.ReadOptionsAllowMissing(env.Options{
		EnvKeySecureServer: "",
	})
	return conf[EnvKeySecureServer] != ""
}

// BaseServerOpts returns the base set of grpc server options as an array
func BaseServerOpts() []grpc.ServerOption {
	return []grpc.ServerOption{grpc.StatsHandler(&BasicStatsHandler{})}
}

// MakeServerOpts returns the default options if the provided options are empty
func MakeServerOpts(opts ...grpc.ServerOption) []grpc.ServerOption {
	if len(opts) == 0 {
		return BaseServerOpts()
	}
	return opts
}

// CreateServer will return a new gRPC server, either secured or not based on the presence of the --secure flag
func CreateServer(port string, opt ...grpc.ServerOption) (*ServerSet, error) {
	if strings.Index(port, ":") < 0 {
		port = ":" + port
	}
	opt = MakeServerOpts(opt...)
	if IsSecure() {
		debug("creating secured server")
		return CreateSecureServer(port, opt...)
	}
	debug("creating internal server")
	return CreateInternalServer(port, opt...)
}

// CreateSecureServer will return a new gRPC server using the
// credentials in the container
func CreateSecureServer(port string, opt ...grpc.ServerOption) (*ServerSet, error) {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		return nil, fmt.Errorf("setting up port: %w", err)
	}
	creds, err := GetServerCredentials()
	if err != nil {
		return nil, err
	}
	opt = append(opt, grpc.Creds(creds))
	return &ServerSet{
		server:   grpc.NewServer(opt...),
		listener: lis,
		secure:   true,
	}, nil
}

// CreateInternalServer will return a new gRPC server with no
// authentication
func CreateInternalServer(port string, opt ...grpc.ServerOption) (*ServerSet, error) {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		return nil, fmt.Errorf("setting up port: %w", err)
	}
	return &ServerSet{
		server:   grpc.NewServer(opt...),
		listener: lis,
		secure:   false,
	}, nil
}
