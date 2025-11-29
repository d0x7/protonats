package impl

import (
	"time"

	"xiam.li/protonats/go/protonats"
)

type CallOpts struct {
	InstanceID      string
	Timeout         time.Duration
	DisableFinisher bool
	ExtraSubject    string
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

func (opts *CallOpts) Subject(subject string) string {
	return Subject(subject, opts.ExtraSubject, opts.InstanceID)
}

// Interface guard
var _ protonats.CallOptions = (*CallOpts)(nil)
