package protonats

import (
	"context"
	"time"

	"github.com/nats-io/nats.go/micro"
)

// ServerUnaryMiddlewareHandler ServerUnaryMiddleware handler function
type ServerUnaryMiddlewareHandler func(ctx context.Context, request micro.Request)

// ServerUnaryMiddleware used to intercept requests before they are processed by server handler function
type ServerUnaryMiddleware func(next ServerUnaryMiddlewareHandler) ServerUnaryMiddlewareHandler

// Those handler interfaces can be implemented by the respective NATS server impl,
//as an alternative to setting a handler via an option when creating the server.

// StatsHandler is an interface that when implemented on the server, will
// be used directly instead of having to use the WithStatsHandler option.
type StatsHandler interface {
	Stats(endpoint *micro.Endpoint) any
}

// DoneHandler is an interface that when implemented on the server, will
// be used directly instead of having to use the WithDoneHandler option.
type DoneHandler interface {
	Done(service micro.Service)
}

// ErrHandler is an interface that when implemented on the server, will
// be used directly instead of having to use the WithErrorHandler option.
type ErrHandler interface {
	Err(service micro.Service, natsErr *micro.NATSError)
}

type Ping struct {
	micro.ServiceIdentity
	Type string
	RTT  time.Duration
}

// ServerOptions to be used when creating a new server.
type ServerOptions interface {
	SetStatsHandler(micro.StatsHandler)
	SetDoneHandler(micro.DoneHandler)
	SetErrorHandler(micro.ErrHandler)
	SetServiceVersion(string)
	WithoutLeaderFns()
	WithoutFollowerFns()
	SetExtraSubject(string)
	SetServerContext(ctx context.Context)
	SetUnaryMiddlewareChain(unaryInterceptors ...ServerUnaryMiddleware)
}

type ServerOption func(options ServerOptions)

// WithStatsHandler sets the stats handler for the server.
// It's an alternative to implementing the StatsHandler interface on the server implementation.
func WithStatsHandler(handler micro.StatsHandler) ServerOption {
	return func(options ServerOptions) {
		options.SetStatsHandler(handler)
	}
}

// WithDoneHandler sets the done handler for the server.
// It's an alternative to implementing the DoneHandler interface on the server implementation.
func WithDoneHandler(handler micro.DoneHandler) ServerOption {
	return func(options ServerOptions) {
		options.SetDoneHandler(handler)
	}
}

// WithErrorHandler sets the error handler for the server.
// It's an alternative to implementing the ErrHandler interface on the server implementation.
func WithErrorHandler(handler micro.ErrHandler) ServerOption {
	return func(options ServerOptions) {
		options.SetErrorHandler(handler)
	}
}

// WithServiceVersion sets the service version for the NATS microservice.
// By default, the service version is set to "1.0.0".
func WithServiceVersion(version string) ServerOption {
	return func(options ServerOptions) {
		options.SetServiceVersion(version)
	}
}

// WithoutLeaderFns disables registration of functions marked as LEADER via protonats.consensus_target.
// Only has an effect when used with the NewMyServiceNATSServer function, and not the respective Leader/Follower functions.
func WithoutLeaderFns() ServerOption {
	return func(options ServerOptions) {
		options.WithoutLeaderFns()
	}
}

// WithoutFollowerFns disables registration of functions marked as FOLLOWER via protonats.consensus_target.
// Only has an effect when used with the NewMyServiceNATSServer function, and not the respective Leader/Follower functions.
func WithoutFollowerFns() ServerOption {
	return func(options ServerOptions) {
		options.WithoutFollowerFns()
	}
}

// WithExtraSubjectSrv sets an extra subject for the subjects used by the microservice.
// Primarily used for consensus algorithms to distinguish between different services, while using the same implementation.
func WithExtraSubjectSrv(extraSubject string) ServerOption {
	return func(options ServerOptions) {
		options.SetExtraSubject(extraSubject)
	}
}

// WithServerContext sets server context.
// Used as a base context for all interceptors and requests.
func WithServerContext(ctx context.Context) ServerOption {
	return func(options ServerOptions) {
		options.SetServerContext(ctx)
	}
}

// WithServerUnaryInterceptorChain sets a chain of unary interceptors.
// Used to process request before passed to server implementation.
func WithServerUnaryInterceptorChain(unaryInterceptorsChain ...ServerUnaryMiddleware) ServerOption {
	return func(options ServerOptions) {
		options.SetUnaryMiddlewareChain(unaryInterceptorsChain...)
	}
}
