package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"json2csv"
	"json2csv/jsonpointer"
	"log"
	"os"

	"github.com/jessevdk/go-flags"
)

var options struct {
	HeaderStyle string `long:"header-style" choice:"jsonpointer" choice:"slash" choice:"dot" choice:"dot-bracket" default:"jsonpointer" description:"Header style"`
	Path        string `long:"path" description:"Target path (JSON Pointer) of the JSON content"`
	Transpose   bool   `long:"transpose" description:"Transponse rows and columns"`
}

var headerStyleTable = map[string]json2csv.KeyStyle{
	"jsonpointer": json2csv.JSONPointerStyle,
	"slash":       json2csv.SlashStyle,
	"dot":         json2csv.DotNotationStyle,
	"dot-bracket": json2csv.DotBracketStyle,
}

// USAGE for go-flags parser.
const USAGE = `[OPTION] [FILE]

Conver JSON FILE or STDIN to CSV.
`

func main() {
	// Hide timestamp because this is CLI application, so just print message for users.
	log.SetFlags(0)

	oparser := flags.NewParser(&options, flags.HelpFlag|flags.PassDoubleDash|flags.PassAfterNonOption)
	oparser.Usage = USAGE
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
	if len(args) > 0 && args[0] != "-" {
		data, err = readJSONFile(args[0])
	} else {
		data, err = readJSON(os.Stdin)
	}
	if err != nil {
		log.Fatal(err)
	}

	if options.Path != "" {
		data, err = jsonpointer.Get(data, options.Path)
		if err != nil {
			log.Fatal(err)
		}
	}

	results, err := json2csv.JSON2CSV(data)
	if err != nil {
		log.Fatal(err)
	}

	headerStyle := headerStyleTable[options.HeaderStyle]
	if err := printCSV(os.Stdout, results, headerStyle, options.Transpose); err != nil {
		log.Fatal(err)
	}
}

func readJSONFile(filename string) (interface{}, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return readJSON(f)
}

func readJSON(r io.Reader) (interface{}, error) {
	buf, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	var data interface{}
	if err := json.Unmarshal(buf, &data); err != nil {
		return nil, err
	}

	return data, nil
}

func printCSV(w io.Writer, results []json2csv.KeyValue, headerStyle json2csv.KeyStyle, transpose bool) error {
	csv := json2csv.NewCSVWriter(w)
	csv.HeaderStyle = headerStyle
	csv.Transpose = transpose
	if err := csv.WriteCSV(results); err != nil {
		return err
	}
	return nil
}
