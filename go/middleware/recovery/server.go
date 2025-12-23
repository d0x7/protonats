package recovery

import (
	"context"
	"log/slog"

	"google.golang.org/protobuf/proto"
	"xiam.li/protonats/go/protonats"
)

// NewServerInterceptor recovers from panics in subsequent interceptors or the handler on the server side.
// It should usually be the first interceptor in the chain.
func NewServerInterceptor(logger *slog.Logger) protonats.UnaryServerInterceptor {
	return func(ctx context.Context, req proto.Message, info *protonats.MethodInfo, next protonats.UnaryHandler) (resp proto.Message, err error) {
		defer func() {
			if r := recover(); r != nil {
				if logger != nil {
					logger.Error(
						"Panic recovered in ProtoNATS server handler",
						slog.Any("panic", r),
						slog.GroupAttrs("route",
							slog.String("service", info.Service),
							slog.String("method", info.Method),
							slog.String("subject", info.Subject),
							slog.String("type", info.Type.String()),
						),
					)
				}

				// Return a safe error to the client
				err = protonats.ErrPanicRecovered.AsServerError()
			}
		}()
		return next(ctx, req)
	}
}
