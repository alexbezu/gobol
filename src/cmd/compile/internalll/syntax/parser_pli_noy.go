// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package syntax

import (
	"fmt"
	"io"
	"strconv"
	"strings"
)

const debug = false
const trace = false

type parser struct {
	file *PosBase
	errh ErrorHandler
	// mode Mode
	scanner_pli

	base   *PosBase // current position base
	first  error    // first error encountered
	errcnt int      // number of errors encountered

	fnest  int    // function nesting level (for error handling)
	xnest  int    // expression nesting level (for complit ambiguity resolution)
	indent []byte // tracing support
}

func (p *parser) init(file *PosBase, r io.Reader, errh ErrorHandler) {
	p.file = file
	p.errh = errh
	// p.mode = mode
	p.scanner_pli.init(
		r,
		// Error and directive handler for scanner.
		// Because the (line, col) positions passed to the
		// handler is always at or after the current reading
		// position, it is safe to use the most recent position
		// base to compute the corresponding Pos value.
		func(line, col uint, msg string) {
			if msg[0] != '/' {
				p.errorAt(p.posAt(line, col), msg)
				return
			}

			// otherwise it must be a comment containing a line or go: directive.
			// //line directives must be at the start of the line (column colbase).
			// /*line*/ directives can be anywhere in the line.
			text := commentText(msg)
			if (col == colbase || msg[1] == '*') && strings.HasPrefix(text, "line ") {
				var pos Pos // position immediately following the comment
				if msg[1] == '/' {
					// line comment (newline is part of the comment)
					pos = MakePos(p.file, line+1, colbase)
				} else {
					// regular comment
					// (if the comment spans multiple lines it's not
					// a valid line directive and will be discarded
					// by updateBase)
					pos = MakePos(p.file, line, col+uint(len(msg)))
				}
				p.updateBase(pos, line, col+2+5, text[5:]) // +2 to skip over // or /*
				return
			}
		},
		directives,
	)

	p.base = file
	p.first = nil
	p.errcnt = 0

	p.fnest = 0
	p.xnest = 0
	p.indent = nil
}

func (p *parser) allowGenerics() bool { return p.mode&AllowGenerics != 0 }

// updateBase sets the current position base to a new line base at pos.
// The base's filename, line, and column values are extracted from text
// which is positioned at (tline, tcol) (only needed for error messages).
func (p *parser) updateBase(pos Pos, tline, tcol uint, text string) {
	i, n, ok := trailingDigits(text)
	if i == 0 {
		return // ignore (not a line directive)
	}
	// i > 0

	if !ok {
		// text has a suffix :xxx but xxx is not a number
		p.errorAt(p.posAt(tline, tcol+i), "invalid line number: "+text[i:])
		return
	}

	var line, col uint
	i2, n2, ok2 := trailingDigits(text[:i-1])
	if ok2 {
		//line filename:line:col
		i, i2 = i2, i
		line, col = n2, n
		if col == 0 || col > PosMax {
			p.errorAt(p.posAt(tline, tcol+i2), "invalid column number: "+text[i2:])
			return
		}
		text = text[:i2-1] // lop off ":col"
	} else {
		//line filename:line
		line = n
	}

	if line == 0 || line > PosMax {
		p.errorAt(p.posAt(tline, tcol+i), "invalid line number: "+text[i:])
		return
	}

	// If we have a column (//line filename:line:col form),
	// an empty filename means to use the previous filename.
	filename := text[:i-1] // lop off ":line"
	trimmed := false
	if filename == "" && ok2 {
		filename = p.base.Filename()
		trimmed = p.base.Trimmed()
	}

	p.base = NewLineBase(pos, filename, trimmed, line, col)
}

func commentText(s string) string {
	if s[:2] == "/*" {
		return s[2 : len(s)-2] // lop off /* and */
	}

	// line comment (does not include newline)
	// (on Windows, the line comment may end in \r\n)
	i := len(s)
	if s[i-1] == '\r' {
		i--
	}
	return s[2:i] // lop off //, and \r at end, if any
}

