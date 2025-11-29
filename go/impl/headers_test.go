package impl_test

import (
	"testing"

	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/require"
	"xiam.li/protonats/go/impl"
	"xiam.li/protonats/go/protonats"
)

func TestResponseOptions(t *testing.T) {
	var capturedHeaders nats.Header
	var capturedID, capturedID2 string

	opts := impl.ProcessCallOptions(
		protonats.ResponseHeader(&capturedHeaders),
		protonats.ResponseInstanceID(&capturedID),
		protonats.ResponseInstanceID(&capturedID2),
	)

	// Simulate receiving a response with headers
	incomingHeaders := nats.Header{
		"Content-Type":       []string{"application/x-protobuf"},
		impl.ServiceIDHeader: []string{"srv-abc-123"},
	}

	opts.ApplyResponseHeaders(incomingHeaders)

	require.NotNil(t, capturedHeaders)
	require.Equal(t, "application/x-protobuf", capturedHeaders.Get("Content-Type"))
	require.Equal(t, "srv-abc-123", capturedHeaders.Get(impl.ServiceIDHeader))

	require.Equal(t, "srv-abc-123", capturedID)
	require.Equal(t, "srv-abc-123", capturedID2)

	t.Run("NilSaftey", func(t *testing.T) {
		// Ensure no panics if headers are missing or pointers are nil
		opts := impl.ProcessCallOptions(protonats.ResponseInstanceID(nil))

		opts.ApplyResponseHeaders(nil)

		opts.ApplyResponseHeaders(nats.Header{impl.ServiceIDHeader: []string{"123"}})
	})
}
