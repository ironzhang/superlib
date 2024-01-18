package ctxutil

import (
	"context"

	"github.com/ironzhang/superlib/trace"
)

// CloneContext clones a new context from o.
func CloneContext(o context.Context) context.Context {
	ctx := context.Background()

	// trace context
	tr, ok := trace.ParseContext(o)
	if ok {
		ctx = trace.WithContext(ctx, tr)
	}

	// request id
	rid, ok := trace.ParseRequestID(o)
	if ok {
		ctx = trace.WithRequestID(ctx, rid)
	}

	return ctx
}
