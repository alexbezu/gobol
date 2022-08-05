// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package syntax

type tokens_pli uint

//go:generate stringer -type tokens_pli -linecomment tokens_pli.go

const (
	_    tokens_pli = iota
	_EOF            // EOF

	// names and literals
	_Name    // name
	_Literal // literal
	_Newline // newline

	// operators and operations
	// _Operator is excluding '*' (_Star)
	_Operator // op
	_AssignOp // op=
	_IncOp    // opop
	_Assign   // =
	// _Define   // :=
	_Star // *

	// delimiters
	_Lparen // (
	_Lbrack // [
	// _Lbrace    // {
	_Rparen // )
	_Rbrack // ]
	// _Rbrace    // }
	_Comma     // ,
	_Semi      // ;
	_Colon     // :
	_Dot       // .
	_DotDotDot // ...

	// keywords
	ZERODIVIDE         // ZERODIVIDE
	ZDIV               // ZDIV
	XU                 // XU
	XN                 // XN
	X                  // X
	WX                 // WX
	WRITE              // WRITE
	WINMAIN            // WINMAIN
	WHILE              // WHILE
	WHENEVER           // WHENEVER
	WHEN               // WHEN
	WAIT               // WAIT
	WIDECHAR           // WIDECHAR
	WCHAR              // WCHAR
	VSAM               // VSAM
	VS                 // VS
	VBS                // VBS
	VB                 // VB
	VARIABLE           // VARIABLE
	VARYINGZ           // VARYINGZ
	VARZ               // VARZ
	VARYING            // VARYING
	VAR                // VAR
	VALUERANGE         // VALUERANGE
	VALUELISTFROM      // VALUELISTFROM
	VALUELIST          // VALUELIST
	VALUE              // VALUE
	VALIDATE           // VALIDATE
	V                  // V
	UPTHRU             // UPTHRU
	UPDATE             // UPDATE
	UNTIL              // UNTIL
	UNSIGNED           // UNSIGNED
	UNS                // UNS
	UNLOCK             // UNLOCK
	UNION              // UNION
	UNDERFLOW          // UNDERFLOW
	UFL                // UFL
	UNDEFINEDFILE      // UNDEFINEDFILE
	UNDF               // UNDF
	UNCONNECTED        // UNCONNECTED
	UNBUFFERED         // UNBUFFERED
	UNBUFF             // UNBUFF
	UNALIGNED          // UNALIGNED
	UNAL               // UNAL
	U                  // U
	TYPE               // TYPE
	TSTACK             // TSTACK
	TRKOFL             // TRKOFL
	TRANSMIT           // TRANSMIT
	TRANSIENT          // TRANSIENT
	TP                 // TP
	TOTAL              // TOTAL
	TO                 // TO
	TITLE              // TITLE
	THREAD             // THREAD
	THEN               // THEN
	TASK               // TASK
	SYSTEM             // SYSTEM
	SUPPRESS           // SUPPRESS
	SUPPORT            // SUPPORT
	SUBSCRIPTRANGE     // SUBSCRIPTRANGE
	SUBRG              // SUBRG
	SUB                // SUB
	STRUCTURE          // STRUCTURE
	STRINGVALUE        // STRINGVALUE
	STRINGSIZE         // STRINGSIZE
	STRZ               // STRZ
	STRINGRANGE        // STRINGRANGE
	STRG               // STRG
	STRING             // STRING
	STREAM             // STREAM
	STORAGE            // STORAGE
	STOP               // STOP
	STDCALL            // STDCALL
	STATIC             // STATIC
	SNAP               // SNAP
	SKIP               // SKIP
	SIZE               // SIZE
	SIS                // SIS
	SIGNED             // SIGNED
	SIGNAL             // SIGNAL
	SET                // SET
	SEQUENTIAL         // SEQUENTIAL
	SEQL               // SEQL
	SELECT             // SELECT
	SCALARVARYING      // SCALARVARYING
	REWRITE            // REWRITE
	REVERT             // REVERT
	REUSE              // REUSE
	RETURNS            // RETURNS
	RETURN             // RETURN
	RETCODE            // RETCODE
	RESIGNAL           // RESIGNAL
	RESERVES           // RESERVES
	RESERVED           // RESERVED
	REREAD             // REREAD
	REPLY              // REPLY
	REPLACE            // REPLACE
	REPEAT             // REPEAT
	REORDER            // REORDER
	RENAME             // RENAME
	RELEASE            // RELEASE
	REGIONAL           // REGIONAL
	REFER              // REFER
	REENTRANT          // REENTRANT
	REDUCIBLE          // REDUCIBLE
	RED                // RED
	RECURSIVE          // RECURSIVE
	RECSIZE            // RECSIZE
	RECORD             // RECORD
	REAL               // REAL
	READ               // READ
	RANGE              // RANGE
	R                  // R
	PUT                // PUT
	PROCEDURE          // PROCEDURE
	PROC               // PROC
	PRIORITY           // PRIORITY
	PRINT              // PRINT
	PRECISION          // PRECISION
	PREC               // PREC
	POSITION           // POSITION
	POS                // POS
	POINTER            // POINTER
	PTR                // PTR
	PICTURE            // PICTURE
	PIC                // PIC
	PENDING            // PENDING
	PASSWORD           // PASSWORD
	PARAMETER          // PARAMETER
	PARM               // PARM
	PAGESIZE           // PAGESIZE
	PAGE               // PAGE
	PACKED             // PACKED
	PACKAGE            // PACKAGE
	P                  // P
	OVERFLOW           // OVERFLOW
	OFL                // OFL
	OUTPUT             // OUTPUT
	OUTONLY            // OUTONLY
	OTHERWISE          // OTHERWISE
	OTHER              // OTHER
	ORDINAL            // ORDINAL
	ORDER              // ORDER
	OPTLINK            // OPTLINK
	OPTIONS            // OPTIONS
	OPTIONAL           // OPTIONAL
	OPEN               // OPEN
	ON                 // ON
	OFFSET             // OFFSET
	NULLINIT           // NULLINIT
	NOZERODIVIDE       // NOZERODIVIDE
	NOZDIV             // NOZDIV
	NOWRITE            // NOWRITE
	NOUNDERFLOW        // NOUNDERFLOW
	NOUFL              // NOUFL
	NOTE               // NOTE
	NOSUBSCRIPTRANGE   // NOSUBSCRIPTRANGE
	NOSUBRG            // NOSUBRG
	NOSTRINGSIZE       // NOSTRINGSIZE
	NOSTRZ             // NOSTRZ
	NOSTRINGRANGE      // NOSTRINGRANGE
	NOSTRG             // NOSTRG
	NOSIZE             // NOSIZE
	NORMAL             // NORMAL
	NOPRINT            // NOPRINT
	NOOVERFLOW         // NOOVERFLOW
	NOOFL              // NOOFL
	NONVARYING         // NONVARYING
	NONVAR             // NONVAR
	NONE               // NONE
	NONCONNECTED       // NONCONNECTED
	NONCONN            // NONCONN
	NONASSIGNABLE      // NONASSIGNABLE
	NONASGN            // NONASGN
	NOLOCK             // NOLOCK
	NOINVALIDOP        // NOINVALIDOP
	NOINLINE           // NOINLINE
	NOINIT             // NOINIT
	NOFIXEDOVERFLOWOFL // NOFIXEDOVERFLOW
	NOEXECOPS          // NOEXECOPS
	NODESCRIPTOR       // NODESCRIPTOR
	NOCONVERSION       // NOCONVERSION
	NOCONV             // NOCONV
	NOCHECK            // NOCHECK
	NOCHARGRAPHIC      // NOCHARGRAPHIC
	NOCHARG            // NOCHARG
	NCP                // NCP
	NAME               // NAME
	MAIN               // MAIN
	M                  // M
	LOOP               // LOOP
	LOCATE             // LOCATE
	LOCAL              // LOCAL
	LITTLEENDIAN       // LITTLEENDIAN
	LIST               // LIST
	LINKAGE            // LINKAGE
	LINESIZE           // LINESIZE
	LINE               // LINE
	LIMITED            // LIMITED
	LIKE               // LIKE
	LEAVE              // LEAVE
	LABEL              // LABEL
	KEYTO              // KEYTO
	KEYLOC             // KEYLOC
	KEYLENGTH          // KEYLENGTH
	KEYFROM            // KEYFROM
	KEYED              // KEYED
	KEY                // KEY
	ITERATE            // ITERATE
	IRREDUCIBLE        // IRREDUCIBLE
	IRRED              // IRRED
	INVALIDOP          // INVALIDOP
	INTO               // INTO
	INTERACTIVE        // INTERACTIVE
	INTER              // INTER
	INTERNAL           // INTERNAL
	INT                // INT
	INPUT              // INPUT
	INOUT              // INOUT
	INONLY             // INONLY
	INLINE             // INLINE
	INITIAL            // INITIAL
	INIT               // INIT
	INI                // INI
	INDFOR             // INDFOR
	INDEXED            // INDEXED
	INDEXAREA          // INDEXAREA
	INCLUDE            // INCLUDE
	IN                 // IN
	IMPORTED           // IMPORTED
	IGNORE             // IGNORE
	IF                 // IF
	IEEE               // IEEE
	I                  // I
	HEXADEC            // HEXADEC
	HANDLE             // HANDLE
	GX                 // GX
	GRAPHIC            // GRAPHIC
	GOTO               // GOTO
	GO                 // GO
	GET                // GET
	GENKEY             // GENKEY
	GENERIC            // GENERIC
	G                  // G
	FS                 // FS
	FROMALIEN          // FROMALIEN
	FROM               // FROM
	FREE               // FREE
	FORTRAN            // FORTRAN
	FORMAT             // FORMAT
	FOREVER            // FOREVER
	FORCE              // FORCE
	FLUSH              // FLUSH
	FLOAT              // FLOAT
	FIXEDOVERFLOW      // FIXEDOVERFLOW
	FOFL               // FOFL
	FIXED              // FIXED
	FINISH             // FINISH
	FILE               // FILE
	FETCHABLE          // FETCHABLE
	FETCH              // FETCH
	FBS                // FBS
	FB                 // FB
	F                  // F
	EXTERNAL           // EXTERNAL
	EXT                // EXT
	EXPORTS            // EXPORTS
	EXIT               // EXIT
	EXEC               // EXEC
	EXCLUSIVE          // EXCLUSIVE
	EXCL               // EXCL
	EVENT              // EVENT
	ERROR              // ERROR
	ENVIRONMENT        // ENVIRONMENT
	ENV                // ENV
	ENTRY              // ENTRY
	ENDPAGE            // ENDPAGE
	ENDFILE            // ENDFILE
	END                // END
	ELSE               // ELSE
	EDIT               // EDIT
	// // [eE]                                                return 'E'
	DOWNTHRU     // DOWNTHRU
	DO           // DO
	DISPLAY      // DISPLAY
	DIRECT       // DIRECT
	DIMACROSS    // DIMACROSS
	DIMENSION    // DIMENSION
	DIM          // DIM
	DETACH       // DETACH
	DESCRIPTORS  // DESCRIPTORS
	DESCRIPTOR   // DESCRIPTOR
	DELETE       // DELETE
	DELAY        // DELAY
	DEFINE       // DEFINE
	DEFAULT      // DEFAULT
	DFT          // DFT
	DEFINED      // DEFINED
	DEF          // DEF
	DECLARE      // DECLARE
	DCL          // DCL
	DECIMAL      // DECIMAL
	DEC          // DEC
	DEACTIVATE   // DEACTIVATE
	DEACT        // DEACT
	DB           // DB
	DATE         // DATE
	DATA         // DATA
	D            // D
	CTLASA       // CTLASA
	CTL360       // CTL360
	COPY         // COPY
	CONVERSION   // CONVERSION
	CONV         // CONV
	CONTROLLED   // CONTROLLED
	CTL          // CTL
	CONSTANT     // CONSTANT
	CONST        // CONST
	CONSECUTIVE  // CONSECUTIVE
	CONNECTED    // CONNECTED
	CONN         // CONN
	CONDITION    // CONDITION
	COND         // COND
	COMPLEX      // COMPLEX
	CPLX         // CPLX
	COLUMN       // COLUMN
	COL          // COL
	COBOL        // COBOL
	CLOSE        // CLOSE
	CHECK        // CHECK
	CHARGRAPHIC  // CHARGRAPHIC
	CHARG        // CHARG
	CHARACTER    // CHARACTER
	CHAR         // CHAR
	CELL         // CELL
	CDECL        // CDECL
	CALL         // CALL
	C            // C
	BYVALUE      // BYVALUE
	BYADDR       // BYADDR
	BY           // BY
	BX           // BX
	BUILTIN      // BUILTIN
	BUFSP        // BUFSP
	BUFNI        // BUFNI
	BUFND        // BUFND
	BUFFOFF      // BUFFOFF
	BUFFERS      // BUFFERS
	BUFF         // BUFF
	BUFFERED     // BUFFERED
	BUF          // BUF
	BLKSIZE      // BLKSIZE
	BKWD         // BKWD
	BIT          // BIT
	BINARY       // BINARY
	BIN          // BIN
	BIGENDIAN    // BIGENDIAN
	BEGIN        // BEGIN
	BASED        // BASED
	BACKWARDS    // BACKWARDS
	B4           // B4
	B3           // B3
	B2           // B2
	B1           // B1
	B            // B
	AUTOMATIC    // AUTOMATIC
	AUTO         // AUTO
	ATTENTION    // ATTENTION
	ATTN         // ATTN
	ATTACH       // ATTACH
	ASSIGNABLE   // ASSIGNABLE
	ASGN         // ASGN
	ASSEMBLER    // ASSEMBLER
	ASM          // ASM
	ASCII        // ASCII
	AREA         // AREA
	ANYCONDITION // ANYCONDITION
	ANYCOND      // ANYCOND
	ALLOCATE     // ALLOCATE
	ALLOC        // ALLOC
	ALIGNED      // ALIGNED
	ALIAS        // ALIAS
	ADDBUFF      // ADDBUFF
	ACTIVATE     // ACTIVATE
	ACT          // ACT
	ABNORMAL     // ABNORMAL
	A            // A

	// empty line comment to exclude it from .String
	tokenCount //
)

