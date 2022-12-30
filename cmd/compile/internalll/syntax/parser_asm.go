package syntax

import (
	"io"
	"strconv"
)

// ----------------------------------------------------------------------------
// Nodes
type Node_asm interface {
	Pos() Pos
	aNode()
}

type node_asm struct {
	pos Pos
}

func (n *node_asm) Pos() Pos { return n.pos }
func (*node_asm) aNode()     {}

// ----------------------------------------------------------------------------
// Files
type File_asm struct {
	Lines []Line
	EOF   Pos
	node_asm
}

// ----------------------------------------------------------------------------
// Lines
type (
	Line struct {
		Label string
		Instr string
		// param
		// : ID = (values) PFK=(D0015,'CHAR',1='01',2='02')
		// | ID = value    FEAT=ID or LTH=6 or FILL=C' '
		// | value         ID or 'CHAR' or 15 with keys: 'id' 'char' 'number'
		// | (values)      (TRANS,'PRTRAN3 E    ') and keys is 'list'
		// Params map[string][]Value
		Params []Param
	}

	Param struct {
		ParamName string
		Values    []Value
	}

	// value
	// : NUMBER
	// | NUMBER = Extra
	// | ID
	// | CHAR
	// | X'4020'  CL6'STRING' 5CL72
	Value struct {
		Tok    tokens_asm
		Value  string // NUMBER, ID, CHAR
		Extra  string // STR or nil
		Length int    // length (repetition factor) or nil
	}
)

type parser_asm struct {
	file *PosBase
	errh ErrorHandler
	// mode Mode
	scanner_asm

	base   *PosBase // current position base
	first  error    // first error encountered
	errcnt int      // number of errors encountered

	fnest  int    // function nesting level (for error handling)
	xnest  int    // expression nesting level (for complit ambiguity resolution)
	indent []byte // tracing support
}

func (p *parser_asm) init(file *PosBase, r io.Reader, errh ErrorHandler) {
	p.file = file
	p.errh = errh
	// p.mode = mode
	p.scanner_asm.init(r, nil, directives)

	p.base = file
	p.first = nil
	p.errcnt = 0

	p.fnest = 0
	p.xnest = 0
	p.indent = nil
}

//lines
func (p *parser_asm) fileOrNil() *File_asm {
	f := new(File_asm)
	f.pos = p.pos()

	for p.tok != _EOF_asm {
		for p.tok == _Newline_asm {
			p.next()
			if p.tok == _EOF_asm {
				f.EOF = p.pos()
				return f
			}
		}
		l := p.lines()
		// fmt.Println(l)
		f.Lines = append(f.Lines, l)
	}

	f.EOF = p.pos()
	return f
}

func (p *parser_asm) lines() (l Line) {
	// l.Params = map[string][]Value{}
	l.Params = []Param{}
	switch p.tok {
	case _Label:
		l.Label = p.lit
		p.next()
		fallthrough
	case _Instruction:
		l.Instr = p.lit
	default:
		panic("TODO: lines switch")
	}

	for p.tok != _Newline_asm {
		p.next()
		if p.tok == _EOF_asm {
			break
		}
		if p.tok != _Newline_asm && p.tok != _Comma_asm {
			var param Param
			param.ParamName, param.Values = p.param()
			l.Params = append(l.Params, param)
			// l.Params[id] = values
		}
	}
	return l
}

