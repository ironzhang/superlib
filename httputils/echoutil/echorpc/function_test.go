package echorpc

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"
)

func TestIsExported(t *testing.T) {
	tests := []struct {
		name   string
		expect bool
	}{
		{name: "", expect: false},
		{name: "A", expect: true},
		{name: "Aa", expect: true},
		{name: "Aaaaaa", expect: true},
		{name: "a", expect: false},
		{name: "aA", expect: false},
		{name: "aAAAAA", expect: false},
		{name: "你好", expect: false},
		{name: "A你好", expect: true},
		{name: "你好A", expect: false},
		{name: "_", expect: false},
		{name: "_A", expect: false},
		{name: "A_", expect: true},
	}
	for _, tt := range tests {
		if got, want := isExported(tt.name), tt.expect; got != want {
			t.Errorf("%q: %v != %v", tt.name, got, want)
		} else {
			t.Logf("%q: %v", tt.name, got)
		}
	}
}

func TestIsExportedOrBuiltinType(t *testing.T) {
	type a struct{}
	type A struct{}
	tests := []struct {
		typ    reflect.Type
		expect bool
	}{
		{typ: reflect.TypeOf(""), expect: true},
		{typ: reflect.TypeOf(1), expect: true},
		{typ: reflect.TypeOf(1.0), expect: true},
		{typ: reflect.TypeOf(a{}), expect: false},
		{typ: reflect.TypeOf(&a{}), expect: false},
		{typ: reflect.TypeOf(A{}), expect: true},
		{typ: reflect.TypeOf(&A{}), expect: true},
	}
	for i, tt := range tests {
		if got, want := isExportedOrBuiltinType(tt.typ), tt.expect; got != want {
			t.Errorf("case%d: %v: %v != %v", i, tt.typ.String(), got, want)
		} else {
			t.Logf("case%d: %v: %v", i, tt.typ.String(), got)
		}
	}
}

type TContext interface {
	Deadline() (deadline time.Time, ok bool)
	Done() <-chan struct{}
	Err() error
	Value(key interface{}) interface{}
}

type tcontext struct {
}

func (p *tcontext) Deadline() (deadline time.Time, ok bool) {
	return
}

func (p *tcontext) Done() <-chan struct{} {
	return nil
}

func (p *tcontext) Err() error {
	return nil
}

func (p *tcontext) Value(key interface{}) interface{} {
	return nil
}

type A struct{}

type a struct{}

type correct struct{}

func (correct) Test00(context.Context, int, *int) error {
	return nil
}

func (correct) Test01(context.Context, int, *string) error {
	return nil
}

func (correct) Test02(context.Context, int, *A) error {
	return nil
}

func (correct) Test03(context.Context, int, interface{}) error {
	return nil
}

func (correct) Test10(context.Context, *int, *int) error {
	return nil
}

func (correct) Test11(context.Context, A, *int) error {
	return nil
}

func (correct) Test12(context.Context, *A, *int) error {
	return nil
}

func (correct) Test13(context.Context, interface{}, *int) error {
	return nil
}

func (correct) Test20(TContext, int, *int) error {
	return nil
}

//func (correct) Test21(*tcontext, int, *int) error {
//	return nil
//}

type incorrect struct{}

func (incorrect) Test00() {
}

func (incorrect) Test01(context.Context, int, *int, int) (error, error) {
	return nil, nil
}

func (incorrect) Test10(int, int, *int) *error {
	return nil
}

func (incorrect) Test11(interface{}, int, *int) int {
	return 0
}

func (incorrect) Test20(context.Context, a, *int) bool {
	return false
}

func (incorrect) Test21(context.Context, *a, *int) {
}

func (incorrect) Test31(context.Context, int, int) {
}

func (incorrect) Test32(context.Context, int, *a) {
}

var (
	co       correct
	inco     incorrect
	corrects = []interface{}{
		co.Test00,
		co.Test01,
		co.Test02,
		co.Test03,
		co.Test10,
		co.Test11,
		co.Test12,
		co.Test13,
		co.Test20,
		//co.Test21,
	}
	incorrects = []interface{}{
		inco.Test00,
		inco.Test01,
		inco.Test10,
		inco.Test11,
		inco.Test20,
		inco.Test21,
		inco.Test31,
		inco.Test32,
	}
)