// contains reports whether tok is in tokset.
func contains(tokset uint64, tok tokens_pli) bool {
	return tokset&(1<<tok) != 0
}

type LitKind uint8

// TODO(gri) With the 'i' (imaginary) suffix now permitted on integer
//           and floating-point numbers, having a single ImagLit does
//           not represent the literal kind well anymore. Remove it?
const (
	IntLit LitKind = iota
	FloatLit
	ImagLit
	RuneLit
	StringLit
)

type Operator_pli uint

//go:generate stringer -type Operator_pli -linecomment tokens_pli.go

const (
	_    Operator_pli = iota
	NEQ1              // <>
	NEQ2              // ^=
	NEQ3              // ~=
	GE1               // ^<
	GE2               // ~<
	GE3               // >=
	LE1               // ^>
	LE2               // ~>
	LE3               // <=
	// [\+\-\*\|\!\&]\=     return 'SELFOP'
	HANDLEPTR // =>
	PTRPTR    // ->
	POWER     // **
	NOT       // ^
	MOT1      // ~
	AND       // &
	CONCAT    // ||
	CONCAT1   // !!
	OR        // |

	// Def is the : in :=
	// Def // :
	// Not   // !
	// Recv  // <-
	// Tilde // ~

	// precOrOr
	// OrOr // ||

	// precAndAnd
	// AndAnd // &&

	// precCmp
	// Eql // ==
	// Neq // !=
	Lss // <
	Leq // <=
	Gtr // >
	Geq // >=

	// precAdd
	Add // +
	Sub // -
	// Or  // |
	// Xor // ^

	// precMul
	Mul // *
	Div // /
	Rem // %
	And // &
	// AndNot // &^
	Shl // <<
	Shr // >>
)
