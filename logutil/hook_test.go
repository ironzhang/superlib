package logutil

import (
	"context"
	"reflect"
	"testing"

	"github.com/ironzhang/superlib/trace"
)

func TestContextHook(t *testing.T) {
	tr := trace.Trace{
		TraceID:  "TraceID",
		ParentID: "ParentID",
		SpanID:   "SpanID",
	}
	rid := "rid"

	tests := []struct {
		ctx  context.Context
		args []interface{}
	}{
		{
			ctx:  context.Background(),
			args: nil,
		},
		{
			ctx:  trace.WithContext(context.Background(), tr),
			args: []interface{}{"trace", tr},
		},
		{
			ctx:  trace.WithRequestID(context.Background(), rid),
			args: []interface{}{"request_id", rid},
		},
		{
			ctx:  trace.WithRequestID(trace.WithContext(context.Background(), tr), rid),
			args: []interface{}{"trace", tr, "request_id", rid},
		},
	}
	for i, tt := range tests {
		args := ContextHook(tt.ctx)
		if got, want := args, tt.args; !reflect.DeepEqual(got, want) {
			t.Errorf("%d: args: got %v, want %v", i, got, want)
		} else {
			t.Logf("%d: args: got %v", i, got)
		}
	}
}
