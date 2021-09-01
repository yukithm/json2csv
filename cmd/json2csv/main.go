package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/yukithm/json2csv"
	"github.com/yukithm/json2csv/jsonpointer"

	"github.com/urfave/cli"
)

const (
	// ApplicationName is the name of this application.
	ApplicationName = "json2csv"
)

// injected by build process
var version = "unknown"

var headerStyleTable = map[string]json2csv.KeyStyle{
	"jsonpointer": json2csv.JSONPointerStyle,
	"slash":       json2csv.SlashStyle,
	"dot":         json2csv.DotNotationStyle,
	"dot-bracket": json2csv.DotBracketStyle,
}

func main() {
	// Hide timestamp because this is CLI application, so just print message for users.
	log.SetFlags(0)

	cli.AppHelpTemplate = `NAME:
   {{.Name}} - {{.Usage}}

USAGE:
   {{if .UsageText}}{{.UsageText}}{{else}}{{.HelpName}} {{if .Flags}}[OPTIONS]{{end}}{{if .Commands}} command [command options]{{end}} {{if .ArgsUsage}}{{.ArgsUsage}}{{else}}[arguments...]{{end}}{{end}}

   If no files are specified, JSON content is read from STDIN.
   {{if .Version}}{{if not .HideVersion}}
VERSION:
   {{.Version}}
   {{end}}{{end}}{{if len .Authors}}
AUTHOR(S):
   {{range .Authors}}{{ . }}{{end}}
   {{end}}{{if .Commands}}
COMMANDS:{{range .Categories}}{{if .Name}}
  {{.Name}}{{ ":" }}{{end}}{{range .Commands}}
    {{.Name}}{{with .ShortName}}, {{.}}{{end}}{{ "\t" }}{{.Usage}}{{end}}
{{end}}{{end}}{{if .Flags}}
OPTIONS:
   {{range .Flags}}{{.}}
   {{end}}{{end}}{{if .Copyright }}
COPYRIGHT:
   {{.Copyright}}
   {{end}}
`

	app := cli.NewApp()
	app.Name = ApplicationName
	app.Version = version
	app.Usage = "convert JSON to CSV"
	app.ArgsUsage = "[FILE]"
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
		log.Fatalf("readJSON failed, err=%v", err)
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
	if len(results) == 0 {
		return
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

func readJSON(r *os.File) (interface{}, error) {
	decoder := json.NewDecoder(r)
	decoder.UseNumber()

	var data interface{}
	err := decoder.Decode(&data)

	if err == nil {
		return data, nil
	}
	if jsonError, ok := err.(*json.SyntaxError); ok {
		jsonData, _ := ReadFile(r)

		line := 1 + strings.Count(string(jsonData[:jsonError.Offset]), "\n")
		column := 1 + int(jsonError.Offset) - (strings.LastIndex(string(jsonData[:jsonError.Offset]), "\n") + len("\n"))
		return nil, fmt.Errorf("cannot parse JSON schema due to a syntax error at line %d, column %d: %v", line, column, jsonError.Error())
	}
	return nil, fmt.Errorf("parse JSON schema failed, err=%v", err)
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

// ReadFile reads from a *os.File and returns the contents.
// took from stdlib os.ReadFile
// A successful call returns err == nil, not err == EOF.
// Because ReadFile reads the whole file, it does not treat an EOF from Read
// as an error to be reported.
func ReadFile(f *os.File) ([]byte, error) {
	var size int
	if info, err := f.Stat(); err == nil {
		size64 := info.Size()
		if int64(int(size64)) == size64 {
			size = int(size64)
		}
	}
	size++ // one byte for final read at EOF

	// If a file claims a small size, read at least 512 bytes.
	// In particular, files in Linux's /proc claim size 0 but
	// then do not work right if read in small pieces,
	// so an initial read of 1 byte would not work correctly.
	if size < 512 {
		size = 512
	}

	data := make([]byte, 0, size)
	for {
		if len(data) >= cap(data) {
			d := append(data[:cap(data)], 0)
			data = d[:len(data)]
		}
		n, err := f.Read(data[len(data):cap(data)])
		data = data[:len(data)+n]
		if err != nil {
			if err == io.EOF {
				err = nil
			}
			return data, err
		}
	}
}
