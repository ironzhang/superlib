package httpclient

import (
	"net/http"
)

type CallOption func(info *InvokeInfo)

func WithHeader(h http.Header) CallOption {
	return func(info *InvokeInfo) {
		for key, values := range h {
			for _, value := range values {
				info.Header.Add(key, value)
			}
		}
	}
}
