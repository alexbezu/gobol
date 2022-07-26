package macros

import (
	"gobol/src/cmd/compile/internalll/syntax"
)

func open(line syntax.Line, traverse int) string {
	//INOPEN   OPEN  (INDCB,INPUT)
	//asm.OPEN(INDCB, "INPUT")
	values := line.Params[0].Values
	return label(line.Label) + "asm.OPEN(" + values[0].Value + ", \"" + values[1].Value + "\")\n"
}

func get(line syntax.Line, traverse int) (ret string) {
	// READREC  GET   INDCB,PAYREC
	// READREC: EODAD := asm.GET(INDCB, PAYREC)
	// if EODAD {
	// 	goto INCLOS
	// }
	params := line.Params
	ret = label(line.Label) + "EODAD := asm.GET(" + params[0].ParamName + ", " + params[1].ParamName + ")\n"
	gotolabel, ok := _EODADs[params[0].ParamName]
	if ok {
		ret += "\tif EODAD {\n\t\tgoto " + gotolabel + "\n\t}\n"
	}
	return ret
}

func put(line syntax.Line, traverse int) string {
	//PUT   OUTDCB,PTITLE
	params := line.Params
	return label(line.Label) + "asm.PUT(" + params[0].ParamName + ", " + params[1].ParamName + ")\n"
}

func close(line syntax.Line, traverse int) string {
	//LABEL CLOSE (INDCB) | asm.CLOSE(INDCB)
	values := line.Params[0].Values
	return label(line.Label) + "asm.CLOSE(" + values[0].Value + ")\n"
}

func dcb(line syntax.Line, traverse int) string {
	// INDCB    DCB   MACRF=GM,DDNAME=INDD,DSORG=PS,EODAD=INCLOS | var INDCB = &asm.DCB{MACRF: "GM", DDNAME: "INDD", DSORG: "PS"}
	// OUTDCB   DCB   MACRF=PM,DDNAME=OUTDD,DSORG=PS             | var OUTDCB = &asm.DCB{MACRF: "PM", DDNAME: "OUTDD", DSORG: "PS"}
	var strParams string
	if traverse == 0 {
		var EODAD string
		for _, param := range line.Params {
			if param.ParamName == "EODAD" {
				EODAD = param.Values[0].Value
				Labels_set[EODAD] = true
			}
		}
		_EODADs[line.Label] = EODAD
		return ""
	} else if traverse == 1 {
		for _, param := range line.Params {
			strParams += param.ParamName + ": \"" + param.Values[0].Value + "\", "
		}
		strParams = strParams[0 : len(strParams)-2]
	}
	return "var " + line.Label + " = &asm.DCB{" + strParams + "}\n"
}

var _EODADs map[string]string

var _ = func() bool {
	_EODADs = map[string]string{}
	Macroses["OPEN"] = open
	Macroses["GET"] = get
	Macroses["PUT"] = put
	Macroses["CLOSE"] = close
	Macroses["DCB"] = dcb
	return false
}()
