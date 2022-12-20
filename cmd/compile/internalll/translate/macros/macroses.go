package macros

import (
	"crypto/rand"
	"encoding/base64"

	"github.com/alexbezu/gobol/cmd/compile/internalll/syntax"
)

var Macroses = map[string]func(syntax.Line, int) string{}
var Labels_set = map[string]bool{}

func label(lbl string) (ret string) {
	if lbl != "" {
		ret = lbl + ": "
	} else {
		ret = "        "
	}
	return ret
}

func wto(line syntax.Line, traverse int) string {
	//WTO   PTITLE 2 asm.WTO(PTITLE)
	params := line.Params
	return label(line.Label) + "asm.WTO(" + params[0].ParamName + ")\n"
}

func equ(line syntax.Line, traverse int) string {
	//R0       EQU   0
	//const R0 = 0
	params := line.Params
	if params[0].Values[0].Value == "*" {
		return line.Label + ":\n"
	}
	return "const " + line.Label + " = " + params[0].Values[0].Value + "\n"
}

func dc(line syntax.Line, traverse int) (ret string) {
	//EMPTY    DC    CL1' '    | var EMPTY = pl.CHAR(1).INIT(" ")
	//EDWD     DC    X'402020' | var EDWD = ds.X(0x40, 0x20, 0x20)
	//BFIELD   DC    F'-1'    BFIELD = X'FFFFFFFF'
	param := line.Params[0]
	if line.Label == "" {
		buff := make([]byte, 6)
		rand.Read(buff)
		line.Label = "_" + base64.RawURLEncoding.EncodeToString(buff)[:7]
	}
	ret = "var " + line.Label + " = "
	length := "1"
	init := ""

	switch param.Values[0].Tok {
	case syntax.Storage_len:
		length = param.Values[0].Value
	case syntax.Storage:
		length = "1"
	}

	if param.Values[0].Extra != "" {
		init = ".INIT(\"" + param.Values[0].Extra + "\")"
	}

	switch param.ParamName {
	case "A":
	case "B":
	case "C":
		ret += "pl.CHAR(" + length + ")" + init
	case "F":
		ret += "dc.F(" + param.Values[0].Extra + ")"
	case "H":
	case "P":
	case "X":
		init := ""
		for i := 0; i < len(param.Values[0].Extra); i += 2 {
			digit := param.Values[0].Extra[i : i+2]
			init += "0x" + digit + ", "
		}
		ret += "ds.X(" + init[0:len(init)-2] + ")"
	default:
		ret = "syntax error dc " + "\n"
	}

	return ret + "\n"
}

func ds(line syntax.Line, traverse int) (ret string) {
	//PAYREC   DS    0CL80  | var PAYREC = pl.CL0(80)
	//EMPID    DS    CL4    | var EMPID = pl.CHAR(4)
	param := line.Params[0]
	if line.Label == "" {
		buff := make([]byte, 6)
		rand.Read(buff)
		line.Label = "_" + base64.RawURLEncoding.EncodeToString(buff)[:7]
	}
	ret = "var " + line.Label + " = "
	switch param.ParamName {
	case "A":
	case "B":
	case "C":
		switch param.Values[0].Tok {
		case syntax.Storage_len:
			ret += "pl.CHAR(" + param.Values[0].Value + ")"
		case syntax.Storage_rf:
			ret += "pl.CL0(" + param.Values[0].Extra + ")"
		}
	case "F":
	case "H":
	case "P":
	case "X":
	default:
		ret = "syntax error ds " + "\n"
	}

	return ret + "\n"
}

var _ = func() bool {
	Macroses["WTO"] = wto
	Macroses["EQU"] = equ
	Macroses["DC"] = dc
	Macroses["DS"] = ds
	return false
}()
