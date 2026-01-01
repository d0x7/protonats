package otel

import (
	"github.com/nats-io/nats.go"
	"go.opentelemetry.io/otel/propagation"
)

// natsHeaderCarrier adapts nats.Header to satisfy otel.TextMapCarrier
type natsHeaderCarrier struct {
	header nats.Header
}

var _ propagation.TextMapCarrier = (*natsHeaderCarrier)(nil)

func (c natsHeaderCarrier) Get(key string) string {
	return c.header.Get(key)
}

func (c natsHeaderCarrier) Set(key string, value string) {
	c.header.Set(key, value)
}

func (c natsHeaderCarrier) Keys() []string {
	keys := make([]string, 0, len(c.header))
	for k := range c.header {
		keys = append(keys, k)
	}
	return keys
}