func TestCheckInsCorrect(t *testing.T) {
	for i, f := range corrects {
		in0, in1, in2, err := checkIns(reflect.TypeOf(f))
		if err != nil {
			t.Fatalf("%d: checkIns: %v", i, err)
		}
		if in0 != reflect.TypeOf(f).In(0) {
			t.Fatalf("%d: in0: %s != %s", i, in0, reflect.TypeOf(f).In(0))
		}
		if in1 != reflect.TypeOf(f).In(1) {
			t.Fatalf("%d: in1: %s != %s", i, in1, reflect.TypeOf(f).In(1))
		}
		if in2 != reflect.TypeOf(f).In(2) {
			t.Fatalf("%d: in2: %s != %s", i, in2, reflect.TypeOf(f).In(2))
		}
		t.Log(in0, in1, in2)
	}
}

func TestCheckInsIncorrect(t *testing.T) {
	for i, f := range incorrects {
		_, _, _, err := checkIns(reflect.TypeOf(f))
		if err == nil {
			t.Fatalf("%d: checkIns return nil", i)
		} else {
			t.Logf("%d: checkIns: %v", i, err)
		}
	}
}

func TestCheckOutsCorrect(t *testing.T) {
	for i, f := range corrects {
		err := checkOuts(reflect.TypeOf(f))
		if err != nil {
			t.Fatalf("%d: checkOuts: %v", i, err)
		}
	}
}

func TestCheckOutsIncorrect(t *testing.T) {
	for i, f := range incorrects {
		err := checkOuts(reflect.TypeOf(f))
		if err == nil {
			t.Fatalf("%d: checkOuts return nil", i)
		} else {
			t.Logf("%d: checkOuts: %v", i, err)
		}
	}
}

func TestParseFunctionCorrect(t *testing.T) {
	for i, f := range corrects {
		fn, err := parseFunction(f)
		if err != nil {
			t.Fatalf("%d: parse function: %v", i, err)
		}
		if got, want := fn.value, reflect.ValueOf(f); got != want {
			t.Fatalf("%d: value: got %v, want %v", i, got, want)
		}
		if got, want := fn.args, reflect.TypeOf(f).In(1); got != want {
			t.Fatalf("%d: args: got %v, want %v", i, got, want)
		}
		if got, want := fn.reply, reflect.TypeOf(f).In(2); got != want {
			t.Fatalf("%d: reply: got %v, want %v", i, got, want)
		}
		t.Log(fn.value, fn.args, fn.reply)
	}
}

func TestParseFunctionIncorrect(t *testing.T) {
	for i, f := range incorrects {
		_, err := parseFunction(f)
		if err == nil {
			t.Fatalf("%d: parse function return nil", i)
		} else {
			t.Logf("%d: parse function: %v", i, err)
		}
	}
}

func TestFunctionCallReturnNil(t *testing.T) {
	for i, f := range corrects {
		fn, err := parseFunction(f)
		if err != nil {
			t.Fatalf("%d: parse function: %v", i, err)
		}
		args := fn.NewArgs()
		if fn.args.Kind() != reflect.Ptr {
			args = args.Elem()
		}
		err = fn.Call(context.Background(), args, fn.NewReply())
		if err != nil {
			t.Fatalf("%d: call: %v", i, err)
		}
	}
}

func TestFunctionCallReturnErr(t *testing.T) {
	f := func(context.Context, interface{}, interface{}) error { return errors.New("test error") }
	fn, err := parseFunction(f)
	if err != nil {
		t.Fatalf("parse function: %v", err)
	}
	args := fn.NewArgs()
	if fn.args.Kind() != reflect.Ptr {
		args = args.Elem()
	}
	err = fn.Call(context.Background(), args, fn.NewReply())
	if err == nil {
		t.Fatalf("call return nil")
	} else {
		t.Logf("call: %v", err)
	}
}

func TestFunctionCall(t *testing.T) {
	type Request struct {
		A, B int
	}
	type Response struct {
		C int
	}

	var calls int
	f := func(ctx context.Context, req *Request, resp *Response) error {
		calls++
		resp.C = req.A + req.B
		return nil
	}

	fn, err := parseFunction(f)
	if err != nil {
		t.Fatalf("parse function: %v", err)
	}

	var req Request
	var resp Response
	req.A = 1
	req.B = 2
	err = fn.Call(context.Background(), reflect.ValueOf(&req), reflect.ValueOf(&resp))
	if err != nil {
		t.Fatalf("call: %v", err)
	}
	if calls != 1 {
		t.Fatalf("calls(%d) != 1", calls)
	}
	if resp.C != req.A+req.B {
		t.Fatalf("C(%d) != A(%d) + B(%d)", resp.C, req.A, req.B)
	}
	t.Logf("C(%d) == A(%d) + B(%d)", resp.C, req.A, req.B)
}
