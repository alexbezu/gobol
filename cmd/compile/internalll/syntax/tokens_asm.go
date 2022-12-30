// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package syntax

type tokens_asm uint

//go:generate stringer -type tokens_asm tokens_asm.go

const (
	_ tokens_asm = iota

	_EOF_asm     // EOF
	_Newline_asm // NL

	_ID_asm      // id
	_Literal_asm // literal
	_Instruction // instruction
	_Label       // label
	_Char_asm    // char
	_Number_asm  // number
	_Storage     // storage
	_Storage_len // storage_length
	_Storage_rf  // storage_repetition_factor
	_L_macro     // length_macros

	// delimiters
	_Lparen_asm // (
	_Rparen_asm // )
	_Comma_asm  // ,
	_Semi_asm   // ;
	_Assign_asm // =
	_Star_asm   // *
	_Plus_asm   // +

	_Equ     // EQU
	_Comment // comment
	// empty line comment to exclude it from .String
	tokenCount_asm //
)

const (
	ID_asm      = _ID_asm
	L_macro     = _L_macro
	Number_asm  = _Number_asm
	Plus_asm    = _Plus_asm
	Storage     = _Storage
	Storage_len = _Storage_len
	Storage_rf  = _Storage_rf
)

// contains reports whether tok is in tokset.
// func contains_asm(tokset uint64, tok tokens_asm) bool {
// 	return tokset&(1<<tok) != 0
// }

// type Operator_asm uint

// : generate stringer -type Operator_asm -linecomment tokens_pli.go

// const (
// 	_ Operator_asm = iota
// )
