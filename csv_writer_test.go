package json2csv_test

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/yukithm/json2csv"
)

func TestKeyWithTrailingSpace(t *testing.T) {
	b := &bytes.Buffer{}
	wr := json2csv.NewCSVWriter(b)
	responses := []map[string]interface{}{
		{
			" A":  1,
			"B ":  "foo",
			"C  ": "FOO",
		},
		{
			" A":  2,
			"B ":  "bar",
			"C  ": "BAR",
		},
	}

	csvContent, err := json2csv.JSON2CSV(responses) // csvContent seems to be complete!
	if err != nil {
		t.Fatal(err)
	}
	wr.WriteCSV(csvContent)
	wr.Flush()

	got := b.String()
	want := `/ A,/B ,/C  
1,foo,FOO
2,bar,BAR
`

	if got != want {
		t.Errorf("Expected %v, but %v", want, got)
	}
}

func TestDotBracketRootHeaders(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		normal     string
		transposed string
	}{
		{
			name:       "root array",
			input:      `[1,2]`,
			normal:     "[0],[1]\n1,2\n",
			transposed: "[0],1\n[1],2\n",
		},
		{
			name:       "nested root arrays",
			input:      `[[1,2],[3,4]]`,
			normal:     "[0][0],[0][1],[1][0],[1][1]\n1,2,3,4\n",
			transposed: "[0][0],1\n[0][1],2\n[1][0],3\n[1][1],4\n",
		},
		{
			name:       "numeric root key",
			input:      `{"0":{"value":1}}`,
			normal:     "[0].value\n1\n",
			transposed: "[0].value,1\n",
		},
		{
			name:       "named root key",
			input:      `{"values":[1,2]}`,
			normal:     "values[0],values[1]\n1,2\n",
			transposed: "values[0],1\nvalues[1],2\n",
		},
	}
	for _, test := range tests {
		for _, transpose := range []bool{false, true} {
			name := test.name + "/normal"
			want := test.normal
			if transpose {
				name = test.name + "/transposed"
				want = test.transposed
			}
			t.Run(name, func(t *testing.T) {
				var data interface{}
				if err := json.Unmarshal([]byte(test.input), &data); err != nil {
					t.Fatal(err)
				}
				results, err := json2csv.JSON2CSV(data)
				if err != nil {
					t.Fatal(err)
				}
				var output bytes.Buffer
				writer := json2csv.NewCSVWriter(&output)
				writer.HeaderStyle = json2csv.DotBracketStyle
				writer.Transpose = transpose
				if err := writer.WriteCSV(results); err != nil {
					t.Fatal(err)
				}
				if got := output.String(); got != want {
					t.Errorf("Expected %q, but %q", want, got)
				}
			})
		}
	}
}
