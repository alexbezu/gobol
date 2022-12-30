// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file implements scanner, a lexical tokenizer for
// Go source. After initialization, consecutive calls of
// next advance the scanner one token at a time.
//
// This file, source.go, tokens.go, and token_string.go are self-contained
// (`go tool compile scanner.go source.go tokens.go token_string.go` compiles)
// and thus could be made into their own package.

package syntax

import (
	"fmt"
	"io"
	"strconv"
	"unicode"
	"unicode/utf8"
)

var debug_tokens = false

type scanner_asm struct {
	source
	// mode uint
	nl bool

	// current token, valid after calling next()
	line, col uint
	tok       tokens_asm
	inparams  bool
	// paramspace bool
	lit string
}

func (s *scanner_asm) init(src io.Reader, errh func(line, col uint, msg string), mode uint) {
	s.source.init(src, errh)
	s.tok = _Newline_asm
	// s.mode = mode
	// s.nlsemi = false
}

// errorf reports an error at the most recently read character position.
func (s *scanner_asm) errorf(format string, args ...interface{}) {
	s.error(fmt.Sprintf(format, args...))
}

// errorAtf reports an error at a byte column offset relative to the current token start.
func (s *scanner_asm) errorAtf(offset int, format string, args ...interface{}) {
	s.errh(s.line, s.col+uint(offset), fmt.Sprintf(format, args...))
}

//1        10    16                                                      72
//LABEL1   INSTR TYPE=(3270,2),FEAT=IGNORE,SYSMSG=D0011,PFK=(D0015,1='01'X
//               ,2='02',3='03',8='08',10='10',13='01',   COMMENT        X
//               14='02',15='03',20='08',22='10')         COMMENT
func (s *scanner_asm) next() {
	// nl := s.nl
	// s.nl = false
	s.lit = ""

	// skip white space
	s.stop()
	for s.ch == ' ' || s.ch == '\t' || s.ch == '\n' /* && !nl*/ || s.ch == '\r' {
		if s.inparams {
			for s.ch != '\n' {
				if s.ch == -1 {
					s.tok = _EOF_asm
					return
				}
				if s.col == 72 && s.ch == 'X' {
					s.nextch()
					for s.ch == ' ' || s.ch == '\n' || s.ch == '\r' {
						s.nextch()
					}
					goto symbols
				}
				s.nextch()
				s.line, s.col = s.pos()
			}
			s.tok = _Newline_asm
			s.lit = "nl"
			s.inparams = false
			return
		}
		if s.ch == '\n' {
			s.tok = _Newline_asm
			s.lit = "nl"
			s.inparams = false
			s.nextch()
			return
		}
		s.nextch()
	}

	s.line, s.col = s.pos()
	//skip X at 72
	if s.col == 72 && s.ch == 'X' {
		if s.inparams {
			// s.paramspace = false
			s.nextch()
			for s.ch == ' ' || s.ch == '\t' || s.ch == '\n' {
				s.nextch()
			}
		} else {
			s.error("X shouldn't be there")
		}
	}

	//skip * comment
	if s.col == 1 && s.ch == '*' {
		s.nextch()
		for s.ch != '\n' /*|| startCol == s.col*/ {
			s.nextch()
			s.lit += string(s.ch)
			s.line, s.col = s.pos()
		}
		s.tok = _Comment
		return
	}

	if s.ch == '\n' {
		s.tok = _Newline_asm
		s.lit = "nl"
		s.inparams = false
		// s.paramspace = false
		return
	}

symbols:
	s.start()
	s.line, s.col = s.pos()
	if isLetter(s.ch) /* || s.ch >= utf8.RuneSelf && s.atIdentChar(true) */ {
		switch s.ch {
		case 'A', 'B', 'C', 'F', 'H', 'P', 'X': // CL72
			if s.col > 10 { // skip label and instruction
				s.tok = _Storage
				s.lit = string(s.ch)
				s.nextch()
				switch s.ch {
				case 'L':
					s.nextch()
					if unicode.IsNumber(s.ch) {
						s.tok = _Storage_len
					} else {
						s.ident()
						s.tok = _ID_asm
						s.inparams = true
					}
				case '\'':
				default:
					s.ident()
					s.tok = _ID_asm
					s.inparams = true
				}
				if debug_tokens {
					fmt.Println(s.tok, s.lit)
				}
				return
			}
		case 'L':
			if s.col > 10 { // skip label and instruction
				s.nextch()
				if s.ch == '\'' {
					s.tok = _L_macro
					s.nextch()
				} else {
					s.ident()
					s.tok = _ID_asm
					s.inparams = true
				}
				if debug_tokens {
					fmt.Println(s.tok, s.lit)
				}
				return
			}
		}

		// xORc := s.ch
		s.nextch()
		// if (xORc == 'X' || xORc == 'C') && s.ch == '\'' {
		// 	s.hexchar(xORc)
		// 	return
		// }
		s.ident()
		if s.col == 1 {
			s.tok = _Label
			s.inparams = false
		} else if s.col > 1 && (s.tok == _Label || s.tok == _Newline_asm) {
			s.tok = _Instruction
			s.inparams = false
		} else {
			s.tok = _ID_asm
			s.inparams = true
		}
		if debug_tokens {
			fmt.Println(s.tok)
		}
		return
	}

	s.inparams = true
	switch s.ch {
	case -1:
		s.lit = "EOF"
		s.tok = _EOF_asm

	case '\n', '\r':
		s.nextch()
		s.lit = "nl"
		s.tok = _Newline_asm

	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		s.number(false)

	case '\'':
		s.stdString()
		s.tok = _Char_asm

	case '=':
		s.nextch()
		s.tok = _Assign_asm

	case '(':
		s.nextch()
		s.tok = _Lparen_asm

	case ')':
		s.nextch()
		s.tok = _Rparen_asm

	case ',':
		s.nextch()
		s.tok = _Comma_asm

	case '-':
		s.nextch()
		s.lit = "-"
		s.tok = _Char_asm

	case '*':
		s.nextch()
		s.lit = "*"
		s.tok = _Star_asm

	case '+':
		s.nextch()
		s.lit = "+"
		s.tok = _Plus_asm

	default:
		fmt.Printf("%s", string(s.ch))
		panic("TODO:")
	}
	if debug_tokens {
		fmt.Println(s.tok)
	}
}

