package asm

import (
	"fmt"
	"reflect"
	"unsafe"

	"github.com/alexbezu/gobol/pl"
)

type Register interface {
	pl.Objer
}

var PSW uint64
var condition_code byte

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
	// return BC(0b1000)
	return condition_code == 0
}

func BH() bool {
	return BC(2)
}

func BNE() bool {
	return condition_code != 0
	// return BC(0b0111)
}

// 'BAS':        ('4D','R1,D2(X2,B2)',       'RX BD DD'),
// 'BASR':       ('0D','R1,R2',              'RR'),
// 'BASSM':      ('0C','R1,R2',              'RR'),
// 'BC':         ('47','M1,D2(X2,B2)',       'MX BD DD'),
func BC(M1 byte) bool {
	return (condition_code & M1) > 0
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
// The Compare instruction is used to compare a 32-bit signed in a R1, Operand 1, with a 32-bit signed in D2X2B2, Operand 2.
func C(R1 byte, D2 interface{}, X2, B2 byte) {
	var D2X2B2 uintptr
	var b *[4]byte
	switch d2 := D2.(type) {
	case int:
		D2X2B2 = uintptr(d2) + Rui[X2] + Rui[B2]
	case pl.Objer:
		hdr := (*reflect.SliceHeader)(unsafe.Pointer(d2.GetBuff()))
		b = (*[4]byte)(unsafe.Pointer(hdr.Data + uintptr(d2.GetOffset()) + Rui[X2] + Rui[B2]))
		// b := *d2.GetBuff()
		// o := d2.GetOffset() + uint32(Rui[X2]+Rui[B2])
		// _ = b[3+o]
		// D2X2B2 = uintptr(b[3+o]) | uintptr(b[2+o])<<8 | uintptr(b[1+o])<<16 | uintptr(b[0+o])<<24
		_ = b[3]
		D2X2B2 = uintptr(b[3]) | uintptr(b[2])<<8 | uintptr(b[1])<<16 | uintptr(b[0])<<24
	}
	cc := int32(Rui[R1]) - int32(D2X2B2)
	if cc == 0 {
		condition_code = 0
	} else if cc < 0 {
		condition_code = 1
	} else {
		condition_code = 2
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
func CR(R1, R2 byte) {
	cc := int32(Rui[R1]) - int32(Rui[R2])
	if cc == 0 {
		condition_code = 0
	} else if cc < 0 {
		condition_code = 1
	} else {
		condition_code = 2
	}
}

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
// Load value (32 bits) by address D2X2B2 to R1
func L(R1 byte, D2 interface{}, X2, B2 byte) {
	var b *[4]byte
	switch d2 := D2.(type) {
	case int:
		panic("TODO: L(R1 byte, D2 interface{}, X2, B2 byte) case int:")
	case uint32:
		if Rui[B2] > 0xFFFF {
			b = (*[4]byte)(unsafe.Pointer(Rui[B2] + uintptr(d2) + Rui[X2]))
		} else if Rui[X2] > 0xFFFF {
			//b = *(*[]byte)(unsafe.Pointer(Rui[X2] + uintptr(d2) + Rui[B2]))
			b = (*[4]byte)(unsafe.Add(unsafe.Pointer(Rui[X2]), uintptr(d2)+Rui[B2]))
		} else {
			panic("L(no address)")
		}
	case pl.Objer:
		// b := *d2.GetBuff()
		// o := d2.GetOffset() + uint32(Rui[X2]+Rui[B2])
		// _ = b[3+o]
		hdr := (*reflect.SliceHeader)(unsafe.Pointer(d2.GetBuff()))
		b = (*[4]byte)(unsafe.Pointer(hdr.Data + uintptr(d2.GetOffset()) + Rui[X2] + Rui[B2]))

	default:
		panic("L(R1 byte, D2 interface{}, X2, B2 byte)")
	}
	_ = b[3]
	Rui[R1] = uintptr(b[3]) | uintptr(b[2])<<8 | uintptr(b[1])<<16 | uintptr(b[0])<<24
}

// Load real Address of the buffer in to the R1
func LA(R1 byte, D2 interface{}, X2, B2 byte) {
	switch d2 := D2.(type) {
	case int:
		// X2 or B2 allready have an address of the Buff in uintptr
		Rui[R1] = uintptr(d2) + Rui[X2] + Rui[B2]
	case uint32:
		Rui[R1] = uintptr(d2) + Rui[X2] + Rui[B2]
	case pl.Objer:
		hdr := (*reflect.SliceHeader)(unsafe.Pointer(d2.GetBuff()))
		Rui[R1] = hdr.Data + uintptr(d2.GetOffset()) + Rui[X2] + Rui[B2]
		// Rui[R1] = uintptr(unsafe.Pointer(d2.GetBuff())) + uintptr(d2.GetOffset()) + Rui[X2] + Rui[B2]
	default:
		panic("LA(R1 byte, D2 interface{}, X2, B2 byte)")
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
func MVC1(D1 pl.Objer, D2 pl.Objer, length ...byte) {
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

func MVC(D1 interface{}, B1 byte, D2 interface{}, B2 byte, length ...byte) {
	var LL uintptr
	switch d := D1.(type) {
	case pl.Objer:
		LL = uintptr(d.GetSize())
	}
	if len(length) > 0 {
		LL = uintptr(length[0])
	}
	dest := calc_address(D1, B1, 0)
	source := calc_address(D2, B2, 0)
	for i := uintptr(0); i < LL; i++ {
		*(*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(dest)) + i)) = *(*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(source)) + i))
	}
}

func calc_address(D interface{}, X, B byte) (ret *byte) {
	switch d := D.(type) {
	case int:
		addr := uintptr(d) + Rui[X] + Rui[B]
		ret = (*byte)(unsafe.Pointer(addr + Rui[X] + Rui[B]))
	case uint32:
		addr := uintptr(d) + Rui[X] + Rui[B]
		ret = (*byte)(unsafe.Pointer(addr + Rui[X] + Rui[B]))
	case pl.Objer:
		hdr := (*reflect.SliceHeader)(unsafe.Pointer(d.GetBuff()))
		ret = (*byte)(unsafe.Pointer(hdr.Data + uintptr(d.GetOffset()) + Rui[X] + Rui[B]))
	default:
		panic("calc_address(D interface{}, X, B byte)")
	}
	return ret
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
//ST is used to copy the 4 (or 8) bytes from register 1 into the 4 bytes memory location specified by operand 2
func ST(R1 byte, D2 interface{}, X2, B2 byte) {
	var b *[4]byte
	switch d2 := D2.(type) {
	case int:
		if Rui[B2] > 0xFFFF {
			b = (*[4]byte)(unsafe.Pointer(Rui[B2] + uintptr(d2) + Rui[X2]))
		} else if Rui[X2] > 0xFFFF {
			b = (*[4]byte)(unsafe.Add(unsafe.Pointer(Rui[X2]), uintptr(d2)+Rui[B2]))
		} else {
			panic("L(no address)")
		}
	case pl.Objer:
		hdr := (*reflect.SliceHeader)(unsafe.Pointer(d2.GetBuff()))
		b = (*[4]byte)(unsafe.Pointer(hdr.Data + uintptr(d2.GetOffset()) + Rui[X2] + Rui[B2]))
	default:
		panic("ST(R1 byte, D2 interface{}, X2, B2 byte)")
	}
	_ = b[3]
	b[0] = byte(Rui[R1] >> 24)
	b[1] = byte(Rui[R1] >> 16)
	b[2] = byte(Rui[R1] >> 8)
	b[3] = byte(Rui[R1])
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
