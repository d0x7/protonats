package slog

import (
	"context"
	"log/slog"
	"time"

	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"
	"xiam.li/protonats/go/protonats"
)

// NewClientInterceptor logs the method, duration, and any errors for outgoing requests.
// If the provided logger is not enabled for debug level, only errors are logged.
func NewClientInterceptor(logger *slog.Logger) protonats.UnaryClientInterceptor {
	return func(ctx context.Context, info *protonats.MethodInfo, req, reply proto.Message, headers nats.Header, next protonats.UnaryInvoker, opts ...protonats.CallOption) error {
		start := time.Now()

		err := next(ctx, info, req, reply, headers, opts...)

		duration := time.Since(start)

		if err != nil || logger.Enabled(ctx, slog.LevelDebug) {
			route := slog.GroupAttrs("route",
				slog.String("service", info.Service),
				slog.String("method", info.Method),
				slog.String("subject", info.Subject),
				slog.String("type", info.Type.String()),
			)
			durationAttr := slog.Duration("duration", duration)

			if err != nil {
				logger.Error("ProtoNATS client call error", durationAttr, route)
			} else {
				logger.Debug("ProtoNATS client call success", durationAttr, route)
			}
		}

		return err
	}
}

// NewServerInterceptor logs the method, duration, and any errors for incoming requests.
// If the provided logger is not enabled for debug level, only errors are logged.
func NewServerInterceptor(logger *slog.Logger) protonats.UnaryServerInterceptor {
	return func(ctx context.Context, req proto.Message, info *protonats.MethodInfo, next protonats.UnaryHandler) (proto.Message, error) {
		start := time.Now()

		resp, err := next(ctx, req)

		duration := time.Since(start)

		if err != nil || logger.Enabled(ctx, slog.LevelDebug) {
			route := slog.GroupAttrs("route",
				slog.String("service", info.Service),
				slog.String("method", info.Method),
				slog.String("subject", info.Subject),
				slog.String("type", info.Type.String()),
			)
			durationAttr := slog.Duration("duration", duration)

			if err != nil {
				logger.Error("ProtoNATS service handler error", durationAttr, route)
			} else {
				logger.Debug("ProtoNATS service handler success", durationAttr, route)
			}
		}

		return resp, err
	}
}
