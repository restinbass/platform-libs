package grpc_interceptors

import (
	"context"
	"path"
	"time"

	"github.com/restinbass/platform-libs/pkg/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

// LoggerInterceptor -
func LoggerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		method := path.Base(info.FullMethod)
		logger.Info(ctx, "started gRPC method", zap.String("method", method))

		startTime := time.Now()
		resp, err := handler(ctx, req)
		duration := time.Since(startTime)

		if err != nil {
			st, _ := status.FromError(err)
			logger.Info(
				ctx,
				"finished gRPC method",
				zap.String("method", method),
				zap.String("code", st.Code().String()),
				zap.Error(err),
				zap.Duration("time_taken", duration),
			)
		} else {
			logger.Info(
				ctx,
				"successfully finished gRPC method",
				zap.String("method", method),
				zap.Duration("time_taken", duration),
			)
		}

		return resp, err
	}
}
