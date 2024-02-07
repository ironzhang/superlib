package passthrough

import "github.com/ironzhang/superlib/httputils/httpclient/resolver"

type passthroughResolver struct {
}

func (p passthroughResolver) Scheme() string {
	return "passthrough"
}

func (p passthroughResolver) Resolve(endpoint string) (string, error) {
	return endpoint, nil
}

func init() {
	resolver.Register(passthroughResolver{})
}
