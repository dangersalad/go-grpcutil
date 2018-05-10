package grpcutil

import (
	"context"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

type checkerFunc func(context.Context) (context.Context, error)

type serverStreamWrapper struct {
	grpc.ServerStream

	ctx context.Context
}

func (w *serverStreamWrapper) Context() context.Context { return w.ctx }

// StreamChecker returns a gRPC server option that run a checker
// function on incoming stream's contexts. Takes a
// func(context.Context) (context.Context, error).
func StreamChecker(check checkerFunc) grpc.ServerOption {
	checker := grpc.StreamServerInterceptor(func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		ctx, err := check(ss.Context())
		if err != nil {
			return errors.Wrap(err, "checker failed")
		}
		return handler(srv, &serverStreamWrapper{ServerStream: ss, ctx: ctx})
	})
	return grpc.StreamInterceptor(checker)
}

// UnaryChecker returns a gRPC server option that will run a checker
// function on incoming unary request's contexts. Takes a
// func(context.Context) (context.Context, error).
func UnaryChecker(check checkerFunc) grpc.ServerOption {
	checker := grpc.UnaryServerInterceptor(func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		ctx, err := check(ctx)
		if err != nil {
			return nil, errors.Wrap(err, "checker failed")
		}
		return handler(ctx, req)
	})
	return grpc.UnaryInterceptor(checker)
}
