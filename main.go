package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"reflect"
)

func main() {
	buf, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	var data interface{}
	if err := json.Unmarshal(buf, &data); err != nil {
		log.Fatal(err)
	}

	var results []keyValue
	v := valueOf(data)
	switch v.Kind() {
	case reflect.Map:
		result := flatten(v)
		results = append(results, result)
	case reflect.Slice:
		for i := 0; i < v.Len(); i++ {
			result := flatten(v.Index(i))
			results = append(results, result)
		}
	default:
		log.Fatal("Unsupported JSON structure.")
	}

	fmt.Println(results)

	js, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(js))
}

func valueOf(obj interface{}) reflect.Value {
	v := reflect.ValueOf(obj)
	for v.Kind() == reflect.Interface {
		v = v.Elem()
	}
	return v
}
