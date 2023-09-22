package echoutil

import (
	"net/http"

	"github.com/labstack/echo"

	"github.com/ironzhang/superlib/httputils"
	"github.com/ironzhang/superlib/trace"
)

// TraceConfig trace 中间件配置
type TraceConfig struct {
	NewTraceID func() string
	NewSpanID  func() string
}

func (p *TraceConfig) MiddlewareFunc(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		r := c.Request()
		tr := p.parseTrace(r.Header)
		ctx := trace.WithContext(r.Context(), tr)
		if rid := c.Response().Header().Get(httputils.HeaderXRequestID); rid != "" {
			ctx = trace.WithRequestID(ctx, rid)
		}
		c.SetRequest(r.WithContext(ctx))
		return next(c)
	}
}

func (p *TraceConfig) parseTrace(header http.Header) trace.Trace {
	traceID := header.Get(httputils.HeaderXTraceID)
	if traceID == "" {
		traceID = p.NewTraceID()
	}
	parentID := header.Get(httputils.HeaderXSpanID)
	if parentID == "" {
		parentID = p.NewSpanID()
	}
	return trace.Trace{
		TraceID:  traceID,
		ParentID: parentID,
		SpanID:   p.NewSpanID(),
	}
}

// TraceMiddleware Trace 中间件
func TraceMiddleware() echo.MiddlewareFunc {
	cfg := TraceConfig{
		NewTraceID: trace.NewTraceID,
		NewSpanID:  trace.NewSpanID,
	}
	return cfg.MiddlewareFunc
}
