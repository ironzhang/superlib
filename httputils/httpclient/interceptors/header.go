package interceptors

import (
	"context"
	"net/http"

	"github.com/ironzhang/superlib/httputils/httpclient"
)

type WithHeaderConfig struct {
	Header http.Header
}

func (p *WithHeaderConfig) Interceptor(ctx context.Context, info *httpclient.InvokeInfo, args, reply interface{}, invoker httpclient.Invoker) error {
	for k, vs := range p.Header {
		for _, v := range vs {
			if k != "" && v != "" {
				info.Header.Add(k, v)
			}
		}
	}
	return invoker(ctx, info, args, reply)
}

func WithHeaderInterceptor(h http.Header) httpclient.Interceptor {
	cfg := WithHeaderConfig{
		Header: h,
	}
	return cfg.Interceptor
}
