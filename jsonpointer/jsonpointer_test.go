package jsonpointer

import (
	"reflect"
	"testing"
)

var testNewJSONPointerCases = []struct {
	pointer  string
	expected []Token
	err      string
}{
	{`/foo`, []Token{`foo`}, ``},
	{`/foo~0bar`, []Token{`foo~bar`}, ``},
	{`/foo~1bar`, []Token{`foo/bar`}, ``},
	{`/foo/bar`, []Token{`foo`, `bar`}, ``},
	{`/foo/0/bar`, []Token{`foo`, `0`, `bar`}, ``},
	{`foo`, nil, `Invalid JSON Pointer "foo"`},
	{``, nil, `Invalid JSON Pointer ""`},
}

func TestNewJSONPointer(t *testing.T) {
	for caseIndex, testCase := range testNewJSONPointerCases {
		pointer, err := NewJSONPointer(testCase.pointer)
		actual := []Token(pointer)
		if err != nil {
			if err.Error() != testCase.err {
				t.Errorf("%d: Expected %v, but %v", caseIndex, testCase.err, err)
			}
		} else if !reflect.DeepEqual(actual, testCase.expected) {
			t.Errorf("%d: Expected %v, but %v", caseIndex, testCase.expected, actual)
		}
	}
}

var testStringsCases = []struct {
	pointer  string
	expected []string
}{
	{`/foo`, []string{`foo`}},
	{`/foo~0bar`, []string{`foo~bar`}},
	{`/foo~1bar`, []string{`foo/bar`}},
	{`/foo/bar`, []string{`foo`, `bar`}},
	{`/foo/0/bar`, []string{`foo`, `0`, `bar`}},
}

func TestStrings(t *testing.T) {
	for caseIndex, testCase := range testStringsCases {
		pointer, err := NewJSONPointer(testCase.pointer)
		if err != nil {
			t.Error(err)
		}
		actual := pointer.Strings()
		if !reflect.DeepEqual(actual, testCase.expected) {
			t.Errorf("%d: Expected %v, but %v", caseIndex, testCase.expected, actual)
		}
	}
}

var testEscapedStringsCases = []struct {
	pointer  string
	expected []string
}{
	{`/foo`, []string{`foo`}},
	{`/foo~0bar`, []string{`foo~0bar`}},
	{`/foo~1bar`, []string{`foo~1bar`}},
	{`/foo/bar`, []string{`foo`, `bar`}},
	{`/foo/0/bar`, []string{`foo`, `0`, `bar`}},
}

func TestEscapedStrings(t *testing.T) {
	for caseIndex, testCase := range testEscapedStringsCases {
		pointer, err := NewJSONPointer(testCase.pointer)
		if err != nil {
			t.Error(err)
		}
		actual := pointer.EscapedStrings()
		if !reflect.DeepEqual(actual, testCase.expected) {
			t.Errorf("%d: Expected %v, but %v", caseIndex, testCase.expected, actual)
		}
	}
}

var testStringCases = []struct {
	pointer  string
	expected string
}{
	{`/foo`, `/foo`},
	{`/foo~0bar`, `/foo~0bar`},
	{`/foo~1bar`, `/foo~1bar`},
	{`/foo/bar`, `/foo/bar`},
	{`/foo/0/bar`, `/foo/0/bar`},
}

func TestString(t *testing.T) {
	for caseIndex, testCase := range testStringCases {
		pointer, err := NewJSONPointer(testCase.pointer)
		if err != nil {
			t.Error(err)
		}
		actual := pointer.String()
		if actual != testCase.expected {
			t.Errorf("%d: Expected %v, but %v", caseIndex, testCase.expected, actual)
		}
	}
}

var testDotNotationCases = []struct {
	pointer         string
	expected        string
	expectedBracket string
}{
	{`/foo`, `foo`, `foo`},
	{`/foo~0bar`, `foo~bar`, `foo~bar`},
	{`/foo~1bar`, `foo/bar`, `foo/bar`},
	{`/foo/bar`, `foo.bar`, `foo.bar`},
	{`/foo/0/bar`, `foo.0.bar`, `foo[0].bar`},
}

func TestDotNotation(t *testing.T) {
	for caseIndex, testCase := range testDotNotationCases {
		pointer, err := NewJSONPointer(testCase.pointer)
		if err != nil {
			t.Error(err)
		}
		actual := pointer.DotNotation(false)
		if actual != testCase.expected {
			t.Errorf("%d: Expected %v, but %v", caseIndex, testCase.expected, actual)
		}
		actual = pointer.DotNotation(true)
		if actual != testCase.expectedBracket {
			t.Errorf("%d: Expected %v, but %v", caseIndex, testCase.expectedBracket, actual)
		}
	}
}
