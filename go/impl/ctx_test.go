package impl_test

import (
	"context"
	"testing"

	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/require"
	"xiam.li/protonats/go/impl"
	"xiam.li/protonats/go/protonats"
)

var (
	keys = struct {
		TraceID string
		AppID   string
	}{
		TraceID: "Trace-ID",
		AppID:   "App-ID",
	}
	values = struct {
		TraceID    string
		App1, App2 string
	}{
		TraceID: "abc-1",
		App1:    "my-app",
		App2:    "my-app-v2",
	}
)

func TestHeadersFromContext_Server(t *testing.T) {
	require.Nil(t, protonats.HeadersFromContext(context.Background()))

	headers := nats.Header{"Auth-Token": []string{"secret-123"}}

	ctx := impl.NewContextWithHeaders(context.Background(), headers)

	got := protonats.HeadersFromContext(ctx)

	require.NotNil(t, got)
	require.Equal(t, "secret-123", got.Get("Auth-Token"))
}

func TestWithOutgoingHeader(t *testing.T) {

	t.Run("Single", func(t *testing.T) {
		ctx := protonats.WithOutgoingHeader(context.Background(), keys.TraceID, values.TraceID)

		headers := protonats.HeadersFromOutgoingContext(ctx)

		require.Equal(t, values.TraceID, headers.Get(keys.TraceID))
		require.Len(t, headers, 1)
	})

	t.Run("Merging", func(t *testing.T) {
		// First interceptor adds initial header
		ctx := protonats.WithOutgoingHeader(context.Background(), keys.AppID, values.App1)

		// Second interceptors adds different header
		ctx = protonats.WithOutgoingHeader(ctx, keys.TraceID, values.TraceID)

		// Third interceptor adds duplicate key, which should merge
		ctx = protonats.WithOutgoingHeader(ctx, keys.AppID, values.App2)

		headers := protonats.HeadersFromOutgoingContext(ctx)

		require.Len(t, headers, 2)

		require.Len(t, headers.Values(keys.TraceID), 1)
		require.Equal(t, values.TraceID, headers.Get(keys.TraceID))

		require.Len(t, headers.Values(keys.AppID), 2)
		require.ElementsMatch(t, []string{values.App1, values.App2}, headers.Values(keys.AppID))
	})

	t.Run("MapMerge", func(t *testing.T) {
		ctx := protonats.WithOutgoingHeader(context.Background(), keys.TraceID, values.TraceID)

		newHeaders := nats.Header{keys.AppID: []string{values.App1}}

		ctx = protonats.WithOutgoingHeaders(ctx, newHeaders)

		headers := protonats.HeadersFromOutgoingContext(ctx)

		require.Len(t, headers, 2)

		require.Len(t, headers.Values(keys.TraceID), 1)
		require.Equal(t, values.TraceID, headers.Get(keys.TraceID))

		require.Len(t, headers.Values(keys.AppID), 1)
		require.Equal(t, values.App1, headers.Get(keys.AppID))
	})
}

func TestImmutability_ModificationCheck(t *testing.T) {
	parentCtx := protonats.WithOutgoingHeader(context.Background(), "Key", "Parent")

	// Child context adds a value
	childCtx := protonats.WithOutgoingHeader(parentCtx, "Key", "Child")

	// GrandChild context adds another value
	inheritedChildCtx := protonats.WithOutgoingHeader(childCtx, "SecondChild", "GrandChild")

	parentHeaders := protonats.HeadersFromOutgoingContext(parentCtx)
	childHeaders := protonats.HeadersFromOutgoingContext(childCtx)
	grandChildHeaders := protonats.HeadersFromOutgoingContext(inheritedChildCtx)

	// Verify each context has expected number of headers

	require.Len(t, parentHeaders, 1)
	require.Contains(t, parentHeaders, "Key")
	require.Len(t, parentHeaders.Values("Key"), 1)
	require.Equal(t, "Parent", parentHeaders.Get("Key"))

	require.Len(t, childHeaders, 1)
	require.Contains(t, childHeaders, "Key")
	require.Len(t, childHeaders.Values("Key"), 2)
	require.ElementsMatch(t, []string{"Parent", "Child"}, childHeaders.Values("Key"))

	require.Len(t, grandChildHeaders, 2)
	require.Contains(t, grandChildHeaders, "Key")
	require.Contains(t, grandChildHeaders, "SecondChild")
	require.Len(t, grandChildHeaders.Values("Key"), 2)
	require.ElementsMatch(t, []string{"Parent", "Child"}, grandChildHeaders.Values("Key"))
	require.Len(t, grandChildHeaders.Values("SecondChild"), 1)
	require.Equal(t, "GrandChild", grandChildHeaders.Get("SecondChild"))
}

