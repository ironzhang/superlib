package echorpc

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo"

	"github.com/ironzhang/superlib/testutil"
)

type Args struct {
	A, B int
}

type Reply struct {
	C int
}

type Arith int

func (t *Arith) Add(ctx context.Context, args Args, reply *Reply) error {
	reply.C = args.A + args.B
	return nil
}

func (t *Arith) Sub(ctx context.Context, args Args, reply *Reply) error {
	reply.C = args.A - args.B
	return nil
}

func (t *Arith) Mul(ctx context.Context, args *Args, reply *Reply) error {
	reply.C = args.A * args.B
	return nil
}

func (t *Arith) Div(ctx context.Context, args Args, reply *Reply) error {
	if args.B == 0 {
		return errors.New("divide by zero")
	}
	reply.C = args.A / args.B
	return nil
}

func NewTestEcho(t *testing.T) *echo.Echo {
	var a Arith
	e := echo.New()
	e.Debug = true
	e.POST("/add", HandlerFunc(a.Add))
	e.POST("/sub", HandlerFunc(a.Sub))
	e.POST("/mul", HandlerFunc(a.Mul))
	e.POST("/div", HandlerFunc(a.Div))
	return e
}

func CallArith(addr, method, path string, a, b int) (c int, err error) {
	args := Args{a, b}

	var buf bytes.Buffer
	if err = json.NewEncoder(&buf).Encode(args); err != nil {
		return 0, err
	}

	req, err := http.NewRequest(method, fmt.Sprintf("%s%s", addr, path), &buf)
	if err != nil {
		return 0, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		data, _ := ioutil.ReadAll(resp.Body)
		return 0, fmt.Errorf("%d: %s: %s", resp.StatusCode, http.StatusText(resp.StatusCode), data)
	}

	var reply Reply
	if err = json.NewDecoder(resp.Body).Decode(&reply); err != nil {
		return 0, err
	}
	return reply.C, nil
}

func TestEcho(t *testing.T) {
	e := NewTestEcho(t)
	s := httptest.NewServer(e)

	tests := []struct {
		method string
		path   string
		err    string
		a      int
		b      int
		c      int
	}{
		{method: "POST", path: "/add", err: "", a: 1, b: 2, c: 3},
		{method: "POST", path: "/sub", err: "", a: 1, b: 2, c: -1},
		{method: "POST", path: "/mul", err: "", a: 1, b: 2, c: 2},
		{method: "POST", path: "/div", err: "", a: 1, b: 2, c: 0},
		{method: "POST", path: "/div", err: "divide by zero", a: 1, b: 0, c: 0},
	}
	for i, tt := range tests {
		c, err := CallArith(s.URL, tt.method, tt.path, tt.a, tt.b)
		if got, want := err, tt.err; !testutil.MatchError(t, got, want) {
			t.Fatalf("%d: error: got %v, want %v", i, got, want)
		}
		if err != nil {
			t.Logf("%d: call arith: %v", i, err)
			continue
		}
		if got, want := c, tt.c; got != want {
			t.Fatalf("%d: c: got %v, want %v", i, got, want)
		}
		t.Logf("%d: c: %v", i, c)
	}
}
