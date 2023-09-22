package trace

import (
	"context"

	"github.com/labstack/gommon/random"

	"github.com/ironzhang/superlib/uuid"
)

const (
	TraceID  = "traceid"
	ParentID = "parentid"
	SpanID   = "spanid"
)

type Trace struct {
	TraceID  string
	ParentID string
	SpanID   string
}

type traceKey struct{}

func NewTrace() Trace {
	return Trace{
		TraceID: NewTraceID(),
		SpanID:  NewSpanID(),
	}
}

func WithContext(ctx context.Context, trace Trace) context.Context {
	return context.WithValue(ctx, traceKey{}, trace)
}

func ParseContext(ctx context.Context) (Trace, bool) {
	trace, ok := ctx.Value(traceKey{}).(Trace)
	return trace, ok
}

func NewContext(ctx context.Context) context.Context {
	tr, ok := ParseContext(ctx)
	if !ok {
		return WithContext(ctx, NewTrace())
	}

	tr.ParentID = tr.SpanID
	tr.SpanID = NewSpanID()
	return WithContext(ctx, tr)
}

func NewTraceID() string {
	return uuid.New().String()
}

func NewSpanID() string {
	return uuid.New().String()
}

func NewRequestID() string {
	return random.String(32)
}

type requestIDKey struct{}

func WithRequestID(ctx context.Context, rid string) context.Context {
	return context.WithValue(ctx, requestIDKey{}, rid)
}

func ParseRequestID(ctx context.Context) (string, bool) {
	rid, ok := ctx.Value(requestIDKey{}).(string)
	return rid, ok
}
