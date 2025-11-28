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
	BaseContext              context.Context
	StatsHandler             *micro.StatsHandler
	DoneHandler              *micro.DoneHandler
	ErrorHandler             *micro.ErrHandler
	UnaryInterceptor         protonats.UnaryServerInterceptor

	chainUnaryInterceptors []protonats.UnaryServerInterceptor
}

func ProcessServerOptions(config *micro.Config, opts ...protonats.ServerOption) *ServerOpts {
	options := new(ServerOpts)
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

	// Chain unary interceptors

	var chainedInterceptors protonats.UnaryServerInterceptor
	if len(options.chainUnaryInterceptors) == 0 {
		chainedInterceptors = nil
	} else if len(options.chainUnaryInterceptors) == 1 {
		chainedInterceptors = options.chainUnaryInterceptors[0]
	} else {
		chainedInterceptors = chainUnaryInterceptors(options.chainUnaryInterceptors)
	}
	options.UnaryInterceptor = chainedInterceptors

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

func (opts *ServerOpts) SetExtraSubject(extraSubject string) {
	opts.ExtraSubject = extraSubject
}

func (opts *ServerOpts) SetBaseContext(ctx context.Context) {
	opts.BaseContext = ctx
}

func (opts *ServerOpts) AddUnaryInterceptor(interceptor protonats.UnaryServerInterceptor) {
	opts.chainUnaryInterceptors = append(opts.chainUnaryInterceptors, interceptor)
}

func (opts *ServerOpts) ChainUnaryInterceptors(interceptors ...protonats.UnaryServerInterceptor) {
	opts.chainUnaryInterceptors = append(opts.chainUnaryInterceptors, interceptors...)
}

func (opts *ServerOpts) Subject(subject, suffix string) micro.EndpointOpt {
	return micro.WithEndpointSubject(_subject(subject, opts.ExtraSubject, suffix))
}

func (opts *ServerOpts) Ctx() context.Context {
	if opts.BaseContext != nil {
		return opts.BaseContext
	}
	return context.Background()
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
