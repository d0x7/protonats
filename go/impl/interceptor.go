package impl

import (
	"context"

	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"
	"xiam.li/protonats/go/protonats"
)

// Copyright: https://github.com/grpc/grpc-go/blob/master/server.go / clientconn.go

func chainUnaryInterceptors(interceptors []protonats.UnaryServerInterceptor) protonats.UnaryServerInterceptor {
	return func(ctx context.Context, req proto.Message, info *protonats.MethodInfo, handler protonats.UnaryHandler) (proto.Message, error) {
		return interceptors[0](ctx, req, info, getChainUnaryHandler(interceptors, 0, info, handler))
	}
}

func getChainUnaryHandler(interceptors []protonats.UnaryServerInterceptor, curr int, info *protonats.MethodInfo, finalHandler protonats.UnaryHandler) protonats.UnaryHandler {
	if curr == len(interceptors)-1 {
		return finalHandler
	}
	return func(ctx context.Context, req proto.Message) (proto.Message, error) {
		return interceptors[curr+1](ctx, req, info, getChainUnaryHandler(interceptors, curr+1, info, finalHandler))
	}
}

func getChainUnaryInvoker(interceptors []protonats.UnaryClientInterceptor, curr int, finalInvoker protonats.UnaryInvoker) protonats.UnaryInvoker {
	if curr == len(interceptors)-1 {
		return finalInvoker
	}
	return func(ctx context.Context, info *protonats.MethodInfo, req, reply proto.Message, headers nats.Header, opts ...protonats.CallOption) error {
		return interceptors[curr+1](ctx, info, req, reply, headers, getChainUnaryInvoker(interceptors, curr+1, finalInvoker), opts...)
	}
}
