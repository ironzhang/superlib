package echoutil

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	"github.com/ironzhang/superlib/httputils"
	"github.com/ironzhang/superlib/trace"
)

type TraceMiddlewareHandler struct {
	hit     bool
	context context.Context
}

func (h *TraceMiddlewareHandler) Handle(c echo.Context) error {
	h.hit = true
	h.context = c.Request().Context()
	return nil
}

func TestTraceMiddleware(t *testing.T) {
	var cfg TraceConfig
	cfg.NewTraceID = func() string {
		return "DefaultTraceID"
	}
	cfg.NewSpanID = func() string {
		return "DefaultSpanID"
	}

	req1, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatalf("new request: %v", err)
	}

	req2, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	req2.Header.Set(httputils.HeaderXTraceID, "TraceID-Request2")
	req2.Header.Set(httputils.HeaderXSpanID, "SpanID-Request2")

	req3, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	req3.Header.Set(echo.HeaderXRequestID, "RequestID-Request3")

	tests := []struct {
		middlewares []echo.MiddlewareFunc
		request     *http.Request
		ok          bool
		trace       trace.Trace
		rid         string
	}{
		{
			middlewares: nil,
			request:     req1,
			ok:          false,
		},
		{
			middlewares: []echo.MiddlewareFunc{cfg.MiddlewareFunc},
			request:     req1,
			ok:          true,
			trace: trace.Trace{
				TraceID:  "DefaultTraceID",
				ParentID: "DefaultSpanID",
				SpanID:   "DefaultSpanID",
			},
			rid: "",
		},
		{
			middlewares: []echo.MiddlewareFunc{cfg.MiddlewareFunc},
			request:     req2,
			ok:          true,
			trace: trace.Trace{
				TraceID:  "TraceID-Request2",
				ParentID: "SpanID-Request2",
				SpanID:   "DefaultSpanID",
			},
			rid: "",
		},
		{
			middlewares: []echo.MiddlewareFunc{middleware.RequestID(), cfg.MiddlewareFunc},
			request:     req3,
			ok:          true,
			trace: trace.Trace{
				TraceID:  "DefaultTraceID",
				ParentID: "DefaultSpanID",
				SpanID:   "DefaultSpanID",
			},
			rid: "RequestID-Request3",
		},
	}
	for i, tt := range tests {
		h := TraceMiddlewareHandler{}
		e := echo.New()
		e.GET("/", h.Handle, tt.middlewares...)
		w := httptest.NewRecorder()
		e.ServeHTTP(w, tt.request)

		if !h.hit {
			t.Errorf("%d: trace middleware handler does not hit", i)
			continue
		}
		tr, ok := trace.ParseContext(h.context)
		if got, want := ok, tt.ok; got != want {
			t.Errorf("%d: ok: got %v, want %v", i, got, want)
			continue
		}
		if !ok {
			continue
		}
		if got, want := tr, tt.trace; got != want {
			t.Errorf("%d: trace: got %v, want %v", i, got, want)
			continue
		}
		rid, _ := trace.ParseRequestID(h.context)
		if got, want := rid, tt.rid; got != want {
			t.Errorf("%d: request_id: got %v, want %v", i, got, want)
			continue
		}
		t.Logf("%d: trace=%v, rid=%q", i, tr, rid)
	}
}
