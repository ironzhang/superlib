package httpclient

import (
	"context"

	"github.com/ironzhang/superlib/httputils"
)

type AccessLogOptions struct {
	ClientMode httputils.AccessLogPrintMode
	ServerMode httputils.AccessLogPrintMode
}

// SetAccessLogOptions 设置本次调用访问日志模式
func SetAccessLogOptions(ctx context.Context, opts AccessLogOptions) context.Context {
	return context.WithValue(ctx, httputils.HeaderAccessLogPrintMode, opts)
}

// ParseAccessLogOptions 从 context 中解析本次调用的访问日志模式
func ParseAccessLogOptions(ctx context.Context) AccessLogOptions {
	v, ok := ctx.Value(httputils.HeaderAccessLogPrintMode).(AccessLogOptions)
	if !ok {
		return AccessLogOptions{}
	}
	return v
}
