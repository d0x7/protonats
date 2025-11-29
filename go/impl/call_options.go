package impl

import (
	"time"

	"github.com/nats-io/nats.go"
	"xiam.li/protonats/go/protonats"
)

type CallOpts struct {
	InstanceID      string
	Timeout         time.Duration
	DisableFinisher bool
	ExtraSubject    string

	responseHeaderPointers     []*nats.Header
	responseInstanceIDPointers []*string
}

func (opts *CallOpts) WithoutFinisher() {
	opts.DisableFinisher = true
}

func ProcessCallOptions(opts ...protonats.CallOption) *CallOpts {
	options := new(CallOpts)
	for _, opt := range opts {
		opt(options)
	}
	return options
}

func (opts *CallOpts) SetInstanceID(id string) {
	opts.InstanceID = id
}

func (opts *CallOpts) SetTimeout(timeout time.Duration) {
	opts.Timeout = timeout
}

func (opts *CallOpts) HasInstanceID() bool {
	return opts.InstanceID != ""
}

func (opts *CallOpts) HasTimeout() bool {
	return opts.Timeout != 0
}

func (opts *CallOpts) GetTimeoutOr(duration time.Duration) time.Duration {
	if opts.HasTimeout() {
		return opts.Timeout
	}
	return duration
}

func (opts *CallOpts) SetExtraSubject(subject string) {
	opts.ExtraSubject = subject
}

func (opts *CallOpts) AddResponseHeaderPointer(target *nats.Header) {
	if target != nil {
		opts.responseHeaderPointers = append(opts.responseHeaderPointers, target)
	}
}

func (opts *CallOpts) AddResponseInstanceIDPointer(target *string) {
	if target != nil {
		opts.responseInstanceIDPointers = append(opts.responseInstanceIDPointers, target)
	}
}

func (opts *CallOpts) Subject(subject string) string {
	return Subject(subject, opts.ExtraSubject, opts.InstanceID)
}

func (opts *CallOpts) ApplyResponseHeaders(h nats.Header) {
	if h == nil {
		return
	}

	if len(opts.responseHeaderPointers) > 0 {
		for _, target := range opts.responseHeaderPointers {
			// Copy to avoid mutation issues
			*target = make(nats.Header, len(h))
			for k, v := range h {
				valCopy := make([]string, len(v))
				copy(valCopy, v)
				(*target)[k] = valCopy
			}
		}
	}

	if len(opts.responseInstanceIDPointers) > 0 {
		if id := h.Get(ServiceIDHeader); id != "" {
			for _, target := range opts.responseInstanceIDPointers {
				*target = id
			}
		}
	}
}

// Interface guard
var _ protonats.CallOptions = (*CallOpts)(nil)
