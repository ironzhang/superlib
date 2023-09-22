package httpclient

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/ironzhang/tlog"

	"github.com/ironzhang/superlib/codes"
)

func interceptor1(ctx context.Context, info *InvokeInfo, args, reply interface{}, invoker Invoker) error {
	tlog.Infow("enter interceptor1")
	err := invoker(ctx, info, args, reply)
	tlog.Infow("exit interceptor1")
	return err
}

func interceptor2(ctx context.Context, info *InvokeInfo, args, reply interface{}, invoker Invoker) error {
	tlog.Infow("enter interceptor2")
	err := invoker(ctx, info, args, reply)
	tlog.Infow("exit interceptor2")
	return err
}

func serveHTTP(w http.ResponseWriter, r *http.Request) {
	tlog.Infow("serve http")
}

func NewTestServer(t *testing.T) string {
	s := httptest.NewServer(http.HandlerFunc(serveHTTP))
	u, err := url.Parse(s.URL)
	if err != nil {
		t.Fatalf("url parse: %v", err)
	}
	return u.Host
}

func TestClientCall(t *testing.T) {
}

func TestClientInvoke(t *testing.T) {
	addr := NewTestServer(t)
	c := Client{
		Addr:         addr,
		Interceptors: []Interceptor{interceptor1, interceptor2},
	}

	err := c.Invoke(context.Background(), "POST", "/path", nil, nil, nil)
	if err != nil {
		t.Fatalf("invoke: %v", err)
	}
}

func TestDefaultResultParserStatusNotOK(t *testing.T) {
	c := &Client{
		Codec: JSONCodec{},
	}
	resp := http.Response{}

	body := []byte(`{"Code":1,"Error":"some error info"}`)
	resp.Body = ioutil.NopCloser(bytes.NewReader(body))
	resp.StatusCode = 500
	resp.ContentLength = int64(len(body))

	err := defaultResultParser(context.Background(), c, &resp, nil)
	messageError, ok := err.(codes.Error)
	if !ok {
		t.Fatalf("default result parse want message error, but not")
	}
	if messageError.Code() != codes.Unknown {
		t.Fatalf("default result parse want error code unkonw, but got %d", messageError.Code())
	}
	t.Log("message error", messageError)
}

func TestDefaultResultParserStatusNotOKWithCodecError(t *testing.T) {
	c := &Client{
		Codec: JSONCodec{},
	}
	resp := http.Response{}

	body := []byte(`not a json`)
	resp.Body = ioutil.NopCloser(bytes.NewReader(body))
	resp.StatusCode = 500
	resp.ContentLength = int64(len(body))

	err := defaultResultParser(context.Background(), c, &resp, nil)
	if err == nil {
		t.Fatalf("default result parse want error, got nil")
	}
	t.Log("got error", err)
}

type result struct {
	Data string
}

func TestDefaultResultParserStatusOkWithCodecError(t *testing.T) {
	c := &Client{
		Codec: JSONCodec{},
	}
	resp := http.Response{}

	resp.Body = ioutil.NopCloser(bytes.NewReader([]byte(`not a json`)))
	resp.StatusCode = 200

	var r result

	err := defaultResultParser(context.Background(), c, &resp, &r)
	if err == nil {
		t.Fatalf("default result parse want error, got nil")
	}
	t.Log("got error", err)
}

func TestDefaultResultParserStatusOk(t *testing.T) {
	c := &Client{
		Codec: JSONCodec{},
	}
	resp := http.Response{}

	resp.Body = ioutil.NopCloser(bytes.NewReader([]byte(`{"Data":"some data"}`)))
	resp.StatusCode = 200

	var r result

	err := defaultResultParser(context.Background(), c, &resp, &r)
	if err != nil {
		t.Fatalf("default result parse got error %s", err.Error())
	}
	if r.Data != "some data" {
		t.Fatalf("default result parse want some data, got %s", r.Data)
	}
	t.Log("result", r)
}
