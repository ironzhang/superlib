package trace

import (
	"context"
	"testing"
)

func TestTraceContext(t *testing.T) {
	trace := Trace{
		TraceID:  NewTraceID(),
		ParentID: NewSpanID(),
		SpanID:   NewSpanID(),
	}
	tests := []struct {
		ctx   context.Context
		ok    bool
		trace Trace
	}{
		{
			ctx: context.Background(),
			ok:  false,
		},
		{
			ctx:   WithContext(context.Background(), trace),
			ok:    true,
			trace: trace,
		},
	}
	for i, tt := range tests {
		tr, ok := ParseContext(tt.ctx)
		if got, want := ok, tt.ok; got != want {
			t.Errorf("%d: ok: got %v, want %v", i, got, want)
		} else {
			t.Logf("%d: ok: got %v", i, got)
		}
		if !ok {
			continue
		}
		if got, want := tr, tt.trace; got != want {
			t.Errorf("%d: trace: got %v, want %v", i, got, want)
		} else {
			t.Logf("%d: trace: got %v", i, got)
		}
	}
}

func TestRequestID(t *testing.T) {
	requestID := "request_id"
	tests := []struct {
		ctx context.Context
		ok  bool
		rid string
	}{
		{
			ctx: context.Background(),
			ok:  false,
		},
		{
			ctx: WithRequestID(context.Background(), requestID),
			ok:  true,
			rid: requestID,
		},
	}
	for i, tt := range tests {
		rid, ok := ParseRequestID(tt.ctx)
		if got, want := ok, tt.ok; got != want {
			t.Errorf("%d: ok: got %v, want %v", i, got, want)
		} else {
			t.Logf("%d: ok: got %v", i, got)
		}
		if !ok {
			continue
		}
		if got, want := rid, tt.rid; got != want {
			t.Errorf("%d: rid: got %v, want %v", i, got, want)
		} else {
			t.Logf("%d: rid: got %v", i, got)
		}
	}
}
