// stolen from Eli Bendersky [https://eli.thegreenplace.net]
package test

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

var pretext = `//go:build ignore
package main
import (
    "github.com/alexbezu/gobol/asm"
    "github.com/alexbezu/gobol/asm/dc"
    "github.com/alexbezu/gobol/asm/ds"
)
`

var postfunc1 = `var _ = func() bool {
	test = `
var postfunc2 = `
	return false
}()`

func TestInstructions(t *testing.T) {
	// Find the paths of all input files in the data directory.
	paths, err := filepath.Glob(filepath.Join("instructionsASM", "*.hlasm"))
	if err != nil {
		t.Fatal(err)
	}

	for _, path := range paths {
		_, filename := filepath.Split(path)
		testname := filename[:len(filename)-len(filepath.Ext(path))]

		// Each path turns into a test: the test name is the filename without the extension.
		t.Run(testname, func(t *testing.T) {
			out, err := exec.Command("../gobol", "asm", testname+".go", path).Output()
			if err != nil {
				switch e := err.(type) {
				case *exec.Error:
					t.Fatal("failed gobol executing:", err)
				case *exec.ExitError:
					t.Error("Parser failed rc =", e.ExitCode())
				default:
					panic(err)
				}
			} else {
				src := pretext + string(out) + postfunc1 + testname + postfunc2

				// src = string(format.Source(src))
				// if err != nil {
				// 	t.Fatal("error formatting:", err)
				// }

				f, err := os.OpenFile(testname+".go", os.O_CREATE|os.O_WRONLY, 0644)
				if err != nil {
					t.Fatal("failed opening:"+testname+".go", err)
				}
				if _, err := f.WriteString(src); err != nil {
					t.Fatal("failed writing source code:"+testname+".go", err)
				}
				f.Close()
				out, err = exec.Command("go", "run", "main.go", testname+".go", path).CombinedOutput()
				if err != nil {
					switch err.(type) {
					case *exec.Error:
						t.Fatal("failed test compiling:", err)
					case *exec.ExitError:
						t.Error("Test run failed. " + string(out))
						return
					default:
						panic(err)
					}
				}
				os.Remove(testname + ".go")
			}
		})
	}
}
