package asm

import (
	"fmt"

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

// 'A':          ('5A','R1,D2(X2,B2)',       'RX BD DD'),
// 'AH':         ('4A','R1,D2(X2,B2)',       'RX BD DD'),
// 'AL':         ('5E','R1,D2(X2,B2)',       'RX BD DD'),
// 'ALR':        ('1E','R1,R2',              'RR'),
func ALR(R1 int, R2 int) {

}

// 'AP':         ('FA','D1(L1,B1),D2(L2,B2)','L1L2 BD DD BD DD'),
func AP(D1 *pl.Fixed_dec, D2 *pl.Fixed_dec, len ...byte) {
	D1.Add_Sub_Mul_Packed_code(D2, '+')
}

// 'AR':         ('1A','R1,R2',              'RR'),
func AR(R1 int, R2 int) {
	R[1] = R[1].P(R2)
}

// 'BAL':        ('45','R1,D2(X2,B2)',       'RX BD DD'),
// 'BALR':       ('05','R1,R2',              'RR'),
func BALR(R1 int, R2 int) {

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
func C(R1 byte, DXB2 pl.Objer) {

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
func L(R1 byte, DXB2 pl.Objer) {
	R[R1] = DXB2
}

// 'LA':         ('41','R1,D2(X2,B2)',       'RX BD DD'),
func LA(R1 byte, DXB2 pl.Objer) {
	switch r := R[R1].(type) {
	case *pl.Char:
		r.BASED(DXB2)
	case *pl.CL_0:
		r.BASED(DXB2)
	case *pl.Int:
	default:
		panic("TODO: asm LA")
	}
	R[R1] = R[R1].P(DXB2.GetOffset())
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
	// MVC  0(L'EMPID,R1),EMPID     MVC(asm.R[R1].P(0), EMPID, asm.L(EMPID))
	// MVC  TARGET(40),SOURCE       MVC(TARGET.P(40), SOURCE)
	// MVC  TARGET+10(20),SOURCE+3  MVC(TARGET.P(40), SOURCE.P(3), 20) SubStringSet
	// MVC  DST(3,R14),SOURCE       MVC(asm.R[14].P(DST), SOURCE, 3) MOVE 3 BYTES FROM SOURCE TO THE ADDRESS BEGINNING AT REGISTER 14 PLUS DST displasement (offset)
	// MVC  TARGET+1(132),TARGET    MVC(TARGET.P(1), TARGET, 132) MOVES THE FIRST BLANK TO CHARACTER 2, WHICH MOVES TO CHAR. 3, THEN 4 AND SO ON
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
// 'SP':         ('FB','D1(L1,B1),D2(L2,B2)','L1L2 BD DD BD DD'),
// 'SR':         ('1B','R1,R2',              'RR'),
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
