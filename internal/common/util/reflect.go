package util

import "reflect"

func Type(val any) string {
	if val == nil {
		return "nil"
	}

	typ := reflect.TypeOf(val)
	if typ.Kind() == reflect.Pointer {
		return "*" + Type(reflect.Indirect(reflect.ValueOf(val)).Interface())
	}

	if name := typ.Name(); name != "" {
		return name
	}

	return typ.Kind().String()
}
