package recovery

import (
	"context"
	"log/slog"

	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"
	"xiam.li/protonats/go/protonats"
)

// NewClientInterceptor recovers from panics in subsequent interceptors or the handler on the client side.
// It should usually be the first interceptor in the chain.
func NewClientInterceptor(logger *slog.Logger) protonats.UnaryClientInterceptor {
	return func(ctx context.Context, info *protonats.MethodInfo, req, reply proto.Message, header nats.Header, invoker protonats.UnaryInvoker, opts ...protonats.CallOption) (err error) {
		defer func() {
			if r := recover(); r != nil {
				if logger != nil {
					logger.Error(
						"Panic recovered in ProtoNATS client handler",
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
		return invoker(ctx, info, req, reply, header, opts...)
	}
}
