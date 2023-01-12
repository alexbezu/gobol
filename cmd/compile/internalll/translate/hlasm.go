package translate

import (
	// "github.com/alexbezu/gobol/asm/macros"
	"github.com/alexbezu/gobol/cmd/compile/internalll/syntax"
	"github.com/alexbezu/gobol/cmd/compile/internalll/translate/macros"
)

type instruction_f func(syntax.Line) string
type Translator_asm struct {
	Src   string
	funcs map[string]instruction_f
}

var funcs = map[string]instruction_f{
	// "A": a,
	// "AH":
	// "AL":
	// "ALR":
	"AP": ap,
}

func (t *Translator_asm) Precompile_tree(tree *syntax.File_asm) {
	for _, line := range tree.Lines {
		f, ok := macros.Macroses[line.Instr]
		if ok {
			f(line, 0)
		}
	}
}

func (t *Translator_asm) Compile_tree(tree *syntax.File_asm) {
	// t.clear_labels(tree)
	for _, line := range tree.Lines {
		// f := funcs[line.Instr]
		// f(line)
		switch line.Instr {
		case "_Comment":
			t.Src += "// " + line.Label
		// case "A":
		case "A":
			t.Src += t.a(line)
		// case "AH":
		// case "AL":
		// case "ALR":
		case "AP":
			t.Src += t.ap(line)
		// case "AR":
		// case "BAL":
		case "BALR":
			t.Src += t.balr(line)
		// case "BAS":
		// case "BASR":
		// case "BASSM":
		// case "BC":
		// case "BCR":
		// case "BCT":
		// case "BCTR":
		// case "BSM":
		// case "BXH":
		// case "BXLE":
		// case "C":
		case "C":
			t.Src += t.c(line)
		// case "CDS":
		// case "CH":
		// case "CL":
		// case "CLC":
		// case "CLCL":
		// case "CLI":
		// case "CLM":
		// case "CLR":
		// case "CP":
		case "CR":
			t.Src += t.cr(line)
		// case "CS":
		// case "CVB":
		// case "CVD":
		// case "D":
		// case "DP":
		// case "DR":
		case "ED":
			t.Src += t.ed(line)
		// case "EDMK":
		// case "EX":
		// case "IC":
		// case "ICM":
		// case "L":
		case "L":
			t.Src += t.l(line)
		case "LA":
			t.Src += t.la(line)
		// case "LCR":
		// case "LH":
		// case "LM":
		// case "LNR":
		// case "LPR":
		// case "LR":
		// case "LTR":
		// case "M":
		// case "MH":
		// case "MP":
		// case "MR":
		case "MVC":
			t.Src += t.mvc(line)
		// case "MVCIN":
		// case "MVCL":
		// case "MVI":
		// case "MVN":
		// case "MVO":
		// case "MVZ":
		// case "N":
		// case "NC":
		// case "NI":
		// case "NR":
		// case "O":
		// case "OC":
		// case "OI":
		// case "OR":
		case "PACK":
			t.Src += t.pack(line)
		// case "S":
		// case "SH":
		// case "SL":
		// case "SLA":
		// case "SLDA":
		// case "SLDL":
		// case "SLL":
		// case "SLR":
		// case "SP":
		// case "SR":
		// case "SRA":
		// case "SRDA":
		// case "SRDL":
		// case "SRL":
		// case "SRP":
		case "ST":
			t.Src += t.st(line)
		// case "STC":
		// case "STCM":
		// case "STH":
		// case "STM":
		// case "SVC":
		// case "TM":
		// case "TR":
		// case "TRT":
		// case "UNPK":
		// case "X":
		// case "XC":
		// case "XI":
		// case "XR":
		// case "ZAP":

		case "B":
			t.Src += t.b(line)
		case "BR":
			t.Src += t.br(line)
		// case "NOP":
		// case "NOPR":
		// case "BH":
		// case "BHR":
		// case "BL":
		// case "BLR":
		// case "BE":
		// case "BER":
		// case "BNH":
		// case "BNHR":
		// case "BNL":
		// case "BNLR":
		case "BNE":
			t.Src += t.bne(line)
		// case "BNER":
		// case "BO":
		// case "BOR":
		// case "BP":
		// case "BPR":
		// case "BM":
		// case "BMR":
		// case "BNP":
		// case "BNPR":
		// case "BNM":
		// case "BNMR":
		// case "BNZ":
		// case "BNZR":
		// case "BZ":
		// case "BZR":
		// case "BNO":
		// case "BNOR":
		default:
			f, ok := macros.Macroses[line.Instr]
			if ok {
				t.Src += f(line, 1)
			} else {
				t.Src += "// todo: " + line.Instr + "\n"
			}
		}
	}
}

