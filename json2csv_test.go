package json2csv

import (
	"bytes"
	"encoding/json"
	"reflect"
	"testing"
)

// Decode JSON with UseNumber option.
func json2obj(jsonstr string) (interface{}, error) {
	r := bytes.NewReader([]byte(jsonstr))
	d := json.NewDecoder(r)
	d.UseNumber()
	var obj interface{}
	if err := d.Decode(&obj); err != nil {
		return nil, err
	}
	return obj, nil
}

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
			{"/id": json.Number("1"), "/name": "foo"},
			{"/id": json.Number("2"), "/name": "bar"},
		},
		``,
	},
	{
		`[
			{"id": 1, "name/a": "foo"},
			{"id": 2, "name~b": "bar"}
		]`,
		[]KeyValue{
			{"/id": json.Number("1"), "/name~1a": "foo"},
			{"/id": json.Number("2"), "/name~0b": "bar"},
		},
		``,
	},
	{
		`[
			{"id":1, "values":["a", "b"]},
			{"id":2, "values":["x"]}
		]`,
		[]KeyValue{
			{"/id": json.Number("1"), "/values/0": "a", "/values/1": "b"},
			{"/id": json.Number("2"), "/values/0": "x"},
		},
		``,
	},
	{
		`[
			{"id":1, "values":[]},
			{"id":2, "values":["x"]}
		]`,
		[]KeyValue{
			{"/id": json.Number("1")},
			{"/id": json.Number("2"), "/values/0": "x"},
		},
		``,
	},
	{
		`[
			{"id":1, "values":{}},
			{"id":2, "values":["x"]}
		]`,
		[]KeyValue{
			{"/id": json.Number("1")},
			{"/id": json.Number("2"), "/values/0": "x"},
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
			{"/id": json.Number("123"), "/values/0/foo": "FOO", "/values/1/bar": "BAR"},
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
	{
		`{"large_int_value": 146163870300}`,
		[]KeyValue{{"/large_int_value": json.Number("146163870300")}},
		``,
	},
	{
		`{"float_value": 146163870.300}`,
		[]KeyValue{{"/float_value": json.Number("146163870.300")}},
		``,
	},
	{`"foo"`, nil, `Unsupported JSON structure.`},
	{`123`, nil, `Unsupported JSON structure.`},
	{`true`, nil, `Unsupported JSON structure.`},
}

func TestJSON2CSV(t *testing.T) {
	for caseIndex, testCase := range testJSON2CSVCases {
		obj, err := json2obj(testCase.json)
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
