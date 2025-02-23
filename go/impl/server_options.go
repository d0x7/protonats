package impl

import (
	"context"
	"time"

	"github.com/nats-io/nats.go/micro"

	"xiam.li/protonats/go/protonats"
)

type ServerOpts struct {
	InstanceID               string
	Timeout                  time.Duration
	WithoutLeaderFunctions   bool
	WithoutFollowerFunctions bool
	ServiceVersion           string
	ExtraSubject             string
	StatsHandler             *micro.StatsHandler
	DoneHandler              *micro.DoneHandler
	ErrorHandler             *micro.ErrHandler
	ServerContext            context.Context
	UnaryMiddlewareChain     []protonats.ServerUnaryMiddleware
}

func ProcessServerOptions(config *micro.Config, opts ...protonats.ServerOption) *ServerOpts {
	options := new(ServerOpts)

	// Set defaults
	options.ServerContext = context.Background()

	for _, opt := range opts {
		opt(options)
	}
	if options.ServiceVersion != "" {
		config.Version = options.ServiceVersion
	}
	if options.StatsHandler != nil {
		config.StatsHandler = *options.StatsHandler
	}
	if options.DoneHandler != nil {
		config.DoneHandler = *options.DoneHandler
	}
	if options.ErrorHandler != nil {
		config.ErrorHandler = *options.ErrorHandler
	}
	return options
}

func (opts *ServerOpts) SetStatsHandler(handler micro.StatsHandler) {
	opts.StatsHandler = &handler
}

func (opts *ServerOpts) SetDoneHandler(handler micro.DoneHandler) {
	opts.DoneHandler = &handler
}

func (opts *ServerOpts) SetErrorHandler(handler micro.ErrHandler) {
	opts.ErrorHandler = &handler
}

func (opts *ServerOpts) SetServiceVersion(serviceVersion string) {
	opts.ServiceVersion = serviceVersion
}

func (opts *ServerOpts) WithoutLeaderFns() {
	opts.WithoutLeaderFunctions = true
}

func (opts *ServerOpts) WithoutFollowerFns() {
	opts.WithoutFollowerFunctions = true
}

func (opts *ServerOpts) SetServerContext(ctx context.Context) {
	opts.ServerContext = ctx
}

func (opts *ServerOpts) SetUnaryMiddlewareChain(unaryMiddlewares ...protonats.ServerUnaryMiddleware) {
	opts.UnaryMiddlewareChain = unaryMiddlewares
}

func (opts *ServerOpts) SetExtraSubject(extraSubject string) {
	opts.ExtraSubject = extraSubject
}

func (opts *ServerOpts) Subject(subject, suffix string) micro.EndpointOpt {
	return micro.WithEndpointSubject(_subject(subject, opts.ExtraSubject, suffix))
}

func _subject(subject, extra, suffix string) string {
	if extra != "" {
		subject += "." + extra
	}
	if suffix != "" {
		subject += "." + suffix
	}
	return subject
}

// Interface guard
var _ protonats.ServerOptions = (*ServerOpts)(nil)

func ApplyServerUnaryMiddlewares(handler protonats.ServerUnaryMiddlewareHandler,
	middlewares ...protonats.ServerUnaryMiddleware) protonats.ServerUnaryMiddlewareHandler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}
	return handler
}
