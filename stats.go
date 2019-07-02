package grpcutil

import (
	"context"
	"fmt"
	"google.golang.org/grpc/stats"
	"google.golang.org/grpc/status"
	"regexp"
	"time"
)

// BasicStatsHandler is a stats handler that logs the method and the
// time it took to run
type BasicStatsHandler struct {
	logBypass *regexp.Regexp
}

func NewStatsHandler(bypass *regexp.Regexp) *BasicStatsHandler {
	return &BasicStatsHandler{
		logBypass: bypass,
	}
}

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

	ctxMethod := ctx.Value(contextKeyMethod)

	method, ok := ctxMethod.(string)
	if !ok {
		logf("method %s is not a string", ctxMethod)
	}

	if h.logBypass.MatchString(method) {
		return
	}

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
