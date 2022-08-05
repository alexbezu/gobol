package macros

import (
	"github.com/alexbezu/gobol/cmd/compile/internalll/syntax"
)

func csect(line syntax.Line, traverse int) string {
	return ""
}

var _ = func() bool {
	Macroses["CSECT"] = csect
	return false
}()
