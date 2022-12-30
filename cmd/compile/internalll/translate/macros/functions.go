package macros

import "github.com/alexbezu/gobol/cmd/compile/internalll/syntax"

func subentry(line syntax.Line, traverse int) string {
	return "func " + line.Label + "() uintptr {\n"
}

func subexit(line syntax.Line, traverse int) string {
	return "}\n"
}

var _ = func() bool {
	Macroses["SUBENTRY"] = subentry
	Macroses["SUBEXIT"] = subexit
	return false
}()
