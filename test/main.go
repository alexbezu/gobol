//go:build ignore

package main

import (
	"os"
)

var test func() uintptr

func main() {
	if test != nil {
		ret := int(test())
		os.Exit(ret)
	}
	os.Exit(1)
}
