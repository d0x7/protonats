package protonats

import (
	"time"

	"github.com/nats-io/nats.go"
)

// CallOptions to be used when calling a method.
type CallOptions interface {
	SetInstanceID(string)
	SetTimeout(time.Duration)
	WithoutFinisher()
	SetExtraSubject(subject string)

	AddResponseHeaderPointer(target *nats.Header)
	AddResponseInstanceIDPointer(target *string)
}

type CallOption func(options CallOptions)

// WithInstanceID routes the call to the specified instance.
func WithInstanceID(id string) CallOption {
	return func(options CallOptions) {
		options.SetInstanceID(id)
	}
}

// WithTimeout overrides the default timeout for the call.
func WithTimeout(timeout time.Duration) CallOption {
	return func(options CallOptions) {
		options.SetTimeout(timeout)
	}
}

// WithoutFinisher disables the use of a finisher, which would return early
// after the first response when no further response is received after 250 ms
func WithoutFinisher() CallOption {
	return func(options CallOptions) {
		options.WithoutFinisher()
	}
}

// WithExtraSubject sets an extra subject for the subjects used by the microservice.
// Primarily used for consensus algorithms to distinguish between different services, while using the same implementation.
// Deprecated: previously used for consensus integration, shouldn't be used and will be removed in future versions.
func WithExtraSubject(extraSubject string) CallOption {
	return func(options CallOptions) {
		options.SetExtraSubject(extraSubject)
	}
}

// ResponseHeader sets a pointer to receive the response headers.
func ResponseHeader(target *nats.Header) CallOption {
	return func(options CallOptions) {
		options.AddResponseHeaderPointer(target)
	}
}

// ResponseInstanceID sets a pointer to receive the instance ID of the responding service instance.
func ResponseInstanceID(target *string) CallOption {
	return func(options CallOptions) {
		options.AddResponseInstanceIDPointer(target)
	}
}