func (p *parser_asm) param() (key string, vals []Value) {
	vals = make([]Value, 0)
	switch p.tok {
	case _ID_asm:
		key = p.lit
		p.next()
		switch p.tok {
		case _Assign_asm:
			// PFK=(
			// SYSMSG=D0011
			// FILL=C' '
			p.next()
			switch p.tok {
			case _Lparen_asm:
				// PFK=(...
				vals = append(vals, p.values()...)
				p.next()
			case _ID_asm, _Char_asm, _Number_asm:
				// SYSMSG=D0011
				// FILL=C' '
				vals = append(vals, Value{Tok: p.tok, Value: p.lit})
				p.next() //should be , or nl
			case _Storage: // DC X'0204'
				v := Value{Tok: p.tok, Value: p.lit}
				// key = p.lit
				p.next()
				if p.tok == _Char_asm {
					v.Extra = p.lit
					p.next()
				}
				vals = append(vals, v)
			default:
				panic("TODO: params switch 3")
			}
		case _Char_asm:
			// C' P A Y R O L L  R E P O R T'
			vals = append(vals, Value{Tok: p.tok, Value: p.lit})
			key = "char"
			p.next()
		case _Lparen_asm:
			vals = append(vals, p.values()...)
			p.next()
		case _Plus_asm:
			v := Value{Tok: _Plus_asm, Extra: key}
			p.next()
			v.Value = p.lit
			vals = append(vals, v)
			if p.tok == _ID_asm { //    MVC   TARGET+10(20,R1),SOURCE+3
				vals = append(vals, p.values()...)
			}
			// p.next()
		case _Comma_asm, _Newline_asm:
			// NOGEN (empty key)
			return key, vals
		default:
			panic("TODO: params switch 2")
		}
	case _Storage: // DC X'0204'
		v := Value{Tok: p.tok}
		key = p.lit
		p.next()
		if p.tok == _Char_asm {
			v.Extra = p.lit
			p.next()
		}
		vals = append(vals, v)
	case _Storage_len: // DC CL02
		key = p.lit
		v := Value{Tok: p.tok}
		p.next()
		v.Value = p.lit // length
		p.next()
		if p.tok == _Char_asm {
			v.Extra = p.lit
			p.next()
		}
		vals = append(vals, v)
	case _Storage_rf: // DC 5CL12
		v := Value{Tok: p.tok, Value: p.lit}
		p.next()
		key = p.lit
		p.next()
		v.Extra = p.lit
		vals = append(vals, v)
	case _Char_asm:
		// 'CHAR'
		vals = append(vals, Value{Tok: p.tok, Value: p.lit})
		key = "char"
		p.next()
	case _Lparen_asm:
		// (...
		key = "list"
		vals = append(vals, p.values()...)
		p.next()
	case _Number_asm:
		key = p.tok.String()
		vals = append(vals, Value{Tok: p.tok, Value: p.lit})
	case _Star_asm:
		key = "star"
		vals = append(vals, Value{Tok: p.tok, Value: p.lit})
	case _Assign_asm:
		v := Value{Tok: p.tok}
		key = "assign"
		p.next()
		switch p.tok {
		case _Storage:
			v.Tok = p.tok
			key = "_Storage"
			v.Value = p.lit
			p.next()
			if p.tok == _Char_asm {
				v.Extra = p.lit
			} else {
				panic("only Char in _Storage is allowed")
			}
		case _Storage_len:
			v.Tok = p.tok
			p.next()
			l, err := strconv.Atoi(p.lit)
			if err == nil {
				v.Length = l
			}
		default:
			panic("after _Assign_asm param")
		}
		vals = append(vals, v)
	case _L_macro:
		key = p.tok.String()
		v := Value{Tok: p.tok}
		p.next()
		if p.tok == _ID_asm {
			v.Value = p.lit
			vals = append(vals, v)
		} else {
			panic("only ID is allowed in _L_macro")
		}
		p.next()
		if p.tok == _Lparen_asm {
			vals = append(vals, p.values()...)
		} else if p.tok == _Comma_asm || p.tok == _Newline_asm {
			return key, vals
		} else {
			panic("only () is allowed in _L_macro after ID")
		}
	case _Comma_asm, _Newline_asm:
		return key, vals
	default:
		panic("TODO: params switch 1")
	}

	return key, vals
}

// (ID,1='01')
// ((3,28,CURSOR))
// (,OVB)
func (p *parser_asm) values() (ret []Value) {
	dblParen := false
	p.next()
	for {
		switch p.tok {
		case _Lparen_asm: // ((...
			dblParen = true
			p.next()
		case _Rparen_asm: // ...))
			if dblParen {
				p.next()
			}
			return ret
		case _Number_asm: // ((3,28,CURSOR))
			val := Value{Tok: p.tok, Value: p.lit}
			p.next()
			if p.tok == _Assign_asm { // (ID,1='01')
				p.next()
				val.Extra = p.lit
				p.next()
			}
			ret = append(ret, val)
		case _ID_asm, _Char_asm: // ('CHAR',OVB)
			ret = append(ret, Value{Tok: p.tok, Value: p.lit})
			p.next()
		case _L_macro:
			val := Value{Tok: p.tok}
			p.next()
			val.Value = p.lit
			ret = append(ret, val)
			p.next()
		case _Comma_asm:
			p.next()
		case _Newline_asm:
			p.next()
			return ret
		default:
			panic("TODO: values")
		}
	}
}

func (p *parser_asm) pos() Pos { return p.posAt(p.line, p.col) }

// posAt returns the Pos value for (line, col) and the current position base.
func (p *parser_asm) posAt(line, col uint) Pos {
	return MakePos(p.base, line, col)
}
