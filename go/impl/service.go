package impl

import (
	"log/slog"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/micro"

	"xiam.li/protonats/go/protonats"
)

func NewService(name string, conn *nats.Conn, impl any, opts ...protonats.ServerOption) (micro.Service, *ServerOpts,
	error) {
	config := micro.Config{
		Name:    name,
		Version: "1.0.0",
	}

	// Check if the service implements any of the handler interfaces
	// but do it before applying options, so these can still override the handlers
	if statsHandler, isStatsHandler := impl.(protonats.StatsHandler); isStatsHandler {
		config.StatsHandler = statsHandler.Stats
		slog.Debug("Service implements StatsHandler; using service's Stats method")
	}
	if doneHandler, isDoneHandler := impl.(protonats.DoneHandler); isDoneHandler {
		config.DoneHandler = doneHandler.Done
		slog.Debug("Service implements DoneHandler; using service's Done method")
	}
	if errHandler, isErrHandler := impl.(protonats.ErrHandler); isErrHandler {
		config.ErrorHandler = errHandler.Err
		slog.Debug("Service implements ErrHandler; using service's Err method")
	}

	// Apply options
	options := ProcessServerOptions(&config, opts...)

	// Create the service
	service, err := micro.AddService(conn, config)
	if err != nil {
		return nil, nil, err
	}

	return service, options, nil
}
