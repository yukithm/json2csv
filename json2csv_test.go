package json2csv

import (
	"encoding/json"
	"reflect"
	"testing"
)

var testJSON2CSVCases = []struct {
	json     string
	expected []KeyValue
	err      string
}{
	{
		`[
			{"id": 1, "name": "foo"},
			{"id": 2, "name": "bar"}
		]`,
		[]KeyValue{
			{"/id": 1.0, "/name": "foo"},
			{"/id": 2.0, "/name": "bar"},
		},
		``,
	},
	{
		`[
			{"id": 1, "name/a": "foo"},
			{"id": 2, "name~b": "bar"}
		]`,
		[]KeyValue{
			{"/id": 1.0, "/name~1a": "foo"},
			{"/id": 2.0, "/name~0b": "bar"},
		},
		``,
	},
	{
		`[
			{"id":1, "values":["a", "b"]},
			{"id":2, "values":["x"]}
		]`,
		[]KeyValue{
			{"/id": 1.0, "/values/0": "a", "/values/1": "b"},
			{"/id": 2.0, "/values/0": "x"},
		},
		``,
	},
	{
		`[
			{"id":1, "values":[]},
			{"id":2, "values":["x"]}
		]`,
		[]KeyValue{
			{"/id": 1.0},
			{"/id": 2.0, "/values/0": "x"},
		},
		``,
	},
	{
		`[
			{"id":1, "values":{}},
			{"id":2, "values":["x"]}
		]`,
		[]KeyValue{
			{"/id": 1.0},
			{"/id": 2.0, "/values/0": "x"},
		},
		``,
	},
	{
		`{
			"id": 123,
			"values": [
				{"foo": "FOO"},
				{"bar": "BAR"}
			]
		}`,
		[]KeyValue{
			{"/id": 123.0, "/values/0/foo": "FOO", "/values/1/bar": "BAR"},
		},
		``,
	},
	{
		`[]`,
		[]KeyValue{},
		``,
	},
	{
		`{}`,
		[]KeyValue{},
		``,
	},
	{`"foo"`, nil, `Unsupported JSON structure.`},
	{`123`, nil, `Unsupported JSON structure.`},
	{`true`, nil, `Unsupported JSON structure.`},
}

func TestJSON2CSV(t *testing.T) {
	for caseIndex, testCase := range testJSON2CSVCases {
		var obj interface{}
		err := json.Unmarshal([]byte(testCase.json), &obj)
		if err != nil {
			t.Fatal(err)
		}

		actual, err := JSON2CSV(obj)
		if err != nil {
			if err.Error() != testCase.err {
				t.Errorf("%d: Expected %v, but %v", caseIndex, testCase.err, err)
			}
		} else if !reflect.DeepEqual(testCase.expected, actual) {
			t.Errorf("%d: Expected %#v, but %#v", caseIndex, testCase.expected, actual)
		}
	}
}
