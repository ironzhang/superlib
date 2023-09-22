package bind

import (
	"errors"
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

func BindQuery(a interface{}, tag string) (url.Values, error) {
	typ := reflect.TypeOf(a)
	val := reflect.ValueOf(a)

	if typ.Kind() != reflect.Struct {
		return nil, errors.New("binding element must be a struct")
	}

	var err error
	values := make(url.Values)
	for i := 0; i < typ.NumField(); i++ {
		tfield := typ.Field(i)
		if !isExported(tfield.Name) {
			continue
		}
		vfield := val.Field(i)
		name := parseTagName(tfield, tag)
		if name == "" {
			continue
		}
		if err = setValue(values, name, vfield); err != nil {
			return nil, err
		}
	}
	return values, nil
}

func isExported(name string) bool {
	r, _ := utf8.DecodeRuneInString(name)
	return unicode.IsUpper(r)
}

func addBoolValue(values url.Values, key string, rv reflect.Value) {
	values.Add(key, strconv.FormatBool(rv.Bool()))
}

func addFloatValue(values url.Values, key string, rv reflect.Value) {
	values.Add(key, strconv.FormatFloat(rv.Float(), 'f', -1, 64))
}

func addIntValue(values url.Values, key string, rv reflect.Value) {
	values.Add(key, strconv.FormatInt(rv.Int(), 10))
}

func addUintValue(values url.Values, key string, rv reflect.Value) {
	values.Add(key, strconv.FormatUint(rv.Uint(), 10))
}

func addValue(values url.Values, key string, rv reflect.Value) error {
	switch k := rv.Kind(); k {
	case reflect.Ptr:
		return addValue(values, key, rv.Elem())
	case reflect.String:
		values.Add(key, rv.String())
	case reflect.Bool:
		addBoolValue(values, key, rv)
	case reflect.Float32, reflect.Float64:
		addFloatValue(values, key, rv)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		addIntValue(values, key, rv)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		addUintValue(values, key, rv)
	default:
		return fmt.Errorf("unsupported type %v", k)
	}
	return nil
}

func setValue(values url.Values, key string, rv reflect.Value) (err error) {
	if k := rv.Kind(); k == reflect.Slice || k == reflect.Array {
		for i, n := 0, rv.Len(); i < n; i++ {
			if err = addValue(values, key, rv.Index(i)); err != nil {
				return err
			}
		}
		return nil
	}
	return addValue(values, key, rv)
}

func parseTagName(tfield reflect.StructField, tagkey string) string {
	if len(tagkey) == 0 {
		return tfield.Name
	}
	tag := tfield.Tag.Get(tagkey)

	// 找不到tag
	if tag == "" {
		return tfield.Name
	}

	//根据","对tag进行切分
	fields := strings.Split(tag, ",")
	if len(fields) <= 0 {
		return tfield.Name
	}
	name := fields[0]
	if name == "" {
		return tfield.Name
	}

	//不处理
	if name == "-" {
		return ""
	}
	return name
}
