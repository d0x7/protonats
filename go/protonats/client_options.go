package protonats

import (
	"context"

	"github.com/nats-io/nats.go"
)

// ClientOptions to be used when creating a new client.
type ClientOptions interface {
	SetClientContext(ctx context.Context)
	SetUnaryMiddlewareChain(unaryMiddlewares ...ClientUnaryMiddleware)
}

type ClientOption func(options ClientOptions)

// ClientUnaryMiddlewareHandler ClientUnaryMiddleware handler function
type ClientUnaryMiddlewareHandler func(ctx context.Context, request *nats.Msg) (*nats.Msg, error)

// ClientUnaryMiddleware used to intercept requests before they are processed by client handler function
type ClientUnaryMiddleware func(next ClientUnaryMiddlewareHandler) ClientUnaryMiddlewareHandler

// WithClientContext sets client base context.
// Used as a base context for all interceptors and requests.
func WithClientContext(ctx context.Context) ClientOption {
	return func(options ClientOptions) {
		options.SetClientContext(ctx)
	}
}

// WithClientUnaryInterceptorChain sets a chain of unary interceptors.
// Used to process request before passed to server implementation.
func WithClientUnaryInterceptorChain(unaryMiddlewares ...ClientUnaryMiddleware) ClientOption {
	return func(options ClientOptions) {
		options.SetUnaryMiddlewareChain(unaryMiddlewares...)
	}
}