func (t *Translator_asm) clear_labels(tree *syntax.File_asm) {
	// var label_set map[string]bool
	for _, line := range tree.Lines {
		switch line.Instr {
		case "B", "BH", "BHR", "BL", "BLR", "BE", "BER", "BNH", "BNHR", "BNL", "BNLR", "BNE", "BNER", "BO", "BOR", "BP", "BPR", "BM", "BMR", "BNP", "BNPR", "BNM", "BNMR", "BNZ", "BNZR", "BZ", "BZR", "BNO", "BNOR":
			// bloom?
			macros.Labels_set[line.Params[0].ParamName] = true
		case "EQU", "DC", "DS", "DCB":
			macros.Labels_set[line.Label] = true
		}
	}
	for i, line := range tree.Lines {
		_, ok := macros.Labels_set[line.Label]
		if !ok {
			tree.Lines[i].Label = ""
		}
	}
}

func (t *Translator_asm) label(lbl string) (ret string) {
	if lbl != "" {
		ret = lbl + ": "
	} else {
		ret = "        "
	}
	return ret
}

func r(param syntax.Param) (r string) {
	switch len(param.Values) {
	case 0:
		r = param.ParamName
	case 1:
		switch param.ParamName {
		case "_Number_asm":
			r = param.Values[0].Value
		default:
			panic("TODO: r1d2x2b2 default")
		}
	case 2:
		panic("TODO: r1d2x2b2 params[0].Values case 2")
	}
	return r
}

// TODO: combine dc() in mocroses
func storage(v syntax.Value) (d string) {
	switch v.Value {
	case "F":
		d = "dc.F(" + v.Extra + ")"
	case "X":
		d = macros.X_generator(v.Extra)
	case "C":
		init := ""
		if v.Extra != "" {
			init = ".INIT(\"" + v.Extra + "\")"
		}
		d = "pl.CHAR(" + ")" + init
	default:
		panic("TODO: r1d2x2b2 case _Storage")
	}
	return d
}

func dxb(param syntax.Param) (D, X, B string) {
	D, X, B = "0", "0", "0"
	switch param.ParamName {
	case "_Number_asm":
		D = param.Values[0].Value
	case "_L_macro":
		D = param.Values[0].Value + ".GetSize()"
	case "_Storage":
		D = storage(param.Values[0])
	default:
		D = param.ParamName
	}

	Xset := false
	for _, v := range param.Values {
		switch v.Tok {
		case syntax.Number_asm, syntax.ID_asm:
			if Xset {
				B = v.Value
			} else {
				X = v.Value
				Xset = true
			}
		case syntax.Plus_asm:
			D += ".P(" + v.Value + ")"
		case syntax.Storage, syntax.L_macro:
		default:
			panic("dlb range param.Values default")
		}
	}
	return D, X, B
}

func db(param syntax.Param) (D, B string) {
	D, B = "0", "0"
	switch param.ParamName {
	case "_Number_asm":
		D = param.Values[0].Value
	case "_L_macro":
		D = param.Values[0].Value + ".GetSize()"
	case "_Storage":
		D = storage(param.Values[0])
	default:
		D = param.ParamName
	}

	for _, v := range param.Values {
		switch v.Tok {
		case syntax.Number_asm, syntax.ID_asm:
			B = v.Value
		case syntax.Plus_asm:
			D += ".P(" + v.Value + ")"
		case syntax.Storage, syntax.L_macro:
		default:
			panic("dlb range param.Values default")
		}
	}

	return D, B
}

func r1d2x2b2(instruction string, line syntax.Line) string {
	var R1, D2, X2, B2 string
	X2 = "0"
	B2 = "0"
	params := line.Params
	if len(params) != 2 {
		panic("parser error, more than two params in r1d2x2b2")
	}

	R1 = r(line.Params[0])
	D2, X2, B2 = dxb(line.Params[1])

	return instruction + R1 + ", " + D2 + ", " + X2 + ", " + B2 + ")\n"
}