func trailingDigits(text string) (uint, uint, bool) {
	// Want to use LastIndexByte below but it's not defined in Go1.4 and bootstrap fails.
	i := strings.LastIndex(text, ":") // look from right (Windows filenames may contain ':')
	if i < 0 {
		return 0, 0, false // no ":"
	}
	// i >= 0
	n, err := strconv.ParseUint(text[i+1:], 10, 0)
	return uint(i + 1), uint(n), err == nil
}

func (p *parser) got(tok tokens_pli) bool {
	if p.tok == tok {
		p.next()
		return true
	}
	return false
}

func (p *parser) want(tok tokens_pli) {
	if !p.got(tok) {
		p.syntaxError("expecting " + tokstring(tok))
		p.advance()
	}
}

// gotAssign is like got(_Assign) but it also accepts ":="
// (and reports an error) for better parser error recovery.
func (p *parser) gotAssign() bool {
	switch p.tok {
	// case _Define:
	// 	p.syntaxError("expecting =")
	// 	fallthrough
	case _Assign:
		p.next()
		return true
	}
	return false
}

// ----------------------------------------------------------------------------
// Error handling

// posAt returns the Pos value for (line, col) and the current position base.
func (p *parser) posAt(line, col uint) Pos {
	return MakePos(p.base, line, col)
}

// error reports an error at the given position.
func (p *parser) errorAt(pos Pos, msg string) {
	err := Error{pos, msg}
	if p.first == nil {
		p.first = err
	}
	p.errcnt++
	if p.errh == nil {
		panic(p.first)
	}
	p.errh(err)
}

// syntaxErrorAt reports a syntax error at the given position.
func (p *parser) syntaxErrorAt(pos Pos, msg string) {
	if trace {
		p.print("syntax error: " + msg)
	}

	if p.tok == _EOF && p.first != nil {
		return // avoid meaningless follow-up errors
	}

	// add punctuation etc. as needed to msg
	switch {
	case msg == "":
		// nothing to do
	case strings.HasPrefix(msg, "in "), strings.HasPrefix(msg, "at "), strings.HasPrefix(msg, "after "):
		msg = " " + msg
	case strings.HasPrefix(msg, "expecting "):
		msg = ", " + msg
	default:
		// plain error - we don't care about current token
		p.errorAt(pos, "syntax error: "+msg)
		return
	}

	// determine token string
	var tok string
	switch p.tok {
	case _Name, _Semi:
		tok = p.lit
	case _Literal:
		tok = "literal " + p.lit
	case _Operator:
		tok = p.op.String()
	case _AssignOp:
		tok = p.op.String() + "="
	case _IncOp:
		tok = p.op.String()
		tok += tok
	default:
		tok = tokstring(p.tok)
	}

	p.errorAt(pos, "syntax error: unexpected "+tok+msg)
}

// tokstring returns the English word for selected punctuation tokens
// for more readable error messages. Use tokstring (not tok.String())
// for user-facing (error) messages; use tok.String() for debugging
// output.
func tokstring(tok tokens_pli) string {
	switch tok {
	case _Comma:
		return "comma"
	case _Semi:
		return "semicolon or newline"
	}
	return tok.String()
}

// Convenience methods using the current token position.
func (p *parser) pos() Pos               { return p.posAt(p.line, p.col) }
func (p *parser) error(msg string)       { p.errorAt(p.pos(), msg) }
func (p *parser) syntaxError(msg string) { p.syntaxErrorAt(p.pos(), msg) }

