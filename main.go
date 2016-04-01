package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"reflect"
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

	tokens := Tokens{"foo", "0", "bar", "foo/bar", "foo~bar", "foo\"bar", "foo\\bar"}
	fmt.Println("         JSONPointer:", tokens.JSONPointer())
	fmt.Println("         DotNotation:", tokens.DotNotation(false))
	fmt.Println("DotNotation(bracket):", tokens.DotNotation(true))

	k2, err := ParseJSONPointer(tokens.JSONPointer())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Parsed:", k2)
	if reflect.DeepEqual(tokens, k2) {
		fmt.Println("tokens == k2")
	}
}
