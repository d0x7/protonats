package impl_test

import (
	"testing"
	"time"

	"xiam.li/protonats/go/impl"
	"xiam.li/protonats/go/protonats"
)

func TestCallOpts_WithInstanceID(t *testing.T) {
	opts := new(impl.CallOpts)
	protonats.WithInstanceID("test")(opts)
	if opts.InstanceID != "test" {
		t.Error("InstanceID not set correctly")
	}
}

func TestCallOpts_WithTimeout(t *testing.T) {
	opts := new(impl.CallOpts)
	protonats.WithTimeout(100 * time.Millisecond)(opts)
	if opts.Timeout != 100*time.Millisecond {
		t.Error("Timeout not set correctly")
	}
}

func TestCallOpts_WithExtraSubject(t *testing.T) {
	opts := new(impl.CallOpts)
	protonats.WithExtraSubject("test")(opts)
	if opts.ExtraSubject != "test" {
		t.Error("ExtraSubject not set correctly")
	}
}

func TestProcessCallOptions(t *testing.T) {
	opts := impl.ProcessCallOptions(
		protonats.WithInstanceID("test"),
		protonats.WithTimeout(100*time.Millisecond),
	)
	if opts.InstanceID != "test" {
		t.Error("InstanceID not set correctly")
	}
	if opts.Timeout != 100*time.Millisecond {
		t.Error("Timeout not set correctly")
	}
}
