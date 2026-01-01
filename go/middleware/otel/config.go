package otel

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

type config struct {
	tracerProvider trace.TracerProvider
	propagator     propagation.TextMapPropagator
	tracerName     string
}

type Option func(*config)

func WithTracerProvider(provider trace.TracerProvider) Option {
	return func(c *config) {
		c.tracerProvider = provider
	}
}

func WithPropagator(propagator propagation.TextMapPropagator) Option {
	return func(c *config) {
		c.propagator = propagator
	}
}

func newConfig(opts ...Option) config {
	c := config{
		tracerName: "xiam.li/protonats", // Name of the instrumentation library
	}
	for _, opt := range opts {
		opt(&c)
	}
	// Fallback to globals if not provided
	if c.tracerProvider == nil {
		c.tracerProvider = otel.GetTracerProvider()
	}
	if c.propagator == nil {
		c.propagator = otel.GetTextMapPropagator()
	}
	return c
}
