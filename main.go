package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"reflect"

	"json2csv/jsonpointer"
)

func main() {
	buf, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	var data interface{}
	if err := json.Unmarshal(buf, &data); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%v\n", data)

	pointer := jsonpointer.JSONPointer{"foo", "0", "bar", "foo/bar", "foo~bar", "foo\"bar", "foo\\bar"}
	fmt.Println("         JSONPointer:", pointer)
	fmt.Println("         DotNotation:", pointer.DotNotation(false))
	fmt.Println("DotNotation(bracket):", pointer.DotNotation(true))

	pointer2, err := jsonpointer.NewJSONPointer(pointer.String())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Parsed:", pointer2)
	if reflect.DeepEqual(pointer, pointer2) {
		fmt.Println("pointer == pointer2")
	}
}
