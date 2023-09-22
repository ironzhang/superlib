package interceptors

import (
	"context"

	"github.com/ironzhang/superlib/httputils"
	"github.com/ironzhang/superlib/httputils/httpclient"
	"github.com/ironzhang/superlib/trace"
)

func TraceInterceptor() httpclient.Interceptor {
	return func(ctx context.Context, info *httpclient.InvokeInfo, args, reply interface{}, invoker httpclient.Invoker) error {
		// trace
		tr, ok := trace.ParseContext(ctx)
		if !ok {
			tr.TraceID = trace.NewTraceID()
			tr.SpanID = trace.NewSpanID()
			ctx = trace.WithContext(ctx, tr)
		}
		info.Header.Set(httputils.HeaderXTraceID, tr.TraceID)
		info.Header.Set(httputils.HeaderXSpanID, tr.SpanID)

		// request id
		rid := trace.NewRequestID()
		ctx = trace.WithRequestID(ctx, rid)
		info.Header.Set(httputils.HeaderXRequestID, rid)

		return invoker(ctx, info, args, reply)
	}
}
