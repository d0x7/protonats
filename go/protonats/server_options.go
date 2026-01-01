package protonats

import (
	"context"
	"time"

	"github.com/nats-io/nats.go/micro"
)

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
	SetExtraSubject(string)
	AddUnaryInterceptor(UnaryServerInterceptor)
	ChainUnaryInterceptors(...UnaryServerInterceptor)
	SetBaseContext(context.Context)
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

// WithExtraSubjectSrv sets an extra subject for the subjects used by the microservice.
// Primarily used for consensus algorithms to distinguish between different services, while using the same implementation.
// Deprecated: previously used for consensus integration, shouldn't be used and will be removed in future versions.
func WithExtraSubjectSrv(extraSubject string) ServerOption {
	return func(options ServerOptions) {
		options.SetExtraSubject(extraSubject)
	}
}

// WithBaseContext sets the base context for the server, used as the parent context in all handlers.
func WithBaseContext(ctx context.Context) ServerOption {
	return func(options ServerOptions) {
		options.SetBaseContext(ctx)
	}
}

func UnaryInterceptor(interceptor UnaryServerInterceptor) ServerOption {
	return func(options ServerOptions) {
		options.AddUnaryInterceptor(interceptor)
	}
}

func ChainUnaryInterceptors(interceptors ...UnaryServerInterceptor) ServerOption {
	return func(options ServerOptions) {
		options.ChainUnaryInterceptors(interceptors...)
	}
}
