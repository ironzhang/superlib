package testutil

import (
	"reflect"
)

type CompareFilterFunc func(pkg, typ, field string) bool

type visit struct {
	a1  uintptr
	a2  uintptr
	typ reflect.Type
}

func compare(v1, v2 reflect.Value, visited map[visit]bool, depth int, filter CompareFilterFunc) bool {
	if !v1.IsValid() || !v2.IsValid() {
		return v1.IsValid() == v2.IsValid()
	}
	if v1.Type() != v2.Type() {
		return false
	}

	hard := func(v1, v2 reflect.Value) bool {
		switch v1.Kind() {
		case reflect.Ptr, reflect.Map, reflect.Slice, reflect.Interface:
			return !v1.IsNil() && !v2.IsNil()
		default:
			return false
		}
	}

	if hard(v1, v2) {
		addr1 := v1.Pointer()
		addr2 := v2.Pointer()
		if addr1 > addr2 {
			addr1, addr2 = addr2, addr1
		}

		typ := v1.Type()
		v := visit{addr1, addr2, typ}
		if visited[v] {
			return true
		}
		visited[v] = true
	}

	switch v1.Kind() {
	case reflect.Array:
		for i := 0; i < v1.Len(); i++ {
			if !compare(v1.Index(i), v2.Index(i), visited, depth+1, filter) {
				return false
			}
		}
		return true
	case reflect.Slice:
		if v1.IsNil() != v2.IsNil() {
			return false
		}
		if v1.Len() != v2.Len() {
			return false
		}
		if v1.Pointer() == v2.Pointer() {
			return true
		}
		for i := 0; i < v1.Len(); i++ {
			if !compare(v1.Index(i), v2.Index(i), visited, depth+1, filter) {
				return false
			}
		}
		return true

	case reflect.Interface:
		if v1.IsNil() || v2.IsNil() {
			return v1.IsNil() == v2.IsNil()
		}
		return compare(v1.Elem(), v2.Elem(), visited, depth+1, filter)

	case reflect.Ptr:
		if v1.Pointer() == v2.Pointer() {
			return true
		}
		return compare(v1.Elem(), v2.Elem(), visited, depth+1, filter)

	case reflect.Struct:
		typ := v1.Type()
		for i, n := 0, v1.NumField(); i < n; i++ {
			if filter != nil && filter(typ.PkgPath(), typ.Name(), typ.Field(i).Name) {
				continue
			}
			if !compare(v1.Field(i), v2.Field(i), visited, depth+1, filter) {
				return false
			}
		}
		return true

	case reflect.Map:
		if v1.IsNil() != v2.IsNil() {
			return false
		}
		if v1.Len() != v2.Len() {
			return false
		}
		if v1.Pointer() == v2.Pointer() {
			return true
		}
		for _, k := range v1.MapKeys() {
			val1 := v1.MapIndex(k)
			val2 := v2.MapIndex(k)
			if !val1.IsValid() || !val2.IsValid() || compare(val1, val2, visited, depth+1, filter) {
				return false
			}
		}
		return true
	case reflect.Func:
		if v1.IsNil() && v2.IsNil() {
			return true
		}
		return false
	default:
		return baseCompare(v1, v2)
	}
}

func baseCompare(v1, v2 reflect.Value) bool {
	switch k := v1.Kind(); k {
	case reflect.Bool:
		return v1.Bool() == v2.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v1.Int() == v2.Int()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v1.Uint() == v2.Uint()
	case reflect.Float32, reflect.Float64:
		return v1.Float() == v2.Float()
	case reflect.Complex64, reflect.Complex128:
		return v1.Complex() == v2.Complex()
	case reflect.String:
		return v1.String() == v2.String()
	default:
		return v1.Interface() == v2.Interface()
	}
}

func Compare(x, y interface{}, filter CompareFilterFunc) bool {
	if x == nil || y == nil {
		return x == y
	}
	v1 := reflect.ValueOf(x)
	v2 := reflect.ValueOf(y)
	if v1.Type() != v2.Type() {
		return false
	}
	return compare(v1, v2, make(map[visit]bool), 0, filter)
}
