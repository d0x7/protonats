package otel

import (
	"context"

	"github.com/nats-io/nats.go"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/semconv/v1.37.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/protobuf/proto"
	"xiam.li/protonats/go/protonats"
)

var ProtoNATSRPCKey = semconv.RPCServiceKey.String("protonats")

// NewTracingClientInterceptor returns an interceptor that traces outgoing requests
// and propagates the trace context via NATS headers.
func NewTracingClientInterceptor(opts ...Option) protonats.UnaryClientInterceptor {
	cfg := newConfig(opts...)
	tracer := cfg.tracerProvider.Tracer(cfg.tracerName)

	return func(ctx context.Context, info *protonats.MethodInfo, req, reply proto.Message, headers nats.Header, next protonats.UnaryInvoker, callOpts ...protonats.CallOption) error {
		// Using info.FullMethod() as span name
		ctx, span := tracer.Start(ctx, info.FullMethod(),
			trace.WithSpanKind(trace.SpanKindClient),
			trace.WithAttributes(
				ProtoNATSRPCKey,
				semconv.RPCServiceKey.String(info.Service),
				semconv.RPCMethodKey.String(info.Method),
			),
		)
		defer span.End()

		// Inject Trace Context into NATS Headers and then call next
		cfg.propagator.Inject(ctx, natsHeaderCarrier{header: headers})

		// Call the next invoker and set status
		if err := next(ctx, info, req, reply, headers, callOpts...); err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			return err
		} else {
			span.SetStatus(codes.Ok, "OK")
			return nil
		}
	}
}

// NewTracingServerInterceptor returns an interceptor that extracts trace context
// from NATS headers and manages the server-side span.
func NewTracingServerInterceptor(opts ...Option) protonats.UnaryServerInterceptor {
	cfg := newConfig(opts...)
	tracer := cfg.tracerProvider.Tracer(cfg.tracerName)

	return func(ctx context.Context, req proto.Message, info *protonats.MethodInfo, next protonats.UnaryHandler) (proto.Message, error) {
		// Extract Headers from Context
		headers := protonats.HeadersFromContext(ctx)
		if headers == nil {
			headers = make(nats.Header)
		}

		ctx = cfg.propagator.Extract(ctx, natsHeaderCarrier{header: headers})

		ctx, span := tracer.Start(ctx, info.FullMethod(),
			trace.WithSpanKind(trace.SpanKindServer),
			trace.WithAttributes(
				ProtoNATSRPCKey,
				semconv.RPCServiceKey.String(info.Service),
				semconv.RPCMethodKey.String(info.Method),
			),
		)
		defer span.End()

		// Call the next handler and set status
		if resp, err := next(ctx, req); err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			return nil, err
		} else {
			span.SetStatus(codes.Ok, "OK")
			return resp, err
		}
	}
}
