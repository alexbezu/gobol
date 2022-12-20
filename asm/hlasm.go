package asm

import (
	"fmt"
	"unsafe"

	"github.com/alexbezu/gobol/pl"
)

type Register interface {
	pl.Objer
}

var PSW uint64

const PSW_CC1 = uint64(1 << 18)
const PSW_CC2 = uint64(1 << 19)

var usings map[pl.Objer]byte

var R = [16]Register{pl.CL0(1), pl.CL0(1), pl.CL0(1), pl.CL0(1), pl.CL0(1), pl.CL0(1), pl.CL0(1), pl.CL0(1), pl.CL0(1), pl.CL0(1), pl.CL0(1), pl.CL0(1), pl.CL0(1), pl.CL0(1), pl.CL0(1), pl.CL0(1)}

var Rui = [16]uintptr{}

var Rup = [16]unsafe.Pointer{
	unsafe.Pointer(Rui[0]),
	unsafe.Pointer(Rui[1]),
	unsafe.Pointer(Rui[2]),
	unsafe.Pointer(Rui[3]),
	unsafe.Pointer(Rui[4]),
	unsafe.Pointer(Rui[5]),
	unsafe.Pointer(Rui[6]),
	unsafe.Pointer(Rui[7]),
	unsafe.Pointer(Rui[8]),
	unsafe.Pointer(Rui[9]),
	unsafe.Pointer(Rui[10]),
	unsafe.Pointer(Rui[11]),
	unsafe.Pointer(Rui[12]),
	unsafe.Pointer(Rui[13]),
	unsafe.Pointer(Rui[14]),
	unsafe.Pointer(Rui[15])}

var Rany = [16]interface{}{Rui[0], Rui[1], Rui[2], Rui[3]}

var asmbuf = [4096]byte{}

var REG byte

// 'A':          ('5A','R1,D2(X2,B2)',       'RX BD DD'),
// Add signed full-word at address D2(X2,B2) to value in register R1. CC: arithmetic, 3 on overflow.
// 'AH':         ('4A','R1,D2(X2,B2)',       'RX BD DD'),
// Add signed half-word at address D2(X2,B2) to value in register R1. CC: arithmetic, 3 on overflow
// 'AHI'          R1,I2 [A7A,RI]
// Add signed half-word constant I2 to value in register R1. To subtract a constant, use negative value. CC: arithmetic, 3 on overflow
// 'AL':         ('5E','R1,D2(X2,B2)',       'RX BD DD'),
// Add unsigned full-word at address D2(X2,B2) to unsigned value in register R1. CC: low bit of CC is set if the result is nonzero, high bit of CC is set if there was overflow (carry).
// 'ALR':        ('1E','R1,R2',              'RR'),
// Add unsigned value in register R2 to unsigned value in register R1. CC: same as AL instruction.
func ALR(R1 byte, R2 byte) {
	Rui[R1] += Rui[R2]
}

// 'AP':         ('FA','D1(L1,B1),D2(L2,B2)','L1L2 BD DD BD DD'),
// Add packed decimal number at address D2(B2) (of length L2 bytes) to packed decimal number at address D1(B1) (of length L1 bytes). CC: arithmetic, 3 on overflow.
func AP(D1 *pl.Fixed_dec, D2 *pl.Fixed_dec, len ...byte) {
	D1.Add_Sub_Mul_Packed_code(D2, '+')
}

// 'AR':         ('1A','R1,R2',              'RR'),
// Add signed value in register R2 to value in register R1. CC: arithmetic, 3 on overflow.
func AR(R1 byte, R2 byte) {
	// R[1] = R[1].P(R2)
	Rui[R1] = uintptr(int(Rui[R1]) + int(Rui[R2]))
}

// 'BAL':        ('45','R1,D2(X2,B2)',       'RX BD DD'),
// Store the address of the next instruction (from PSW) into register R1 and then branch to address D2(X2,B2).
// Depending on addressing mode, unused high bits of address stored into R1 will contain other information from PSW. CC: no change.
// 'BALR':       ('05','R1,R2',              'RR'),
// Similar to the BAL instruction, except the address of branch is taken from register R2 (but when R2 is zero, no branch is taken).
func BALR(R1 byte, R2 byte) bool {
	return false
}

