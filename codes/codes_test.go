package codes_test

import (
	"testing"

	"github.com/ironzhang/superlib/codes"
)

func TestRegisterCodes(t *testing.T) {
	tests := []struct {
		code codes.Code
		desc string
	}{
		{code: -1, desc: "-1"},
		{code: -2, desc: "-2"},
		{code: -3, desc: "-3"},
		{code: -4, desc: "-4"},
	}

	for _, tt := range tests {
		codes.Register(tt.code, tt.desc)
	}
	for _, tt := range tests {
		if got, want := tt.code.String(), tt.desc; got != want {
			t.Errorf("%q != %q", got, want)
		}
	}
}

func TestPrintCodes(t *testing.T) {
	codes := []codes.Code{
		codes.OK,
		codes.Unknown,
		codes.Internal,
		codes.NotFound,
		codes.InvalidParams,
	}
	for _, code := range codes {
		t.Logf("%d=%s\n", code, code)
	}
}
