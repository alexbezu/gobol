package macros

import "github.com/alexbezu/gobol/cmd/compile/internalll/syntax"

func USING(line syntax.Line, traverse int) string {
	return ""
}

var _ = func() bool {
	Macroses["USING"] = print
	return false
}()
