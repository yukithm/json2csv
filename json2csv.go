// Package json2csv provides JSON to CSV functions.
package json2csv

import (
	"errors"
	"reflect"
)

// JSON2CSV converts JSON to CSV.
func JSON2CSV(data interface{}) ([]KeyValue, error) {
	results := []KeyValue{}
	v := valueOf(data)
	switch v.Kind() {
	case reflect.Map:
		if v.Len() > 0 {
			result, err := flatten(v)
			if err != nil {
				return nil, err
			}
			results = append(results, result)
		}
	case reflect.Slice:
		if isObjectArray(v) {
			for i := 0; i < v.Len(); i++ {
				result, err := flatten(v.Index(i))
				if err != nil {
					return nil, err
				}
				results = append(results, result)
			}
		} else if v.Len() > 0 {
			result, err := flatten(v)
			if err != nil {
				return nil, err
			}
			if result != nil {
				results = append(results, result)
			}
		}
	default:
		return nil, errors.New("Unsupported JSON structure.")
	}

	return results, nil
}

func isObjectArray(obj interface{}) bool {
	value := valueOf(obj)
	if value.Kind() != reflect.Slice {
		return false
	}

	len := value.Len()
	if len == 0 {
		return false
	}
	for i := 0; i < len; i++ {
		if valueOf(value.Index(i)).Kind() != reflect.Map {
			return false
		}
	}

	return true
}