// Advance consumes tokens until it finds a token of the stopset or followlist.
// The stopset is only considered if we are inside a function (p.fnest > 0).
// The followlist is the list of valid tokens that can follow a production;
// if it is empty, exactly one (non-EOF) token is consumed to ensure progress.
func (p *parser) advance(followlist ...tokens_pli) {
	if trace {
		p.print(fmt.Sprintf("advance %s", followlist))
	}

	// compute follow set
	// (not speed critical, advance is only called in error situations)
	var followset uint64 = 1 << _EOF // don't skip over EOF
	if len(followlist) > 0 {
		for _, tok := range followlist {
			followset |= 1 << tok
		}
	}

	for !contains(followset, p.tok) {
		if trace {
			p.print("skip " + p.tok.String())
		}
		p.next()
		if len(followlist) == 0 {
			break
		}
	}

	if trace {
		p.print("next " + p.tok.String())
	}
}

// usage: defer p.trace(msg)()
func (p *parser) trace(msg string) func() {
	p.print(msg + " (")
	const tab = ". "
	p.indent = append(p.indent, tab...)
	return func() {
		p.indent = p.indent[:len(p.indent)-len(tab)]
		if x := recover(); x != nil {
			panic(x) // skip print_trace
		}
		p.print(")")
	}
}

func (p *parser) print(msg string) {
	fmt.Printf("%5d: %s%s\n", p.line, p.indent, msg)
}

// pl1stmtlist
// 	: pl1stmt ';'
// 	| pl1stmtlist pl1stmt ';' EOF ;
func (p *parser) pl1stmtList() *File {
	if trace {
		defer p.trace("pl1stmtList")()
	}

	f := new(File)
	f.pos = p.pos()

	for p.tok != _EOF {
		// switch p.tok {
		// case _Const:
		// default:
		// }
		// f.TopProc = p.name()
		s := p.pl1stmt()
		// f.DeclList = p.appendGroup(f.DeclList, p.importDecl)
		if s == nil {
			break
		}
		// l = append(l, s)
	}

	f.EOF = p.pos()
	return f
}

// pl1stmt
// 	: stmt
// 	| varname ':' stmt
// 	| stmtscope
// 	| varname ':' stmtscope
// 	| varname ':' procedurestmt
// 	| endstmt
// 	| varname ':' endstmt
//  ;
func (p *parser) pl1stmt() (l []Stmt) {
	if trace {
		defer p.trace("pl1stmtList")()
	}

	for p.tok != _EOF {
		// switch p.tok {
		// case _Const:
		// default:
		// }
		s := p.pl1stmt()
		if s == nil {
			break
		}
		// l = append(l, s)
	}
	return
}

// stmt
// 	: /* empty */
// 	| ALLOCATE allocateoptionlist
// 	| assignstmt
// 	| attachstmt
// 	| CALL calloptionlist
// 	| CLOSE closegrouplist
// 	| dclstmt
// 	| defaultstmt
// 	| DEFINE ALIAS varname dcloptionlist
// 	| defineordinalstmt
// 	| definestructurestmt
// 	| DELAY '(' expr ')'
// 	| DELETE deleteoptionlist
// 	| DETACH THREAD '(' varnameref ')'
// 	| displaystmt
// 	| ELSE block
// 	| EXEC varname INCLUDE varname
// 	| EXEC varname WHENEVER varname GO TO varname
// 	| EXIT
// 	| FETCH fetchoptioncommalist
// 	| flushstmt
// 	| FORMAT_STMT '(' editformatlist ')'
// 	| FREE freeoption
// 	| getstmt
// 	| gotostmt
// 	| IF expr THEN block
// 	| iteratestmt
// 	| leavestmt
// 	| LOCATE varnameref locateoptionlist
// 	| onstmt
// 	| OPEN opengrouplist
// 	| OTHERWISE block
// 	| putstmt
// 	| READ readoptionlist
// 	| releasestmt
// 	| returnstmt
// 	| revertstmt
// 	| rewritestmt
// 	| SELECT '(' expr ')' pl1stmtlist END
// 	| signalstmt
// 	| stopstmt
// 	| waitstmt
// 	| WHEN '(' varnamedimensioncommalist ')' block
// 	| writestmt
// 	| unlockstmt
// 	;
