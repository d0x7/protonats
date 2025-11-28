package impl

import (
	"context"

	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"
	"xiam.li/protonats/go/protonats"
)

type ClientOpts struct {
	UnaryInterceptor protonats.UnaryClientInterceptor

	chainUnaryInterceptors []protonats.UnaryClientInterceptor
}

func ProcessClientOptions(opts ...protonats.ClientOption) *ClientOpts {
	options := new(ClientOpts)
	for _, opt := range opts {
		opt(options)
	}

	// Chain unary interceptors

	var chainedInterceptors protonats.UnaryClientInterceptor
	if len(options.chainUnaryInterceptors) == 0 {
		chainedInterceptors = nil
	} else if len(options.chainUnaryInterceptors) == 1 {
		chainedInterceptors = options.chainUnaryInterceptors[0]
	} else {
		chainedInterceptors = func(ctx context.Context, info *protonats.MethodInfo, req, reply proto.Message, headers nats.Header, invoker protonats.UnaryInvoker, opts ...protonats.CallOption) error {
			return options.chainUnaryInterceptors[0](ctx, info, req, reply, headers, getChainUnaryInvoker(options.chainUnaryInterceptors, 0, invoker), opts...)
		}
	}
	options.UnaryInterceptor = chainedInterceptors

	return options
}

func (opts *ClientOpts) AddUnaryInterceptor(interceptor protonats.UnaryClientInterceptor) {
	opts.chainUnaryInterceptors = append(opts.chainUnaryInterceptors, interceptor)
}

func (opts *ClientOpts) ChainUnaryInterceptors(interceptors ...protonats.UnaryClientInterceptor) {
	opts.chainUnaryInterceptors = append(opts.chainUnaryInterceptors, interceptors...)
}

// Interface guard
var _ protonats.ClientOptions = (*ClientOpts)(nil)
