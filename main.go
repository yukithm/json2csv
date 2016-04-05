package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"os"
	"reflect"

	"github.com/jessevdk/go-flags"
)

var options struct {
	HeaderStyle string `long:"header-style" choice:"jsonpointer" choice:"slash" choice:"dot" choice:"dot-bracket" default:"jsonpointer" description:"Header style"`
}

var headerStyleTable = map[string]keyStyle{
	"jsonpointer": JSONPointerStyle,
	"slash":       SlashStyle,
	"dot":         DotNotationStyle,
	"dot-bracket": DotBracketStyle,
}

func main() {
	oparser := flags.NewParser(&options, flags.HelpFlag|flags.PassDoubleDash|flags.PassAfterNonOption)
	args, err := oparser.Parse()
	if err != nil {
		if e, ok := err.(*flags.Error); ok && e.Type == flags.ErrHelp {
			os.Stdout.WriteString(e.Message + "\n")
			os.Exit(0)
		} else {
			log.Fatal(err)
		}
	}

	var data interface{}
	if len(args) > 0 {
		if data, err = readJSONFile(args[0]); err != nil {
			log.Fatal(err)
		}
	} else {
		if data, err = readJSON(); err != nil {
			log.Fatal(err)
		}
	}

	var results []keyValue
	v := valueOf(data)
	switch v.Kind() {
	case reflect.Map:
		result := flatten(v)
		results = append(results, result)
	case reflect.Slice:
		for i := 0; i < v.Len(); i++ {
			result := flatten(v.Index(i))
			results = append(results, result)
		}
	default:
		log.Fatal("Unsupported JSON structure.")
	}

	headerStyle := headerStyleTable[options.HeaderStyle]
	if err := printCSV(os.Stdout, results, headerStyle); err != nil {
		log.Fatal(err)
	}
}

func readJSON() (interface{}, error) {
	buf, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		return nil, err
	}

	var data interface{}
	if err := json.Unmarshal(buf, &data); err != nil {
		return nil, err
	}

	return data, nil
}

func readJSONFile(filename string) (interface{}, error) {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var data interface{}
	if err := json.Unmarshal(buf, &data); err != nil {
		return nil, err
	}

	return data, nil
}

func printCSV(w io.Writer, results []keyValue, headerStyle keyStyle) error {
	csv := NewCSVWriter(w)
	csv.headerStyle = headerStyle
	if err := csv.WriteCSV(results); err != nil {
		return err
	}
	return nil
}
