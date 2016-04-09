package json2csv

import (
	"fmt"
	"reflect"
)

func valueOf(obj interface{}) reflect.Value {
	v, ok := obj.(reflect.Value)
	if !ok {
		v = reflect.ValueOf(obj)
	}

	for v.Kind() == reflect.Interface && !v.IsNil() {
		v = v.Elem()
	}
	return v
}

func toString(obj interface{}) string {
	return fmt.Sprintf("%v", obj)
}
