package codes

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"testing"
)

func TestNewError(t *testing.T) {
	tests := []struct {
		code Code
		err  error
		estr string
	}{
		{
			code: -1,
			err:  fmt.Errorf("error1"),
			estr: "{code: -1, desc: code(-1), wrap: error1}",
		},
		{
			code: Unknown,
			err:  fmt.Errorf("error2"),
			estr: "{code: 1, desc: unknown, wrap: error2}",
		},
	}
	for i, tt := range tests {
		e := NewError(tt.code, tt.err)
		if got, want := e.Code(), tt.code; got != want {
			t.Errorf("%d: code: got %v, want %v", i, got, want)
		}
		if got, want := e.Unwrap(), tt.err; got.Error() != want.Error() {
			t.Errorf("%d: unwrap: got %v, want %v", i, got, want)
		}
		if got, want := e.Error(), tt.estr; got != want {
			t.Errorf("%d: error: got %v, want %v", i, got, want)
		}
	}
}

func TestErrorf(t *testing.T) {
	tests := []struct {
		code   Code
		format string
		args   []interface{}
		estr   string
	}{
		{
			code:   -1,
			format: "error %d",
			args:   []interface{}{1},
			estr:   "{code: -1, desc: code(-1), wrap: error 1}",
		},
		{
			code:   Unknown,
			format: "error %s %s",
			args:   []interface{}{"1", "2"},
			estr:   "{code: 1, desc: unknown, wrap: error 1 2}",
		},
	}
	for i, tt := range tests {
		e := Errorf(tt.code, tt.format, tt.args...)
		if got, want := e.Code(), tt.code; got != want {
			t.Errorf("%d: code: got %v, want %v", i, got, want)
		}
		if got, want := e.Unwrap(), fmt.Errorf(tt.format, tt.args...); got.Error() != want.Error() {
			t.Errorf("%d: unwrap: got %v, want %v", i, got, want)
		}
		if got, want := e.Error(), tt.estr; got != want {
			t.Errorf("%d: error: got %v, want %v", i, got, want)
		}
	}
}

func TestGetErrorCode(t *testing.T) {
	tests := []struct {
		err  error
		code Code
	}{
		{
			err:  nil,
			code: OK,
		},
		{
			err:  NewError(Unknown, errors.New("error")),
			code: Unknown,
		},
		{
			err:  NewError(Internal, NewError(Unknown, errors.New("error"))),
			code: Internal,
		},
		{
			err:  fmt.Errorf("error: %w", NewError(Unknown, errors.New("error"))),
			code: Unknown,
		},
		{
			err:  fmt.Errorf("error: %w", fmt.Errorf("error: %w", NewError(-1, errors.New("error")))),
			code: -1,
		},
		{
			err: fmt.Errorf("http client do: %w", &url.Error{
				Op:  "Get",
				URL: "http://127.0.0.1:8000/v1/api/cluster/list",
				Err: context.Canceled,
			}),
			code: Canceled,
		},
	}
	for i, tt := range tests {
		code := GetErrorCode(tt.err)
		if got, want := code, tt.code; got != want {
			t.Errorf("%d: code: got %v, want %v", i, got, want)
		} else {
			t.Logf("%d: code: got %v", i, got)
		}
	}
}
