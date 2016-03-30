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

	keys := Keys{"foo", 0, "bar", "foo/bar", "foo~bar", "foo\"bar", "foo\\bar"}
	fmt.Println("         JSONPointer:", keys.JSONPointer())
	fmt.Println("         DotNotation:", keys.DotNotation(false))
	fmt.Println("DotNotation(bracket):", keys.DotNotation(true))

	k2, err := ParseJSONPointer(keys.JSONPointer())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Parsed:", k2)
	if reflect.DeepEqual(keys, k2) {
		fmt.Println("keys == k2")
	}
}