func TestWithOutgoingHeaderMap(t *testing.T) {
	t.Run("AddToEmpty", func(t *testing.T) {
		input := map[string]string{
			keys.TraceID: values.TraceID,
			keys.AppID:   values.App1,
		}

		ctx := protonats.WithOutgoingHeaderMap(context.Background(), input)

		headers := protonats.HeadersFromOutgoingContext(ctx)

		require.Len(t, headers, 2)

		require.Len(t, headers.Values(keys.TraceID), 1)
		require.Equal(t, values.TraceID, headers.Get(keys.TraceID))

		require.Equal(t, values.App1, headers.Get(keys.AppID))
		require.Len(t, headers.Values(keys.AppID), 1)
	})

	t.Run("MergeExisting", func(t *testing.T) {
		ctx := protonats.WithOutgoingHeader(context.Background(), keys.TraceID, values.TraceID)

		input := map[string]string{keys.AppID: values.App1}
		ctx = protonats.WithOutgoingHeaderMap(ctx, input)

		headers := protonats.HeadersFromOutgoingContext(ctx)

		require.Len(t, headers, 2)

		require.Len(t, headers.Values(keys.TraceID), 1)
		require.Equal(t, values.TraceID, headers.Get(keys.TraceID))

		require.Equal(t, values.App1, headers.Get(keys.AppID))
		require.Len(t, headers.Values(keys.AppID), 1)
	})

	t.Run("AppendDuplicate", func(t *testing.T) {
		ctx := protonats.WithOutgoingHeader(context.Background(), keys.AppID, values.App1)

		input := map[string]string{keys.AppID: values.App2}
		ctx = protonats.WithOutgoingHeaderMap(ctx, input)

		headers := protonats.HeadersFromOutgoingContext(ctx)

		require.Len(t, headers, 1)
		require.Len(t, headers.Values(keys.AppID), 2)
		require.ElementsMatch(t, []string{values.App1, values.App2}, headers.Values(keys.AppID))
	})
}

func TestSetOutgoingHeaders(t *testing.T) {
	t.Run("ReplaceHeaders", func(t *testing.T) {
		ctx := protonats.WithOutgoingHeader(context.Background(), keys.AppID, values.App1)

		// Replacement
		newHeaders := nats.Header{
			keys.AppID:   []string{values.App2},
			keys.TraceID: []string{values.TraceID},
		}

		ctx = protonats.SetOutgoingHeaders(ctx, newHeaders)

		headers := protonats.HeadersFromOutgoingContext(ctx)

		require.Len(t, headers, 2)

		require.Len(t, headers.Values(keys.AppID), 1)
		require.Equal(t, values.App2, headers.Get(keys.AppID))

		require.Len(t, headers.Values(keys.TraceID), 1)
		require.Equal(t, values.TraceID, headers.Get(keys.TraceID))
	})

	t.Run("ReplaceMap", func(t *testing.T) {
		ctx := protonats.WithOutgoingHeader(context.Background(), keys.AppID, values.App1)

		// Replacement
		input := map[string]string{keys.TraceID: values.TraceID}

		ctx = protonats.SetOutgoingHeaderMap(ctx, input)

		headers := protonats.HeadersFromOutgoingContext(ctx)

		require.Len(t, headers, 1)

		require.Len(t, headers.Values(keys.TraceID), 1)
		require.Equal(t, values.TraceID, headers.Get(keys.TraceID))
	})

	t.Run("EmptyMap", func(t *testing.T) {
		ctx := protonats.WithOutgoingHeader(context.Background(), keys.AppID, values.App1)

		ctx = protonats.SetOutgoingHeaderMap(ctx, make(map[string]string))

		headers := protonats.HeadersFromOutgoingContext(ctx)

		require.NotNil(t, headers)
		require.Empty(t, headers)
	})
}