func BE() bool {
	return (PSW&PSW_CC1 != PSW_CC1) && (PSW&PSW_CC2 != PSW_CC2)
}

func BH() bool {
	return false
}

func BNE() bool {
	return (PSW&PSW_CC1 == PSW_CC1) || (PSW&PSW_CC2 == PSW_CC2)
}

// 'BAS':        ('4D','R1,D2(X2,B2)',       'RX BD DD'),
// 'BASR':       ('0D','R1,R2',              'RR'),
// 'BASSM':      ('0C','R1,R2',              'RR'),
// 'BC':         ('47','M1,D2(X2,B2)',       'MX BD DD'),
func BC(M1 byte) bool {
	// Branch to address D2(X2,B2) depending on current condition code and 4-bit mask M1. Each bit of mask corresponds to possible value of CC.
	// The branch is taken if the bit of the mask corresponding to the current value of condition code is set.
	// Assembler provides additional mnemonics for various common mask types. CC: no change.
	return false
}

// 'BCR':        ('07','M1,R2',              'MR'),
// 'BCT':        ('46','R1,D2(X2,B2)',       'RX BD DD'),
// 'BCTR':       ('06','R1,R2',              'RR'),
func BCTR(R1 int, R2 int) bool {
	R1--
	if R2 != 0 {
		return true
	}
	return false
}

func BO() bool { // branch Ones
	return (PSW&PSW_CC1 == PSW_CC1) && (PSW&PSW_CC2 == PSW_CC2)
}

// 'BSM':        ('0B','R1,R2',              'RR'),
// 'BXH':        ('86','R1,R3,D2(B2)',       'RR BD DD'),
// 'BXLE':       ('87','R1,R3,D2(B2)',       'RR BD DD'),
// 'C':          ('59','R1,D2(X2,B2)',       'RX BD DD'),
func C(R1 byte, D2 interface{}, X2, B2 byte) {
	switch d2 := D2.(type) {
	case int:
		Rui[R1] = uintptr(d2) + Rui[X2] + Rui[B2]
	case pl.Objer:
		Rui[R1] = uintptr(unsafe.Pointer(&d2)) // + Rui[X2] + Rui[B2]
	case pl.Fixed_bin:
		Rui[R1] = uintptr(d2.I32()) + Rui[X2] + Rui[B2]
	}
}

// 'CDS':        ('BB','R1,R3,D2(B2)',       'RR BD DD'),
// 'CH':         ('49','R1,D2(X2,B2)',       'RX BD DD'),
// 'CL':         ('55','R1,D2(X2,B2)',       'RX BD DD'),
// 'CLC':        ('D5','D1(L,B1),D2(B2)',    'LL BD DD BD DD'),
func CLC(DB1 pl.Objer, DB2 pl.Objer, length ...byte) bool {
	LL := DB1.GetSize()
	if len(length) > 0 {
		LL = uint32(length[0])
	}
	destBuff := *DB1.GetBuff()
	sourceBuff := *DB2.GetBuff()
	for i := uint32(0); i < LL; i++ {
		cc := destBuff[int(DB1.GetOffset()+i)] - sourceBuff[int(DB2.GetOffset()+i)]
		if cc < 0 { // Condition Code 1
			PSW = PSW & PSW_CC1
			PSW = PSW &^ PSW_CC2
			return false
		} else if cc > 0 { // cc > 0 Condition Code 2
			PSW = PSW &^ PSW_CC1
			PSW = PSW & PSW_CC2
			return false
		}
	}
	return true
}

// 'CLCL':       ('0F','R1,R2',              'RR'),
// 'CLI':        ('95','D1(B1),I2',          'II BD DD'),
func CLI(D1B1 pl.Objer, I2 byte) bool {
	cc := D1B1.String()[0] - I2
	if cc == 0 { // Condition Code 0
		PSW = PSW &^ PSW_CC1
		PSW = PSW &^ PSW_CC2
		return true
	} else if cc < 0 { // Condition Code 1
		PSW = PSW & PSW_CC1
		PSW = PSW &^ PSW_CC2
	} else { // cc > 0 Condition Code 2
		PSW = PSW &^ PSW_CC1
		PSW = PSW & PSW_CC2
	}
	return false
}

