package echorpc

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"unicode"
	"unicode/utf8"
)

var (
	typeOfError        = reflect.TypeOf((*error)(nil)).Elem()
	typeOfContext      = reflect.TypeOf((*context.Context)(nil)).Elem()
	typeOfNilInterface = reflect.TypeOf((*interface{})(nil)).Elem()
)

// Is this an exported - upper case - name?
func isExported(name string) bool {
	r, _ := utf8.DecodeRuneInString(name)
	return unicode.IsUpper(r)
}

// Is this type exported or a builtin?
func isExportedOrBuiltinType(t reflect.Type) bool {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return isExported(t.Name()) || t.PkgPath() == ""
}

// Function needs 3 ins: context.Context, *args, *reply.
func checkIns(ftype reflect.Type) (in0, in1, in2 reflect.Type, err error) {
	if ftype.NumIn() != 3 {
		err = fmt.Errorf("func has wrong number of ins: %d", ftype.NumIn())
		return
	}
	in0 = ftype.In(0)
	if !in0.Implements(typeOfContext) {
		err = fmt.Errorf("func context type dose not implement context.Context: %s", in0)
		return
	}
	in1 = ftype.In(1)
	if !isExportedOrBuiltinType(in1) {
		err = fmt.Errorf("func args type is not exported: %s", in1)
		return
	}
	in2 = ftype.In(2)
	if in2.Kind() != reflect.Ptr && in2 != typeOfNilInterface {
		err = fmt.Errorf("func reply type is not a pointer or interface{}: %s", in2)
		return
	}
	if !isExportedOrBuiltinType(in2) {
		err = fmt.Errorf("func reply type is not exported: %s", in2)
		return
	}
	return
}

// The return type of the function must be error.
func checkOuts(ftype reflect.Type) error {
	if ftype.NumOut() != 1 {
		return fmt.Errorf("func has wrong number of outs: %d", ftype.NumOut())
	}
	if out0 := ftype.Out(0); out0 != typeOfError {
		return fmt.Errorf("func returns %s not an error", out0)
	}
	return nil
}

type function struct {
	value reflect.Value
	args  reflect.Type
	reply reflect.Type
}

func parseFunction(f interface{}) (*function, error) {
	value := reflect.ValueOf(f)
	ftype := reflect.TypeOf(f)
	if ftype.Kind() != reflect.Func {
		return nil, errors.New("f is not a func")
	}
	_, args, reply, err := checkIns(ftype)
	if err != nil {
		return nil, err
	}
	if err = checkOuts(ftype); err != nil {
		return nil, err
	}
	return &function{value: value, args: args, reply: reply}, nil
}

func (p *function) NewArgs() reflect.Value {
	if p.args.Kind() == reflect.Ptr {
		return reflect.New(p.args.Elem())
	}
	return reflect.New(p.args)
}

func (p *function) NewReply() reflect.Value {
	if isNilInterface(p.reply) {
		return reflect.New(p.reply)
	}
	reply := reflect.New(p.reply.Elem())
	switch reply.Elem().Kind() {
	case reflect.Map:
		reply.Elem().Set(reflect.MakeMap(p.reply.Elem()))
	case reflect.Slice:
		reply.Elem().Set(reflect.MakeSlice(p.reply.Elem(), 0, 0))
	}
	return reply
}

func (p *function) Call(ctx context.Context, args, reply reflect.Value) error {
	in := []reflect.Value{reflect.ValueOf(ctx), args, reply}
	out := p.value.Call(in)
	ret := out[0].Interface()
	if err, ok := ret.(error); ok {
		return err
	}
	return nil
}

func isNilInterface(t reflect.Type) bool {
	return t == typeOfNilInterface
}
