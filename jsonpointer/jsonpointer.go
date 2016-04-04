// Package jsonpointer implements representations for JSON Pointer and tokens.
package jsonpointer

import (
	"fmt"
	"strings"
)

// JSONPointer is a sequence of Token.
type JSONPointer []Token

// NewJSONPointer parses a pointer string and creates a new JSONPointer.
func NewJSONPointer(pointer string) (JSONPointer, error) {
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
