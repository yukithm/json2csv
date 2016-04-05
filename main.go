package main

import (
	"encoding/json"
	"io"
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

	if err := printCSV(os.Stdout, results); err != nil {
		log.Fatal(err)
	}
}

func printCSV(w io.Writer, results []keyValue) error {
	csv := NewCSVWriter(w)
	csv.style = JSONPointerStyle
	if err := csv.WriteCSV(results); err != nil {
		return err
	}
	return nil
}
