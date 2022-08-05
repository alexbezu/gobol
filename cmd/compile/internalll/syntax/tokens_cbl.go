// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package syntax

type tokens_cbl uint

//go:generate stringer -type Operator_cbl tokens_pli.go

const (
	_        tokens_cbl = iota
	_EOF_cbl            // EOF

	// empty line comment to exclude it from .String
	tokenCount_cbl //
)

// contains reports whether tok is in tokset.
// func contains(tokset uint64, tok tokens_pli) bool {
// 	return tokset&(1<<tok) != 0
// }

type Operator_cbl uint

//go:generate stringer -type Operator_cbl -linecomment tokens_pli.go

const (
	_ Operator_cbl = iota
)
