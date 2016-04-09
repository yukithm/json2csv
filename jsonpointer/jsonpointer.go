// Package jsonpointer implements representations for JSON Pointer and tokens.
package jsonpointer

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// JSONPointer is a sequence of Token.
type JSONPointer []Token

// New parses a pointer string and creates a new JSONPointer.
func New(pointer string) (JSONPointer, error) {
	pointer = strings.TrimSpace(pointer)
	if pointer == "" {
		return JSONPointer{}, nil
	}

	if !strings.HasPrefix(pointer, "/") {
		return nil, fmt.Errorf("Invalid JSON Pointer %q", pointer)
	}

	tokens := strings.Split(pointer[1:], "/")
	if len(tokens) == 0 {
		return nil, fmt.Errorf("Invalid JSON Pointer %q", pointer)
	}

	jp := make(JSONPointer, 0, len(tokens))
	for _, token := range tokens {
		jp = append(jp, NewTokenFromEscaped(token))
	}

	return jp, nil
}

// Get retrieves a value from the obj.
func Get(obj interface{}, pointer string) (interface{}, error) {
	p, err := New(pointer)
	if err != nil {
		return nil, err
	}
	return p.Get(obj)
}

// Len returns the length of tokens.
func (p *JSONPointer) Len() int {
	return len(*p)
}

// Append appends the token.
func (p *JSONPointer) Append(token Token) *JSONPointer {
	*p = append(*p, token)
	return p
}

// AppendString appends the token.
func (p *JSONPointer) AppendString(token string) *JSONPointer {
	*p = append(*p, Token(token))
	return p
}

// Pop removes last token and return.
func (p *JSONPointer) Pop() Token {
	if p.Len() == 0 {
		return Token("")
	}

	// t, *p := (*p)[len(*p)-1], (*p)[:len(*p)-1]
	t := (*p)[len(*p)-1]
	*p = (*p)[:len(*p)-1]
	return t
}

// Clone returns a duplicate of the JSONPointer.
func (p JSONPointer) Clone() JSONPointer {
	if p.Len() == 0 {
		return JSONPointer{}
	}

	obj := make(JSONPointer, len(p))
	copy(obj, p)
	return obj
}

// Strings returns an array of each token string.
func (p JSONPointer) Strings() []string {
	tokens := make([]string, 0, len(p))
	for _, token := range p {
		tokens = append(tokens, string(token))
	}
	return tokens
}

// EscapedStrings returns an array of each token string that is escaped.
func (p JSONPointer) EscapedStrings() []string {
	tokens := make([]string, 0, len(p))
	for _, token := range p {
		tokens = append(tokens, token.EscapedString())
	}
	return tokens
}

// String returns JSON Pointer representation.
func (p JSONPointer) String() string {
	s := p.EscapedStrings()
	if len(s) == 0 {
		return ""
	}

	return "/" + strings.Join(p.EscapedStrings(), "/")
}

// DotNotation returns dot-notated representation.
func (p JSONPointer) DotNotation(bracketIndex bool) string {
	if !bracketIndex {
		return strings.Join(p.Strings(), ".")
	}

	tokens := make([]string, 0, len(p))
	for _, token := range p {
		if token.IsIndex() {
			// foo[0] style
			tokens[len(tokens)-1] += fmt.Sprintf("[%s]", token)
		} else {
			tokens = append(tokens, string(token))
		}
	}
	return strings.Join(tokens, ".")
}

// Get retrieves a value from the obj.
func (p JSONPointer) Get(obj interface{}) (value interface{}, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("Invalid JSON Pointer %q", p)
		}
	}()

	v := valueOf(obj)
	for i := 0; i < p.Len(); i++ {
		token := string(p[i])
		switch v.Kind() {
		case reflect.Map:
			v = v.MapIndex(reflect.ValueOf(token))
			v = valueOf(v)
		case reflect.Slice:
			if index, e := strconv.Atoi(token); e == nil {
				v = v.Index(index)
				v = valueOf(v)
			} else {
				return nil, fmt.Errorf("Invalid JSON Pointer %q", p)
			}
		}
	}

	return v.Interface(), nil
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
