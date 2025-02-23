package impl

import (
	"context"
	"testing"

	"github.com/nats-io/nats.go"

	"xiam.li/protonats/go/protonats"
)

func TestClientOpts_WithClientContext(t *testing.T) {
	opts := new(ClientOpts)
	protonats.WithClientContext(context.TODO())(opts)
	if opts.ClientContext != context.TODO() {
		t.Error("ClientContext not set correctly")
	}
}

func TestClientOpts_WithClientUnaryInterceptorChain(t *testing.T) {
	opts := new(ClientOpts)
	dummyClientUnaryInterceptor := func(next protonats.ClientUnaryMiddlewareHandler) protonats.ClientUnaryMiddlewareHandler {
		return func(ctx context.Context, request *nats.Msg) (*nats.Msg, error) {
			return next(ctx, request)
		}
	}
	protonats.WithClientUnaryInterceptorChain(dummyClientUnaryInterceptor)(opts)
	if opts.UnaryMiddlewareChain == nil {
		t.Error("UnaryMiddlewareChain not set correctly")
	}
}

func TestProcessClientOptions(t *testing.T) {
	t.Run("WithMicroConfig", func(t *testing.T) {
		dummyClientUnaryInterceptor := func(next protonats.ClientUnaryMiddlewareHandler) protonats.ClientUnaryMiddlewareHandler {
			return func(ctx context.Context, request *nats.Msg) (*nats.Msg, error) {
				return next(ctx, request)
			}
		}

		opts := ProcessClientOptions(
			protonats.WithClientContext(context.TODO()),
			protonats.WithClientUnaryInterceptorChain(dummyClientUnaryInterceptor),
		)
		if opts.ClientContext != context.TODO() {
			t.Error("ClientContext not set correctly")
		}
		if opts.UnaryMiddlewareChain == nil {
			t.Error("UnaryMiddlewareChain not set correctly")
		}
	})
}
