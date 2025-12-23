package chaos

import (
	"context"
	"errors"
	"math/rand"
	"time"

	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"
	"xiam.li/protonats/go/protonats"
)

type Config struct {
	ErrorRate  float64 // 0.0 to 1.0 (e.g. 0.1 = 10% errors)
	MinLatency time.Duration
	MaxLatency time.Duration
}

// NewClientInterceptor creates a client interceptor that injects faults.
func NewClientInterceptor(cfg Config) protonats.UnaryClientInterceptor {
	return func(ctx context.Context, info *protonats.MethodInfo, req, reply proto.Message, headers nats.Header, next protonats.UnaryInvoker, opts ...protonats.CallOption) error {
		if cfg.MaxLatency > 0 {
			// Random duration between Min and Max
			delta := int64(cfg.MaxLatency - cfg.MinLatency)
			if delta > 0 {
				sleep := cfg.MinLatency + time.Duration(rand.Int63n(delta))
				time.Sleep(sleep)
			} else {
				time.Sleep(cfg.MinLatency)
			}
		}

		// Inject error based on ErrorRate
		if cfg.ErrorRate > 0.0 {
			if rand.Float64() < cfg.ErrorRate {
				return errors.New("simulated chaos monkey error")
			}
		}

		return next(ctx, info, req, reply, headers, opts...)
	}
}
