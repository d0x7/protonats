package impl_test

import (
	"context"
	"testing"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/micro"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"xiam.li/protonats/go/impl"
	"xiam.li/protonats/go/protonats"
)

type TestMsg struct{}

func (m *TestMsg) ProtoReflect() protoreflect.Message {
	return nil
}

func (m *TestMsg) Reset()         {}
func (m *TestMsg) String() string { return "" }
func (m *TestMsg) ProtoMessage()  {}

func TestClientInterceptor_Ordering(t *testing.T) {
	var executionLog []string

	// Interceptor 1 (Outer)
	int1 := func(ctx context.Context, info *protonats.MethodInfo, req, reply proto.Message, headers nats.Header, next protonats.UnaryInvoker, opts ...protonats.CallOption) error {
		executionLog = append(executionLog, "1-pre")
		err := next(ctx, info, req, reply, headers, opts...)
		executionLog = append(executionLog, "1-post")
		return err
	}

	// Interceptor 2 (Inner)
	int2 := func(ctx context.Context, info *protonats.MethodInfo, req, reply proto.Message, headers nats.Header, next protonats.UnaryInvoker, opts ...protonats.CallOption) error {
		executionLog = append(executionLog, "2-pre")
		err := next(ctx, info, req, reply, headers, opts...)
		executionLog = append(executionLog, "2-post")
		return err
	}

	// The Final Invoker (Network Call)
	finalInvoker := func(ctx context.Context, info *protonats.MethodInfo, req, reply proto.Message, headers nats.Header, opts ...protonats.CallOption) error {
		executionLog = append(executionLog, "invoker")
		return nil
	}

	// Build the chain
	clientOpts := protonats.WithChainUnaryInterceptor(int1, int2)
	options := impl.ProcessClientOptions(clientOpts)
	chain := options.UnaryInterceptor

	// Run
	headers := make(nats.Header)
	err := chain(context.Background(), &protonats.MethodInfo{}, &TestMsg{}, &TestMsg{}, headers, finalInvoker)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Verify Order
	expected := []string{"1-pre", "2-pre", "invoker", "2-post", "1-post"}
	for i, v := range expected {
		if i >= len(executionLog) || executionLog[i] != v {
			t.Errorf("Index %d: expected %s, got %v", i, v, executionLog)
		}
	}
}

func TestClientInterceptor_HeaderMutation(t *testing.T) {
	// Interceptor that sets a header
	authInterceptor := func(ctx context.Context, info *protonats.MethodInfo, req, reply proto.Message, headers nats.Header, next protonats.UnaryInvoker, opts ...protonats.CallOption) error {
		headers.Set("Authorization", "Bearer 123")
		return next(ctx, info, req, reply, headers, opts...)
	}

	// Final Invoker verifies the header arrived
	invoker := func(ctx context.Context, info *protonats.MethodInfo, req, reply proto.Message, headers nats.Header, opts ...protonats.CallOption) error {
		val := headers.Get("Authorization")
		if val != "Bearer 123" {
			return errors.New("Header not found in invoker")
		}
		return nil
	}

	options := impl.ProcessClientOptions(protonats.WithUnaryInterceptor(authInterceptor))
	chain := options.UnaryInterceptor

	headers := make(nats.Header)
	err := chain(context.Background(), nil, nil, nil, headers, invoker)

	if err != nil {
		t.Errorf("Header propagation failed: %v", err)
	}
}

func TestClientInterceptor_ErrorStopsChain(t *testing.T) {
	invokerCalled := false

	// This interceptor returns error immediately
	badInterceptor := func(ctx context.Context, info *protonats.MethodInfo, req, reply proto.Message, headers nats.Header, next protonats.UnaryInvoker, opts ...protonats.CallOption) error {
		return errors.New("blocked")
	}

	invoker := func(ctx context.Context, info *protonats.MethodInfo, req, reply proto.Message, headers nats.Header, opts ...protonats.CallOption) error {
		invokerCalled = true
		return nil
	}

	options := impl.ProcessClientOptions(protonats.WithUnaryInterceptor(badInterceptor))
	chain := options.UnaryInterceptor

	// Note: We simulate the geerated code creating the map
	err := chain(context.Background(), nil, nil, nil, make(nats.Header), invoker)

	if err == nil || err.Error() != "blocked" {
		t.Errorf("Expected 'blocked' error, got %v", err)
	}
	if invokerCalled {
		t.Error("Invoker should not have been called")
	}
}

func TestUnaryServerInterceptor_Ordering(t *testing.T) {
	var log []string

	// Outer: Logging
	loggingInt := func(ctx context.Context, req proto.Message, info *protonats.MethodInfo, next protonats.UnaryHandler) (proto.Message, error) {
		log = append(log, "log-start")
		resp, err := next(ctx, req)
		log = append(log, "log-end")
		return resp, err
	}

	// Inner: Auth
	authInt := func(ctx context.Context, req proto.Message, info *protonats.MethodInfo, next protonats.UnaryHandler) (proto.Message, error) {
		log = append(log, "auth-check")
		return next(ctx, req) // Auth doesn't do post-processing, just returns
	}

	// Actual User Implementation
	userHandler := func(ctx context.Context, req proto.Message) (proto.Message, error) {
		log = append(log, "handler")
		return &TestMsg{}, nil
	}

	// Build the chain
	options := impl.ProcessServerOptions(&micro.Config{}, protonats.ChainUnaryInterceptors(loggingInt, authInt))
	chain := options.UnaryInterceptor

	// Run
	_, err := chain(context.Background(), &TestMsg{}, &protonats.MethodInfo{}, userHandler)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Expected: Log wraps Auth, Auth wraps Handler
	expected := []string{"log-start", "auth-check", "handler", "log-end"}

	if len(log) != len(expected) {
		t.Fatalf("Length mismatch. Got %v, expected %v", log, expected)
	}

	for i, v := range expected {
		if log[i] != v {
			t.Errorf("Step %d: got %s, expected %s", i, log[i], v)
		}
	}
}

func TestUnaryServerInterceptor_ContextModification(t *testing.T) {
	// Interceptor injects data into Context
	ctxInt := func(ctx context.Context, req proto.Message, info *protonats.MethodInfo, next protonats.UnaryHandler) (proto.Message, error) {
		newCtx := context.WithValue(ctx, "userID", "u-123")
		return next(newCtx, req)
	}

	// Handler verifies it received the modified context
	handler := func(ctx context.Context, req proto.Message) (proto.Message, error) {
		val := ctx.Value("userID")
		if val != "u-123" {
			return nil, errors.New("context value missing")
		}
		return &TestMsg{}, nil
	}

	options := impl.ProcessServerOptions(&micro.Config{}, protonats.UnaryInterceptor(ctxInt))
	chain := options.UnaryInterceptor

	// Run
	_, err := chain(context.Background(), &TestMsg{}, &protonats.MethodInfo{}, handler)
	if err != nil {
		t.Errorf("Context propagation failed: %v", err)
	}
}
