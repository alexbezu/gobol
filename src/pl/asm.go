package pl

import (
	"strconv"
)

type CL_0 struct {
	Char
}

func CL0(size uint32) (ret *CL_0) {
	tmp := make([]byte, size)
	ret = &CL_0{Char: Char{Obj: Obj{size: size, buff: &tmp}}}
	add2index(ret)
	return ret
}

// func (c *CL_0) P(val interface{}) Objer {
// 	switch val.(type) {

// 	}
// 	return c
// }

// 'MVC':        ('D2','D1(L,B1),D2(B2)',    'LL BD DD BD DD'),
// MVC  0(L'EMPID,R1),EMPID     MVC(asm.R[R1].P(0), EMPID, asm.L(EMPID))
// MVC  TARGET(40),SOURCE       MVC(TARGET.P(40), SOURCE)
// MVC  TARGET+10(20),SOURCE+3  MVC(TARGET.P(40), SOURCE.P(3), 20) SubStringSet
// MVC  DST(3,R14),SOURCE       MVC(asm.R[14].P(DST), SOURCE, 3) MOVE 3 BYTES FROM SOURCE TO THE ADDRESS BEGINNING AT REGISTER 14 PLUS DST displasement (offset)
// MVC  TARGET+1(132),TARGET    MVC(TARGET.P(1), TARGET, 132) MOVES THE FIRST BLANK TO CHARACTER 2, WHICH MOVES TO CHAR. 3, THEN 4 AND SO ON
func MVC(D1 Objer, D2 Objer, length ...byte) {
	d1 := D1.(*Char)
	d2 := D2.(*Char)
	LL := d1.size
	if len(length) > 0 {
		LL = uint32(length[0])
	}
	destBuff := *D1.GetBuff()
	sourceBuff := *D2.GetBuff()
	for i := uint32(0); i < LL; i++ {
		destBuff[int(d1.offset+i)] = sourceBuff[int(d2.offset+i)]
	}
	D1.buf2native()
}

// 'ED':         ('DE','D1(L,B1),D2(B2)',    'LL BD DD BD DD'),
func ED(D1 *Char, D2 *Fixed_dec, length ...byte) {
	LL := D1.size
	if len(length) > 0 {
		LL = uint32(length[0])
	}
	destBuff := (*D1.buff)[D1.offset : D1.offset+D1.size]
	sourceVal := a2e(strconv.Itoa(D2.int))
	srcIdx := len(sourceVal) - 1
	fill_character := destBuff[0]

	for i := int(LL - 1); i >= 0; i-- {
		switch destBuff[i] {
		case 0x40:
			(*D1.buff)[int(D1.offset)+i] = 0x40
		case 0x20: // digit selector
			if srcIdx >= 0 {
				(*D1.buff)[int(D1.offset)+i] = sourceVal[srcIdx]
			} else {
				(*D1.buff)[int(D1.offset)+i] = fill_character
			}
			srcIdx--
		case 0x21: // digit selector and a significance starter
			if srcIdx >= 0 {
				(*D1.buff)[int(D1.offset)+i] = sourceVal[srcIdx]
			} else {
				(*D1.buff)[int(D1.offset)+i] = fill_character
			}
			srcIdx--
		case 0x6B: // EBCDIC code for a comma
			(*D1.buff)[int(D1.offset)+i] = 0x6B
		case 0x4B: // EBCDIC code for a decimal point
			(*D1.buff)[int(D1.offset)+i] = 0x4B
		case 0x60: // EBCDIC code for a minus sign
			if D2.int < 0 {
				(*D1.buff)[int(D1.offset)+i] = 0x60
			}
		case 0x5C: // EBCDIC symbol for an asterisk
			(*D1.buff)[int(D1.offset)+i] = 0x5C
		}
	}
	D1.buf2native()
}

func ASMend(ltorg Objer) {
	var offset uint32
	for _, obj := range orderedObjList {
		switch obj := obj.(type) {
		case *CL_0:
			obj.offset = offset
			copy((*ltorg.GetBuff())[offset:offset+obj.size], *obj.buff)
			obj.buff = ltorg.GetBuff()
			// no offset increment
		case *Char:
			obj.offset = offset
			copy((*ltorg.GetBuff())[offset:offset+obj.size], *obj.buff)
			obj.buff = ltorg.GetBuff()
			offset += obj.size
		case *Fixed_dec:
			obj.offset = offset
			copy((*ltorg.GetBuff())[offset:offset+obj.size], *obj.buff)
			obj.buff = ltorg.GetBuff()
			offset += obj.size
		case *Fixed_bin:
			obj.offset = offset
			copy((*ltorg.GetBuff())[offset:offset+obj.size], *obj.buff)
			obj.buff = ltorg.GetBuff()
			offset += obj.size
		default:
			panic("ASMend unsupported type")
		}
	}
}

