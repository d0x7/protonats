package protonats

import (
	"context"

	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"
)

type MethodInfo struct {
	Subject string
	Service string
	Method  string
	Type    MethodType
}

//go:generate go run github.com/dmarkham/enumer -type=MethodType -trimprefix=MethodType -output=methodtype_enum.go
type MethodType uint8

const (
	MethodTypeUnary MethodType = iota + 1
	MethodTypeStream
	MethodTypeBroadcast
	MethodTypeLeader
	MethodTypeFollower
)

func (m *MethodInfo) FullMethod() string {
	return m.Service + "." + m.Method
}

type UnaryInvoker func(ctx context.Context, info *MethodInfo, req, reply proto.Message, header nats.Header, opts ...CallOption) error

// UnaryClientInterceptor is a function that wraps the execution of a ClientInvoker.
//
// IMPORTANT - Header Handling:
// Do not use protonats.WithOutgoingHeader(ctx) (and every other adjacent method from
// the context.go file) inside an interceptor to set headers.
//
// The context headers are processed before the chain is invoked, and any changes
// made to the context inside the interceptor will be ignored.
//
// Instead, modify the 'headers' map directly:
//
//	// Correct:
//	headers.Set("My-Key", "value")
//
//	// Incorrect (Will be ignored):
//	ctx = protonats.WithOutgoingHeader(ctx, "My-Key", "value")
type UnaryClientInterceptor func(ctx context.Context, info *MethodInfo, req, reply proto.Message, header nats.Header, invoker UnaryInvoker, opts ...CallOption) error

type UnaryHandler func(ctx context.Context, req proto.Message) (proto.Message, error)

// UnaryServerInterceptor acts as a middleware for incoming server requests.
//
// The provided context 'ctx' is pre-populated with incoming NATS headers.
// To read headers (e.g. for Auth or Tracing), use:
//
//	headers := protonats.HeadersFromContext(ctx)
//
// To modify the context passed to the handler (e.g. adding a User object
// after auth), use:
//
//	newCtx := context.WithValue(ctx, "user", myUser)
//	return next(newCtx, req)
//
// If an interceptor returns an error, the execution chain stops, and the
// error is returned to the client (converted to a generic 500 or passed
// through if it is a protonats.ServerError).
//
// It is highly recommended to use the Recovery interceptor (middleware/recovery/ServerInterceptor) as the first
// interceptor in your chain to prevent panics from crashing the whole service.
type UnaryServerInterceptor func(ctx context.Context, req proto.Message, info *MethodInfo, handler UnaryHandler) (proto.Message, error)
