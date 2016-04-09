package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"json2csv"
	"json2csv/jsonpointer"
	"log"
	"os"

	"github.com/codegangsta/cli"
)

// VERSION is the version number of this application.
const VERSION = "0.1.0"

var headerStyleTable = map[string]json2csv.KeyStyle{
	"jsonpointer": json2csv.JSONPointerStyle,
	"slash":       json2csv.SlashStyle,
	"dot":         json2csv.DotNotationStyle,
	"dot-bracket": json2csv.DotBracketStyle,
}

func main() {
	// Hide timestamp because this is CLI application, so just print message for users.
	log.SetFlags(0)

	app := cli.NewApp()
	app.Name = "json2csv"
	app.Version = VERSION
	app.Usage = "convert JSON to CSV"
	app.HideHelp = true
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "header-style",
			Value: "jsonpointer",
			Usage: "header style (jsonpointer, slash, dot, dot-bracket)",
		},
		cli.StringFlag{
			Name:  "path",
			Usage: "target path (JSON Pointer) of the content",
		},
		cli.BoolFlag{
			Name:  "transpose",
			Usage: "transpose rows and columns",
		},
		cli.HelpFlag,
	}

	app.Before = func(c *cli.Context) error {
		if _, ok := headerStyleTable[c.String("header-style")]; !ok {
			return fmt.Errorf("Invalid --header-style value %q", c.String("header-style"))
		}
		return nil
	}

	app.Action = func(c *cli.Context) {
		if c.Bool("help") {
			cli.ShowAppHelp(c)
			return
		}
		mainAction(c)
	}

	app.RunAndExitOnError()
}

func mainAction(c *cli.Context) {
	var data interface{}
	var err error
	if c.NArg() > 0 && c.Args()[0] != "-" {
		data, err = readJSONFile(c.Args()[0])
	} else {
		data, err = readJSON(os.Stdin)
	}
	if err != nil {
		log.Fatal(err)
	}

	if c.String("path") != "" {
		data, err = jsonpointer.Get(data, c.String("path"))
		if err != nil {
			log.Fatal(err)
		}
	}

	results, err := json2csv.JSON2CSV(data)
	if err != nil {
		log.Fatal(err)
	}

	headerStyle := headerStyleTable[c.String("header-style")]
	err = printCSV(os.Stdout, results, headerStyle, c.Bool("transpose"))
	if err != nil {
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
