package main

import (
	"encoding/csv"
	"io"
	"json2csv/jsonpointer"
	"log"
	"sort"
)

type keyStyle uint

// Header style
const (
	JSONPointerStyle keyStyle = iota
	SlashStyle
	DotNotationStyle
	DotBracketStyle
)

// CSVWriter writes CSV data.
type CSVWriter struct {
	*csv.Writer
	headerStyle keyStyle
}

// NewCSVWriter returns new CSVWriter with JSONPointerStyle.
func NewCSVWriter(w io.Writer) *CSVWriter {
	return &CSVWriter{
		csv.NewWriter(w),
		JSONPointerStyle,
	}
}

// WriteCSV writes CSV data.
func (w *CSVWriter) WriteCSV(results []keyValue) error {
	pts := allPointers(results)
	sort.Sort(pts)
	keys := pts.Strings()

	var header []string
	switch w.headerStyle {
	case JSONPointerStyle:
		header = keys
	case SlashStyle:
		header = pts.Slashes()
	case DotNotationStyle:
		header = pts.DotNotations(false)
	case DotBracketStyle:
		header = pts.DotNotations(true)
	}

	if err := w.Write(header); err != nil {
		return err
	}

	for _, result := range results {
		record := toRecord(result, keys)
		if err := w.Write(record); err != nil {
			return err
		}
	}

	w.Flush()
	if err := w.Error(); err != nil {
		return err
	}

	return nil
}

func allPointers(results []keyValue) (pointers pointers) {
	set := make(map[string]bool, 0)
	for _, result := range results {
		for _, key := range result.Keys() {
			if !set[key] {
				set[key] = true
				pointer, err := jsonpointer.NewJSONPointer(key)
				if err != nil {
					log.Fatal(err)
				}
				pointers = append(pointers, pointer)
			}
		}
	}
	return
}

func toRecord(kv keyValue, keys []string) []string {
	record := make([]string, 0, len(keys))
	for _, key := range keys {
		if value, ok := kv[key]; ok {
			record = append(record, toString(value))
		} else {
			record = append(record, "")
		}
	}
	return record
}