// 'CLM':        ('BD','R1,M3,D2(B2)',       'RM BD DD'),
// 'CLR':        ('15','R1,R2',              'RR'),
// 'CP':         ('F9','D1(L1,B1),D2(L2,B2)','L1L2 BD DD BD DD'),
// 'CR':         ('19','R1,R2',              'RR'),
// 'CS':         ('BA','R1,R3,D2(B2)',       'RR BD DD'),
// 'CVB':        ('4F','R1,D2(X2,B2)',       'RX BD DD'),
// 'CVD':        ('4E','R1,D2(X2,B2)',       'RX BD DD'),
// 'D':          ('5D','R1,D2(X2,B2)',       'RX BD DD'),
// 'DP':         ('FD','D1(L1,B1),D2(L2,B2)','L1L2 BD DD BD DD'),
// 'DR':         ('1D','R1,R2',              'RR'),
// 'ED':         ('DE','D1(L,B1),D2(B2)',    'LL BD DD BD DD'),
func ED(D1 *pl.Char, D2 *pl.Fixed_dec, length ...byte) {
	pl.ED(D1, D2, length...)
}

// 'EDMK':       ('DF','D1(L,B1),D2(B2)',    'LL BD DD BD DD'),
// 'EX':         ('44','R1,D2(X2,B2)',       'RX BD DD'),
// 'IC':         ('43','R1,D2(X2,B2)',       'RX BD DD'),
// 'ICM':        ('BF','R1,M3,D2(B2)',       'RM BD DD'),
// 'L':          ('58','R1,D2(X2,B2)',       'RX BD DD'),
func L(R1 byte, D2 interface{}, X2, B2 byte) {
	switch d2 := D2.(type) {
	case int:
		Rui[R1] = uintptr(d2) + Rui[X2] + Rui[B2]
	case pl.Objer:
		Rui[R1] = uintptr(unsafe.Pointer(&d2)) // + Rui[X2] + Rui[B2]
	case pl.Fixed_bin:
		Rui[R1] = uintptr(d2.I32()) + Rui[X2] + Rui[B2]
	}
}

// 'LA':         ('41','R1,D2(X2,B2)',       'RX BD DD'),
// func LA(R1 byte, DXB2 pl.Objer) {
// 	switch r := R[R1].(type) {
// 	case *pl.Char:
// 		r.BASED(DXB2)
// 	case *pl.CL_0:
// 		r.BASED(DXB2)
// 	case *pl.Int:
// 	default:
// 		panic("TODO: asm LA")
// 	}
// 	R[R1] = R[R1].P(DXB2.GetOffset())
// }
func LA(R1 byte, D2 interface{}, X2, B2 byte) {
	// regs[_R1] = calc_address(_B2, _D2, _X2)
	switch d2 := D2.(type) {
	case int:
		Rui[R1] = uintptr(d2) + Rui[X2] + Rui[B2]
	case pl.Objer:
		displacement := d2.GetOffset()
		Rui[R1] = uintptr(displacement) + Rui[X2] + Rui[B2]
	}
}

// 'LCR':        ('13','R1,R2',              'RR'),
// 'LH':         ('48','R1,D2(X2,B2)',       'RX BD DD'),
func LH(R1 byte, DXB2 *pl.Fixed_bin) {
	R[R1] = DXB2
}

