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
}

func (m *MethodInfo) FullMethod() string {
	return m.Service + "." + m.Method
}

type UnaryInvoker func(ctx context.Context, info *MethodInfo, req, reply proto.Message, header nats.Header, opts ...CallOption) error

type UnaryClientInterceptor func(ctx context.Context, info *MethodInfo, req, reply proto.Message, header nats.Header, invoker UnaryInvoker, opts ...CallOption) error

type UnaryHandler func(ctx context.Context, req proto.Message) (proto.Message, error)

type UnaryServerInterceptor func(ctx context.Context, req proto.Message, info *MethodInfo, handler UnaryHandler) (proto.Message, error)
