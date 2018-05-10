package grpcutil

import (
	"context"
	"fmt"
	"google.golang.org/grpc/stats"
	"google.golang.org/grpc/status"
	"time"
)

// BasicStatsHandler is a stats handler that logs the method and the
// time it took to run
type BasicStatsHandler struct{}

type contextKey string

func (c contextKey) String() string {
	return string(c)
}

var (
	contextKeyMethod contextKey = "dangersalad/grpcutil:method"
	contextKeyStart  contextKey = "dangersalad/grpcutil:start"
)

// TagRPC tags the incoming request
func (h *BasicStatsHandler) TagRPC(ctx context.Context, tagInfo *stats.RPCTagInfo) context.Context {
	ctx = context.WithValue(ctx, contextKeyMethod, tagInfo.FullMethodName)
	ctx = context.WithValue(ctx, contextKeyStart, time.Now())
	return ctx
}

// HandleRPC runs after the request has finished
func (h *BasicStatsHandler) HandleRPC(ctx context.Context, i stats.RPCStats) {

	switch info := i.(type) {
	case *stats.End:
		if start, ok := ctx.Value(contextKeyStart).(time.Time); ok {

			diff := info.EndTime.Sub(start)
			diffStr := diff.String()
			if diff > time.Second {
				diffStr = diff.Truncate(time.Millisecond).String()
			} else if diff > time.Millisecond {
				diffStr = fmt.Sprintf("%0.3fms", float64(diff.Nanoseconds())/10000000.0)
			}

			method := ctx.Value(contextKeyMethod)

			if info.Error != nil {
				if gErr, ok := status.FromError(info.Error); ok {
					logf("%s - %s (%s)", method, gErr.Message(), diffStr)
				} else {
					logf("%s - %s (%s)", method, info.Error.Error(), diffStr)
				}
			} else {
				logf("%s (%s)", method, diffStr)
			}
		}
	}
}

// TagConn tags the connection
func (h *BasicStatsHandler) TagConn(ctx context.Context, tagInfo *stats.ConnTagInfo) context.Context {
	return ctx
}

// HandleConn runs after the connection
func (h *BasicStatsHandler) HandleConn(ctx context.Context, i stats.ConnStats) {
}
