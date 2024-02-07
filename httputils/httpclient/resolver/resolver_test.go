package resolver_test

import (
	"testing"

	"github.com/ironzhang/superlib/httputils/httpclient/resolver"
	"github.com/ironzhang/superlib/testutil"
)

type testResolver struct {
}

func (p testResolver) Scheme() string {
	return "test"
}

func (p testResolver) Resolve(endpoint string) (string, error) {
	return endpoint, nil
}

func init() {
	resolver.Register(testResolver{})
}

func TestResolve(t *testing.T) {
	tests := []struct {
		target resolver.Target
		addr   string
		err    string
	}{
		{
			target: resolver.Target{
				Scheme:   "test",
				Endpoint: "127.0.0.1:8000",
			},
			addr: "127.0.0.1:8000",
		},
		{
			target: resolver.Target{
				Endpoint: "127.0.0.1:8000",
			},
			err: "can not find",
		},
	}
	for i, tt := range tests {
		addr, err := resolver.Resolve(tt.target)
		if got, want := err, tt.err; !testutil.MatchError(t, got, want) {
			t.Fatalf("%d: error is not match: got %v, want %v", i, got, want)
		}
		if err != nil {
			t.Logf("%d: resolve: %v", i, err)
			continue
		}
		if got, want := addr, tt.addr; addr != addr {
			t.Fatalf("%d: addr is unexpected: got %v, want %v", i, got, want)
		}
	}
}