func d1lb1d2b2(instruction string, line syntax.Line) (ret string) {
	var D1, L, B1, D2, B2 string
	B1, B2 = "0", "0"
	if len(line.Params) != 2 {
		panic("parser error, more than two params in d1lb1d2b2")
	}
	D1, L, B1 = dxb(line.Params[0])
	D2, B2 = db(line.Params[1])

	ret = instruction + D1 + ", " + B1 + ", " + D2 + ", " + B2
	if L != "0" {
		ret += ", " + L
	}
	return ret + ")\n"
}

func (t *Translator_asm) a(line syntax.Line) string {
	return t.label(line.Label) + r1d2x2b2("asm.A(", line)
}

func ap(line syntax.Line) string {
	var t Translator_asm
	params := line.Params
	return t.label(line.Label) + "asm.AP(" + params[0].ParamName + ", " + params[1].ParamName + ")\n"
}

func (t *Translator_asm) ap(line syntax.Line) string {
	params := line.Params
	return t.label(line.Label) + "asm.AP(" + params[0].ParamName + ", " + params[1].ParamName + ")\n"
}

func (t *Translator_asm) balr(line syntax.Line) string {
	params := line.Params
	return t.label(line.Label) + "asm.BALR(" + params[0].ParamName + ", " + params[1].ParamName + ")\n"
}

func (t *Translator_asm) c(line syntax.Line) string {
	return t.label(line.Label) + r1d2x2b2("asm.C(", line)
}

func (t *Translator_asm) cr(line syntax.Line) string {
	R1 := r(line.Params[0])
	R2 := r(line.Params[1])
	return t.label(line.Label) + "asm.CR(" + R1 + ", " + R2 + ")\n"
}

func (t *Translator_asm) ed(line syntax.Line) string {
	params := line.Params
	return t.label(line.Label) + "asm.ED(" + params[0].ParamName + ", " + params[1].ParamName + ")\n"
}

func (t *Translator_asm) l(line syntax.Line) string {
	return t.label(line.Label) + r1d2x2b2("asm.L(", line)
}

func (t *Translator_asm) la(line syntax.Line) string {
	//41F00000 LA 15,0
	//41F00FFF LA 15,4095
	//41F50FFF LA 15,4095(5)
	//4140D0C0 LA R4,AFIELD     (REG==D is substituted automatically e.g. AFIELD has FFFFF0C0 address and D has FFFFF000 location by USING)
	//4146D0C0 LA R4,AFIELD(R6)
	//41456014 LA R4,20(R5,R6)
	return t.label(line.Label) + r1d2x2b2("asm.LA(", line)
}

//d1lb1d2b2
func (t *Translator_asm) mvc(line syntax.Line) (ret string) {
	// MVC(D1 Objer, D2 Objer, length ...byte)
	// params := line.Params
	// ret = t.label(line.Label) + "asm.MVC(" + params[0].ParamName + ", " + params[1].ParamName + ")\n"
	return t.label(line.Label) + d1lb1d2b2("asm.MVC(", line)
	// return ret
}

func (t *Translator_asm) pack(line syntax.Line) (ret string) {
	//L1L2_BD_DD_BD_DD
	params := line.Params
	ret = t.label(line.Label) + "asm.PACK(" + params[0].ParamName + ", " + params[1].ParamName + ")\n"
	return ret
}

func (t *Translator_asm) st(line syntax.Line) (ret string) {
	//L1L2_BD_DD_BD_DD
	return t.label(line.Label) + r1d2x2b2("asm.ST(", line)
}

// branches

func (t *Translator_asm) b(line syntax.Line) (ret string) {
	// B     READREC | goto     READREC
	params := line.Params
	ret = t.label(line.Label) + "goto " + params[0].ParamName + "\n"
	return ret
}

func (t *Translator_asm) br(line syntax.Line) (ret string) {
	params := line.Params
	if params[0].ParamName == "R14" {
		ret = "return asm.Rui[15]\n"
	} else if params[0].ParamName == "_Number_asm" && params[0].Values[0].Value == "14" {
		ret = "return asm.Rui[15]\n"
	} else {
		ret = t.label(line.Label) + "goto " + params[0].ParamName + "\n"
	}
	return ret
}

func (t *Translator_asm) bne(line syntax.Line) (ret string) {
	params := line.Params
	ret = t.label(line.Label) + "if asm.BNE() {\n\t\tgoto " + params[0].ParamName + "\n\t}\n"
	return ret
}
