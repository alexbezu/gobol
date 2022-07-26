package macros

import (
	"gobol/src/cmd/compile/internalll/syntax"
)

func print(line syntax.Line, traverse int) string {

	return ""
}

var _ = func() bool {
	Macroses["PRINT"] = print
	return false
}()
