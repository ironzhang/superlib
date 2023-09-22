package testutil

import (
	"testing"
)

type S struct {
	A string
	b string
	c int
}

func TestCompare(t *testing.T) {
	tests := []struct {
		x      interface{}
		y      interface{}
		filter CompareFilterFunc
		result bool
	}{
		{
			x:      1,
			y:      1,
			filter: nil,
			result: true,
		},
		{
			x:      1,
			y:      2,
			filter: nil,
			result: false,
		},
		{
			x:      &S{A: "A", b: "b", c: 1},
			y:      &S{A: "A", b: "b", c: 1},
			filter: nil,
			result: true,
		},
		{
			x:      &S{A: "A", b: "b", c: 1},
			y:      &S{A: "A", b: "b", c: 2},
			filter: nil,
			result: false,
		},
		{
			x: &S{A: "A", b: "b", c: 1},
			y: &S{A: "A", b: "b", c: 2},
			filter: func(pkg, typ, field string) bool {
				if pkg == "github.com/ironzhang/superlib/testutil" && typ == "S" && field == "c" {
					return true
				}
				return false
			},
			result: true,
		},
	}
	for i, tt := range tests {
		result := Compare(tt.x, tt.y, tt.filter)
		if got, want := result, tt.result; got != want {
			t.Errorf("%d: x=%v, y=%v: result: got %v, want %v", i, tt.x, tt.y, got, want)
		} else {
			t.Logf("%d: x=%v, y=%v: result: %v", i, tt.x, tt.y, got)
		}
	}
}
