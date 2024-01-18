package ctxutil

import (
	"context"

	"github.com/ironzhang/superlib/trace"
)

func ContextHook(ctx context.Context) (args []interface{}) {
	tr, ok := trace.ParseContext(ctx)
	if ok {
		args = append(args, "trace", tr)
	}
	rid, ok := trace.ParseRequestID(ctx)
	if ok {
		args = append(args, "request_id", rid)
	}
	return args
}