// C\'.+\'                       return 'CHAR' e.g. C' '
// X\'[\dA-F]+\'                 return 'HEX'  e.g. X'0090'
func (s *scanner_asm) hexchar(ch rune) {
	s.nextch()
	for {
		if s.ch == '\'' {
			s.nextch()
			break
		}
		s.nextch()
	}
	if ch == 'X' {
		num := string(s.segment())
		value, err := strconv.ParseInt(num[2:len(num)-1], 16, 64)
		if err != nil {
			s.errorf("strconv.ParseInt %#U", s.ch)
			return
		}
		s.tok = _Number_asm
		s.lit = strconv.Itoa(int(value))
		return
	} //else ch == 'C'
	s.tok = _Char_asm
	s.lit = string(s.segment())
}

func (s *scanner_asm) ident() {
	// accelerate common case (7bit ASCII)
	for isLetter(s.ch) || isDecimal(s.ch) || s.ch == '-' {
		s.nextch()
	}

	// general case
	if s.ch >= utf8.RuneSelf {
		for s.atIdentChar(false) {
			s.nextch()
		}
	}

	lit := s.segment()
	s.lit = string(lit)

}

func (s *scanner_asm) atIdentChar(first bool) bool {
	switch {
	case unicode.IsLetter(s.ch):
		// ok
	case unicode.IsDigit(s.ch):
		if first {
			s.errorf("identifier cannot begin with digit %#U", s.ch)
		}
	case s.ch >= utf8.RuneSelf:
		s.errorf("invalid character %#U in identifier", s.ch)
	default:
		return false
	}
	return true
}

func (s *scanner_asm) digits(base int, invalid *int) (digsep int) {
	if base <= 10 {
		max := rune('0' + base)
		for isDecimal(s.ch) || s.ch == '_' {
			ds := 1
			if s.ch == '_' {
				ds = 2
			} else if s.ch >= max && *invalid < 0 {
				_, col := s.pos()
				*invalid = int(col - s.col) // record invalid rune index
			}
			digsep |= ds
			s.nextch()
		}
	} else {
		for isHex(s.ch) || s.ch == '_' {
			ds := 1
			if s.ch == '_' {
				ds = 2
			}
			digsep |= ds
			s.nextch()
		}
	}
	return
}

func (s *scanner_asm) number(seenPoint bool) {
	base := 10        // number base
	prefix := rune(0) // one of 0 (decimal), '0' (0-octal), 'x', 'o', or 'b'
	digsep := 0       // bit 0: digit present, bit 1: '_' present
	invalid := -1     // index of invalid digit in literal, or < 0

	// integer part
	if !seenPoint {
		/* if s.ch == '0' {
			s.nextch()
			switch lower(s.ch) {
			case 'x':
				s.nextch()
				base, prefix = 16, 'x'
			case 'o':
				s.nextch()
				base, prefix = 8, 'o'
			case 'b':
				s.nextch()
				base, prefix = 2, 'b'
			default:
				base, prefix = 8, '0'
				digsep = 1 // leading 0
			}
		} */
		digsep |= s.digits(base, &invalid)
		switch s.ch {
		case '.':
			if prefix == 'o' || prefix == 'b' {
				s.errorf("invalid radix point in %s literal", baseName(base))
				// ok = false
			}
			s.nextch()
			seenPoint = true
		case 'A', 'B', 'C', 'F', 'H', 'P', 'X': // 5CL72 or 0(3,4)
			s.tok = _Storage_rf
			s.lit = string(s.segment())
			return
		case '(':
			s.tok = _ID_asm
			s.lit = string(s.segment())
			return
		}
	}

	s.tok = _Number_asm
	s.lit = string(s.segment())
}

func (s *scanner_asm) stdString() {
	s.nextch()
	for {
		s.line, s.col = s.pos()
		if s.ch == 'X' && s.col == 72 {
			s.nextch()
		} else if s.col < 15 {
			s.nextch()
		} else if s.ch == '\'' {
			s.nextch()
			break
		} else if s.ch == '\n' {
			s.nextch()
		} else if s.ch < 0 {
			s.errorAtf(0, "string not terminated")
			break
		} else {
			s.lit += string(s.ch)
			s.nextch()
		}
	}
}

// func (s *scanner_asm) storage() {
// 	s.lit = string(s.ch)
// 	s.nextch()

// }

// func Temporardebug() {
// 	// transactions/DSN8IPD.hlasm
// 	f, err := os.Open("test/format.hlasm")
// 	if err != nil {
// 		return
// 	}
// 	defer f.Close()

// 	// var s scanner_asm
// 	// s.init(f, nil, 0)
// 	// for s.tok != _EOF_asm {
// 	// 	s.next()
// 	// 	if s.tok == _Newline_asm {
// 	// 		fmt.Println()
// 	// 	} else {
// 	// 		fmt.Printf("{%s %s}", s.lit, s.tok)
// 	// 	}
// 	// }

// 	var p parser_asm
// 	fbase := NewFileBase("f")
// 	p.init(fbase, f, nil)
// 	p.next()

// 	fmt.Println(p.fileOrNil())
// }
