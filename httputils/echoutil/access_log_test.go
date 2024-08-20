package echoutil

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo"

	"github.com/ironzhang/tlog/zaplog"

	"github.com/ironzhang/superlib/ctxutil"
)

type AccessLogMiddlewareHandler struct {
	Body []byte
}

func (p *AccessLogMiddlewareHandler) ReturnOK(c echo.Context) (err error) {
	req := c.Request()
	res := c.Response()
	p.Body, err = ioutil.ReadAll(req.Body)
	if err != nil {
		return err
	}
	res.Write(p.Body)
	return nil
}

func (p *AccessLogMiddlewareHandler) ReturnBadRequest(c echo.Context) (err error) {
	req := c.Request()
	res := c.Response()
	p.Body, err = ioutil.ReadAll(req.Body)
	if err != nil {
		return err
	}
	res.WriteHeader(http.StatusBadRequest)
	res.Write(p.Body)
	return errors.New("error")
}

func TestAccessLogMiddleware(t *testing.T) {
	//tlog.SetLogger(nil)
	zaplog.StdContextHook = ctxutil.ContextHook

	req1, err := http.NewRequest("GET", "/ok", nil)
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	req2, err := http.NewRequest("GET", "/badrequest", strings.NewReader("req2"))
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	req3, err := http.NewRequest("GET", "/badrequest", strings.NewReader("req3"))
	if err != nil {
		t.Fatalf("new request: %v", err)
	}

	h := AccessLogMiddlewareHandler{}
	e := echo.New()
	e.Use(TraceMiddleware())
	e.Use(AccessLogMiddleware())
	e.GET("/ok", h.ReturnOK)
	e.GET("/badrequest", h.ReturnBadRequest)

	tests := []struct {
		req    *http.Request
		status int
		body   string
	}{
		{
			req:    req1,
			status: http.StatusOK,
			body:   "",
		},
		{
			req:    req2,
			status: http.StatusBadRequest,
			body:   "req2",
		},
		{
			req:    req3,
			status: http.StatusBadRequest,
			body:   "req3",
		},
	}
	for i, tt := range tests {
		res := httptest.NewRecorder()
		e.ServeHTTP(res, tt.req)
		if got, want := res.Code, tt.status; got != want {
			t.Errorf("%d: status: got %v, want %v", i, got, want)
			continue
		}
		if got, want := string(res.Body.Bytes()), tt.body; got != want {
			t.Errorf("%d: body: got %v, want %v", i, got, want)
			continue
		}
		t.Logf("%d: status=%d, body=%s", i, res.Code, res.Body.Bytes())
	}
}
