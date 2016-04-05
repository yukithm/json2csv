package main

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
)

func debugDump(results []keyValue) {
	js, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	debugPrint(string(js))
}

func debugPrint(obj interface{}) {
	fmt.Println(obj)
}

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
