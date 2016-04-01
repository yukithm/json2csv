package main

import (
	"fmt"
	"strconv"
	"strings"
)

// Tokens is a sequence of token.
// Tokens contains only strings even if a token is an index of array.
type Tokens []string

// ParseJSONPointer parses JSON Pointer and return Tokens.
func ParseJSONPointer(pointer string) (Tokens, error) {
	pointer = strings.TrimSpace(pointer)
	if !strings.HasPrefix(pointer, "/") {
		return nil, fmt.Errorf("Invalid JSON Pointer %q", pointer)
	}

	parts := strings.Split(pointer[1:], "/")
	if len(parts) == 0 {
		return nil, fmt.Errorf("Invalid JSON Pointer %q", pointer)
	}

	tokens := make(Tokens, 0, len(parts))
	for _, part := range parts {
		tokens = append(tokens, unescapePart(part))
	}

	return tokens, nil
}

// DotNotation returns dot-notated representation.
func (t *Tokens) DotNotation(bracketIndex bool) string {
	if !bracketIndex {
		return t.join(".", false)
	}

	parts := make([]string, 0, len(*t))
	for _, token := range *t {
		if isInt(token) {
			// foo[0] style
			parts[len(parts)-1] += fmt.Sprintf("[%s]", token)
		} else {
			parts = append(parts, token)
		}
	}
	return strings.Join(parts, ".")
}

// JSONPointer returns JSON Pointer representation.
func (t *Tokens) JSONPointer() string {
	return "/" + t.join("/", true)
}

func (t *Tokens) join(sep string, escapeJSONPointer bool) string {
	parts := make([]string, 0, len(*t))
	for _, token := range *t {
		if escapeJSONPointer {
			token = escapePart(token)
		}
		parts = append(parts, token)
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

func isInt(value string) bool {
	_, err := strconv.Atoi(value)
	return err == nil
}
