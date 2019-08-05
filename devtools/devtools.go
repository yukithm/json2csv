// +build devtools

package devtools

import (
	_ "github.com/mitchellh/gox"
)

//go:generate go build -v -o=./bin/gox github.com/mitchellh/gox
