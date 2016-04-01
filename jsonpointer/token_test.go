package jsonpointer

import "testing"

var testNewTokenFromEscapedCases = []struct {
	token    string
	expected string
}{
	{`foo`, `foo`},
	{`foo~0bar`, `foo~bar`},
	{`foo~1bar`, `foo/bar`},
	{`foo~0bar~0baz~1qux`, `foo~bar~baz/qux`},
}

func TestNewTokenFromEscaped(t *testing.T) {
	for caseIndex, testCase := range testNewTokenFromEscapedCases {
		actual := string(NewTokenFromEscaped(testCase.token))
		if actual != testCase.expected {
			t.Errorf("%d: Expected %v, but %v", caseIndex, testCase.expected, actual)
		}
	}
}

var testUnescapeTokenStringCases = []struct {
	token    string
	expected string
}{
	{`foo`, `foo`},
	{`foo~0bar`, `foo~bar`},
	{`foo~1bar`, `foo/bar`},
	{`foo~0bar~0baz~1qux`, `foo~bar~baz/qux`},
}

func TestUnescapeTokenString(t *testing.T) {
	for caseIndex, testCase := range testUnescapeTokenStringCases {
		actual := UnescapeTokenString(testCase.token)
		if actual != testCase.expected {
			t.Errorf("%d: Expected %v, but %v", caseIndex, testCase.expected, actual)
		}
	}
}

var testEscapedStringCases = []struct {
	token    string
	expected string
}{
	{`foo`, `foo`},
	{`foo~bar`, `foo~0bar`},
	{`foo/bar`, `foo~1bar`},
	{`foo~bar~baz/qux`, `foo~0bar~0baz~1qux`},
}

func TestEscapedString(t *testing.T) {
	for caseIndex, testCase := range testEscapedStringCases {
		actual := Token(testCase.token).EscapedString()
		if actual != testCase.expected {
			t.Errorf("%d: Expected %v, but %v", caseIndex, testCase.expected, actual)
		}
	}
}

var testIsIntCases = []struct {
	token    string
	expected bool
}{
	{"0", true},
	{"1", true},
	{"999", true},
	{"001", true},
	{"-1", true},
	{"+1", true},
	{"1.3", false},
	{"1a", false},
	{"a1", false},
	{"foo", false},
	{"", false},
}

func TestIsInt(t *testing.T) {
	for caseIndex, testCase := range testIsIntCases {
		actual := Token(testCase.token).IsInt()
		if actual != testCase.expected {
			t.Errorf("%d: Expected %v, but %v", caseIndex, testCase.expected, actual)
		}
	}
}

var testIsIndexCases = []struct {
	token    string
	expected bool
}{
	{"0", true},
	{"1", true},
	{"999", true},
	{"001", false},
	{"1.3", false},
	{"-1", false},
	{"+1", false},
	{"1a", false},
	{"a1", false},
	{"foo", false},
	{"", false},
}

func TestIsIndex(t *testing.T) {
	for caseIndex, testCase := range testIsIndexCases {
		actual := Token(testCase.token).IsIndex()
		if actual != testCase.expected {
			t.Errorf("%d: Expected %v, but %v", caseIndex, testCase.expected, actual)
		}
	}
}
