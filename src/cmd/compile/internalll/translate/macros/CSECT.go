package macros

import (
	"gobol/src/cmd/compile/internalll/syntax"
)

func csect(line syntax.Line, traverse int) string {
	return ""
}

var _ = func() bool {
	Macroses["CSECT"] = csect
	return false
}()
