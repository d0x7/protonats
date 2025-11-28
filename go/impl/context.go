package impl

import (
	"context"

	"github.com/nats-io/nats.go"
	"xiam.li/protonats/go/protonats"
)

// NewContextWithHeaders injects NATS headers into the context for use in server interceptors.
func NewContextWithHeaders(ctx context.Context, h nats.Header) context.Context {
	return context.WithValue(ctx, protonats.HeaderKey{}, h)
}
