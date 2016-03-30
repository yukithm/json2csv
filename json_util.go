package main

import (
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"
)

// Keys is a equence of keys/indexes.
// keys contains only string and int.
type Keys []interface{}

// ParseJSONPointer parses JSON Pointer and return Keys.
func ParseJSONPointer(pointer string) (Keys, error) {
	pointer = strings.TrimSpace(pointer)
	if !strings.HasPrefix(pointer, "/") {
		return nil, fmt.Errorf("Invalid JSON Pointer %q", pointer)
	}

	parts := strings.Split(pointer[1:], "/")
	if len(parts) == 0 {
		return nil, fmt.Errorf("Invalid JSON Pointer %q", pointer)
	}

	keys := make(Keys, 0, len(parts))
	for _, part := range parts {
		if index, err := strconv.Atoi(part); err == nil {
			keys = append(keys, index)
		} else {
			keys = append(keys, unescapePart(part))
		}
	}

	return keys, nil
}

// DotNotation returns dot-notated representation.
func (keys *Keys) DotNotation(bracketIndex bool) string {
	if !bracketIndex {
		return keys.join(".", false)
	}

	parts := make([]string, 0, len(*keys))
	for _, k := range *keys {
		switch key := k.(type) {
		case string:
			parts = append(parts, key)
		case int, uint:
			parts[len(parts)-1] += fmt.Sprintf("[%d]", key)
		default:
			v := reflect.ValueOf(k)
			log.Fatalf("Unsupported key type %q", v.Type())
		}
	}
	return strings.Join(parts, ".")
}

// JSONPointer returns JSON Pointer representation.
func (keys *Keys) JSONPointer() string {
	return "/" + keys.join("/", true)
}

func (keys *Keys) join(sep string, escapeJSONPointer bool) string {
	parts := make([]string, 0, len(*keys))
	for _, k := range *keys {
		var part string
		switch key := k.(type) {
		case string:
			part = key
		case int:
			part = strconv.Itoa(key)
		case uint:
			part = strconv.FormatUint(uint64(key), 10)
		default:
			log.Fatalf("Unsupported key type %q", reflect.TypeOf(k))
		}
		if escapeJSONPointer {
			part = escapePart(part)
		}
		parts = append(parts, part)
	}
	return strings.Join(parts, sep)
}

func escapePart(part string) string {
	r := strings.NewReplacer(
		"~", "~0",
		"/", "~1",
	)
	return r.Replace(part)
}

func unescapePart(part string) string {
	r := strings.NewReplacer(
		"~1", "/",
		"~0", "~",
	)
	return r.Replace(part)
}