// 'LM':         ('98','R1,R3,D2(B2)',       'RR BD DD'),
// 'LNR':        ('11','R1,R2',              'RR'),
// 'LPR':        ('10','R1,R2',              'RR'),
// 'LR':         ('18','R1,R2',              'RR'),
// 'LTR':        ('12','R1,R2',              'RR'),
// 'M':          ('5C','R1,D2(X2,B2)',       'RX BD DD'),
// 'MH':         ('4C','R1,D2(X2,B2)',       'RX BD DD'),
// 'MP':         ('FC','D1(L1,B1),D2(L2,B2)','L1L2 BD DD BD DD'),
// 'MR':         ('1C','R1,R2',              'RR'),
// 'MVC':        ('D2','D1(L,B1),D2(B2)',    'LL BD DD BD DD'),
func MVC(D1 pl.Objer, D2 pl.Objer, length ...byte) {
	// MVC  WRIRANGE,=AL2(WRIEQ)    mvc(WRIRANGE,p(...))
	// MVC  0(L'EMPID,R1),EMPID     MVC(asm.R[R1].P(0), EMPID, len(EMPID))
	// MVC  TARGET(40),SOURCE       MVC(TARGET, SOURCE, 40)
	// MVC  TARGET+10(20),SOURCE+3  MVC(TARGET.P(10), SOURCE.P(3), 20) SubStringSet
	// MVC  DST(3,R14),SOURCE       MVC(asm.R[14].P(DST), SOURCE, 3) MOVE 3 BYTES FROM SOURCE TO THE ADDRESS BEGINNING AT REGISTER 14 PLUS DST displasement (offset)
	//                              MVC(DST.B1(R14), SOURCE, 3)
	//								MVC(DST, asm.R[14], SOURCE, asm.B2, 3)
	// MVC  TARGET+1(132),TARGET    MVC(TARGET.P(1), TARGET, 132) MOVES THE FIRST BLANK TO CHARACTER 2, WHICH MOVES TO CHAR. 3, THEN 4 AND SO ON
	// MVC BB(L'=C'RAY'),=C'RAY'
	LL := D1.GetSize()
	if len(length) > 0 {
		LL = uint32(length[0])
	}
	destBuff := *D1.GetBuff()
	sourceBuff := *D2.GetBuff()
	for i := uint32(0); i < LL; i++ {
		destBuff[int(D1.GetOffset()+i)] = sourceBuff[int(D2.GetOffset()+i)]
	}
	var _ = D1.String()
}

func MVC1(D1 pl.Objer, L byte, B1 byte, D2 pl.Objer, B2 byte) {
	dest := calc_address(B1, D1)
	source := calc_address(B2, D2)
	for i := uint32(0); i < uint32(L); i++ {
		asmbuf[dest+i] = asmbuf[source+i]
	}
	var _ = D1.String()
}

func calc_address(B byte, D pl.Objer, X ...byte) uint32 {
	addr := D.GetOffset()

	// if X != 0:
	//     addr = addr + cast_to_type(regs[X],int)

	// addr += cast_to_type(regs[B], int)

	return addr
}

// 'MVCIN':      ('E8','D1(L,B1),D2(B2)',    'LL BD DD BD DD'),
// 'MVCL':       ('0E','R1,R2',              'RR'),
// 'MVI':        ('92','D1(B1),I2',          'II BD DD'),
func MVI(D1 pl.Objer, I2 byte) {
	D1.Set(I2)
}

// 'MVN':        ('D1','D1(L,B1),D2(B2)',    'LL BD DD BD DD'),
// 'MVO':        ('F1','D1(L1,B1),D2(L2,B2)','L1L2 BD DD BD DD'),
// 'MVZ':        ('D3','D1(L,B1),D2(B2)',    'LL BD DD BD DD'),
// 'N':          ('54','R1,D2(X2,B2)',       'RX BD DD'),
// 'NC':         ('D4','D1(L,B1),D2(B2)',    'LL BD DD BD DD'),
// 'NI':         ('94','D1(B1),I2',          'II BD DD'),
// 'NR':         ('14','R1,R2',              'RR'),
// 'O':          ('56','R1,D2(X2,B2)',       'RX BD DD'),
// 'OC':         ('D6','D1(L,B1),D2(B2)',    'LL BD DD BD DD'),
// 'OI':         ('96','D1(B1),I2',          'II BD DD'),
// 'OR':         ('16','R1,R2',              'RR'),
// 'PACK':       ('F2','D1(L1,B1),D2(L2,B2)','L1L2 BD DD BD DD'),
func PACK(D1 *pl.Fixed_dec, D2 pl.Objer, len ...byte) {
	D1.Set(D2)
}

// 'S':          ('5B','R1,D2(X2,B2)',       'RX BD DD'),
func S(R1 byte, DXB2 pl.Objer) {
	// DXB2.Set(R[R1])
}

