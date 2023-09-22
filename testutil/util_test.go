package testutil

import (
	"errors"
	"testing"
)

type TestTask struct {
	n int
}

func (p *TestTask) Do(t testing.TB) error {
	p.n++
	return nil
}

func TestRun(t *testing.T) {
	var t1, t2 TestTask
	tasks := []Task{&t1, &t2}
	Run(t, tasks)
	if got, want := t1.n, 1; got != want {
		t.Errorf("t1: got %v, want %v", got, want)
	}
	if got, want := t2.n, 1; got != want {
		t.Errorf("t2: got %v, want %v", got, want)
	}
}

func TestMatchError(t *testing.T) {
	tests := []struct {
		err    error
		errstr string
		result bool
	}{
		{
			err:    nil,
			errstr: "",
			result: true,
		},
		{
			err:    nil,
			errstr: "error",
			result: false,
		},
		{
			err:    errors.New("error"),
			errstr: "",
			result: false,
		},
		{
			err:    errors.New("error"),
			errstr: "error",
			result: true,
		},
		{
			err:    errors.New("error1"),
			errstr: "error",
			result: true,
		},
		{
			err:    errors.New("error"),
			errstr: "error1",
			result: false,
		},
	}
	for i, tt := range tests {
		ok := MatchError(t, tt.err, tt.errstr)
		if got, want := ok, tt.result; got != want {
			t.Errorf("%d: MatchError(%v, %q): got %v, want: %v", i, tt.err, tt.errstr, got, want)
		} else {
			t.Logf("%d: MatchError(%v, %q): got %v", i, tt.err, tt.errstr, got)
		}
	}
}
