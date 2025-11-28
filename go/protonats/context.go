package protonats

import (
	"context"

	"github.com/nats-io/nats.go"
)

type HeaderKey struct{}
type OutgoingHeaderKey struct{}

// HeadersFromContext retrieves NATS headers from the context for use in server interceptors.
func HeadersFromContext(ctx context.Context) nats.Header {
	if h, ok := ctx.Value(HeaderKey{}).(nats.Header); ok {
		return h
	}
	return nil
}

// HeadersFromOutgoingContext retrieves NATS headers from the context for use in client interceptors.
func HeadersFromOutgoingContext(ctx context.Context) nats.Header {
	if h, ok := ctx.Value(OutgoingHeaderKey{}).(nats.Header); ok {
		return h
	}
	return nil
}

// WithOutgoingHeaders puts NATS headers into the context for use in client interceptors.
// It merges with existing headers if present.
func WithOutgoingHeaders(ctx context.Context, h nats.Header) context.Context {
	existing := HeadersFromOutgoingContext(ctx)
	newHeaders := make(nats.Header)

	if existing != nil {
		for k, v := range existing {
			for _, val := range v {
				newHeaders.Add(k, val)
			}
		}
	}
	for k, v := range h {
		for _, val := range v {
			newHeaders.Add(k, val)
		}
	}

	return context.WithValue(ctx, OutgoingHeaderKey{}, newHeaders)
}

// WithOutgoingHeader adds a single NATS header key-value pair into the context for use in client interceptors.
// It merges with existing headers if present.
func WithOutgoingHeader(ctx context.Context, key, value string) context.Context {
	existing := HeadersFromOutgoingContext(ctx)
	newHeaders := make(nats.Header)

	if existing != nil {
		for k, v := range existing {
			for _, val := range v {
				newHeaders.Add(k, val)
			}
		}
	}

	newHeaders.Add(key, value)
	return context.WithValue(ctx, OutgoingHeaderKey{}, newHeaders)
}

// WithOutgoingHeaderMap adds multiple NATS header key-value pairs into the context for use in client interceptors.
// It merges with existing headers if present.
func WithOutgoingHeaderMap(ctx context.Context, headers map[string]string) context.Context {
	existing := HeadersFromOutgoingContext(ctx)
	newHeaders := make(nats.Header)

	if existing != nil {
		for k, v := range existing {
			for _, val := range v {
				newHeaders.Add(k, val)
			}
		}
	}

	for k, v := range headers {
		newHeaders.Add(k, v)
	}

	return context.WithValue(ctx, OutgoingHeaderKey{}, newHeaders)
}

// SetOutgoingHeaders sets NATS headers into the context for use in client interceptors,
// replacing any existing headers.
func SetOutgoingHeaders(ctx context.Context, h nats.Header) context.Context {
	return context.WithValue(ctx, OutgoingHeaderKey{}, h)
}

// SetOutgoingHeaderMap sets multiple NATS header key-value pairs into the context for use in client interceptors,
// replacing any existing headers.
func SetOutgoingHeaderMap(ctx context.Context, headers map[string]string) context.Context {
	newHeaders := make(nats.Header)
	for k, v := range headers {
		newHeaders.Add(k, v)
	}
	return context.WithValue(ctx, OutgoingHeaderKey{}, newHeaders)
}