// 'SH':         ('4B','R1,D2(X2,B2)',       'RX BD DD'),
// 'SL':         ('5F','R1,D2(X2,B2)',       'RX BD DD'),
// 'SLA':        ('8B','R1,D2(X2,B2)',       'R0 BD DD'),
// 'SLDA':       ('8F','R1,D2(X2,B2)',       'R0 BD DD'),
// 'SLDL':       ('8D','R1,D2(X2,B2)',       'R0 BD DD'),
// 'SLL':        ('89','R1,D2(X2,B2)',       'R0 BD DD'),
// 'SLR':        ('1F','R1,R2',              'RR'),
// Subtract unsigned value in register R2 from unsigned value in register R1.
// CC: same as SL instruction.
func SLR(R1 byte, R2 byte) {
	Rui[R1] = Rui[R1] - Rui[R2]
}

// 'SP':         ('FB','D1(L1,B1),D2(L2,B2)','L1L2 BD DD BD DD'),
func SP(D1 pl.Objer, B1 byte, D2 pl.Objer, B2 byte, len ...byte) {
}

// 'SR':         ('1B','R1,R2',              'RR'),
// Subtract signed value in register R2 from value in register R1.
// CC: arithmetic, 3 on overflow.
func SR(R1 byte, R2 byte) {
	Rui[R1] = uintptr(int(Rui[R1]) - int(Rui[R2]))
}

// 'SRA':        ('8A','R1,D2(X2,B2)',       'R0 BD DD'),
// 'SRDA':       ('8E','R1,D2(X2,B2)',       'R0 BD DD'),
// 'SRDL':       ('8C','R1,D2(X2,B2)',       'R0 BD DD'),
// 'SRL':        ('88','R1,D2(X2,B2)',       'R0 BD DD'),
// 'SRP':        ('F0','D1(L1,B1),D2(B2),I3','LI BD DD BD DD'),
// 'ST':         ('50','R1,D2(X2,B2)',       'RX BD DD'),
func ST(R1 byte, DXB2 pl.Objer) {
	DXB2.Set(R[R1])
}

// 'STC':        ('42','R1,D2(X2,B2)',       'RX BD DD'),
// 'STCM':       ('BE','R1,M3,D2(B2)',       'RM BD DD'),
// 'STH':        ('40','R1,D2(X2,B2)',       'RX BD DD'),
// 'STM':        ('90','R1,R3,D2(B2)',       'RR BD DD'),
// 'SVC':        ('0A','I1',                 'II'),
// 'TM':         ('91','D1(B1),I2',          'II BD DD'),
func TM(D1B1 pl.Objer, I2 byte) {
	b := (*D1B1.GetBuff())[0]
	cc := (b & I2)
	if cc == 0 {
		PSW = PSW &^ PSW_CC1
		PSW = PSW &^ PSW_CC2
	} else if cc == I2 { // 1
		PSW = PSW & PSW_CC1
		PSW = PSW &^ PSW_CC2
	} else { // 3
		PSW = PSW & PSW_CC1
		PSW = PSW & PSW_CC2
	}
}

// 'TR':         ('DC','D1(L,B1),D2(B2)',    'LL BD DD BD DD'),
// 'TRT':        ('DD','D1(L,B1),D2(B2)',    'LL BD DD BD DD'),
// 'UNPK':       ('F3','D1(L1,B1),D2(L2,B2)','L1L2 BD DD BD DD'),
// 'X':          ('57','R1,D2(X2,B2)',       'RX BD DD'),
// 'XC':         ('D7','D1(L,B1),D2(B2)',    'LL BD DD BD DD'),
// 'XI':         ('97','D1(B1),I2',          'II BD DD'),
// 'XR':         ('17','R1,R2',              'RR'),
// 'ZAP':        ('F8','D1(L1,B1),D2(L2,B2)','L1L2 BD DD BD DD')

func WTO(D1 pl.Objer) {
	fmt.Println(D1.String())
}

func END(ltorg pl.Objer) (ret pl.Objer) {
	pl.ASMend(ltorg)
	return nil
}

func LINK(f func(), pointers ...pl.Objer) {
	r1 := &pl.Array{}
	r1.ArraySet(pointers)
	R[1] = r1
	f()
}

func USING(r byte) *pl.CL_0 {
	// R[15] = *pl.CL_0

	return R[r].(*pl.CL_0)
}

func USING2(o pl.Objer, r byte) {
	usings[o] = r
	if o != nil {
		switch o := o.(type) {
		case *pl.CL_0:
			o.BASED(R[r])
		}
	} else {

	}

}
