package interceptors

import (
	"context"
	"reflect"
	"time"

	"github.com/ironzhang/tlog"

	"github.com/ironzhang/superlib/httputils"
	"github.com/ironzhang/superlib/httputils/httpclient"
)

type AccessLogInterceptorConfig struct {
	DisableBody bool
}

func (p *AccessLogInterceptorConfig) Interceptor(ctx context.Context, info *httpclient.InvokeInfo, args, reply interface{}, invoker httpclient.Invoker) error {
	log := tlog.Named("access.out")
	opts := httpclient.ParseAccessLogOptions(ctx)

	// 输出请求日志
	p.printRequest(ctx, log, opts, info, args)

	// 透传访问日志输出模式
	if opts.ServerMode != "" {
		info.Header.Set(httputils.HeaderAccessLogPrintMode, string(opts.ServerMode))
	}

	// 调用 invoker
	start := time.Now()
	err := invoker(ctx, info, args, reply)

	// 输出响应日志
	p.printResponse(ctx, log, opts, info, reply, time.Since(start), err)

	// 返回错误
	return err
}

func (p *AccessLogInterceptorConfig) printRequest(ctx context.Context, log tlog.Logger, opts httpclient.AccessLogOptions, info *httpclient.InvokeInfo, args interface{}) {
	log = log.WithContext(ctx).WithArgs("method", info.Method, "path", info.Path, "remote_addr", info.Addr, "query", info.Query, "header", info.Header)
	if opts.ClientMode == httputils.DisableAccessBody || p.DisableBody || args == nil {
		log.Infof("http client request")
	} else {
		log.Infof("http client request: %s {%+v}", reflect.TypeOf(args).String(), args)
	}
}

func (p *AccessLogInterceptorConfig) printResponse(ctx context.Context, log tlog.Logger, opts httpclient.AccessLogOptions, info *httpclient.InvokeInfo, reply interface{},
	latency time.Duration, err error) {
	log = log.WithContext(ctx).WithArgs("method", info.Method, "path", info.Path, "remote_addr", info.Addr, "latency", latency, "error", err)
	if opts.ClientMode == httputils.DisableAccessBody || p.DisableBody || reply == nil {
		log.Infof("http client response")
	} else {
		log.Infof("http client response: %s {%+v}", reflect.TypeOf(reply).String(), reply)
	}
}

func AccessLogInterceptor() httpclient.Interceptor {
	cfg := AccessLogInterceptorConfig{}
	return cfg.Interceptor
}
