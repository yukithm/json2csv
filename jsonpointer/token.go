package jsonpointer

import (
	"strconv"
	"strings"
)

// Token is each part of a JSON Pointer.
type Token string

// NewTokenFromEscaped returns a new Token from an escaped string.
func NewTokenFromEscaped(token string) Token {
	return Token(UnescapeTokenString(token))
}

// UnescapeTokenString returns unescaped representation of the token.
func UnescapeTokenString(token string) string {
	return strings.NewReplacer(
		"~1", "/",
		"~0", "~",
	).Replace(token)
}

// EscapedString returns escaped representation of the token.
func (t Token) EscapedString() string {
	return strings.NewReplacer(
		"~", "~0",
		"/", "~1",
	).Replace(string(t))
}

// IsInt returns true if the token is an integer like string.
func (t Token) IsInt() bool {
	_, err := strconv.Atoi(string(t))
	return err == nil
}

// IsIndex returns true if the token is an index like string.
func (t Token) IsIndex() bool {
	v := string(t)
	if _, err := strconv.Atoi(v); err != nil {
		return false
	}
	if v[0] == '+' || v[0] == '-' || (len(v) > 1 && v[0] == '0') {
		return false
	}

	return true
}
