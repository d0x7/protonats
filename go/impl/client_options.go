package impl

import (
	"context"

	"xiam.li/protonats/go/protonats"
)

type ClientOpts struct {
	ClientContext        context.Context
	UnaryMiddlewareChain []protonats.ClientUnaryMiddleware
}

func ProcessClientOptions(opts ...protonats.ClientOption) *ClientOpts {
	options := new(ClientOpts)

	// Set defaults
	options.ClientContext = context.Background()

	for _, opt := range opts {
		opt(options)
	}
	return options
}

func (opts *ClientOpts) SetClientContext(ctx context.Context) {
	opts.ClientContext = ctx
}

func (opts *ClientOpts) SetUnaryMiddlewareChain(unaryMiddlewares ...protonats.ClientUnaryMiddleware) {
	opts.UnaryMiddlewareChain = unaryMiddlewares
}

// Interface guard
var _ protonats.ClientOptions = (*ClientOpts)(nil)
