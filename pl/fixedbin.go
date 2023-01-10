package pl

import (
	"errors"
	"strconv"
)

type Fixed_bin struct {
	Obj
	int32
}

func FIXED_BIN(size uint32) (ret *Fixed_bin) {
	if size <= 7 {
		size = 1
	} else if size <= 15 {
		size = 2
	} else if size <= 31 {
		size = 4
	} else {
		panic("wrong size of FIXED_BIN")
	}
	tmp := make([]byte, size)
	ret = &Fixed_bin{Obj: Obj{size: size, buff: &tmp}, int32: 0}
	add2index(ret)
	return ret
}

func (c *Fixed_bin) buf2native() {
	b := *c.buff
	switch c.size {
	case 1:
		_ = b[0+c.offset]
		c.int32 = int32(b[0+c.offset])
	case 2:
		_ = b[1+c.offset]
		c.int32 = int32(uint16(b[1+c.offset]) | uint16(b[0+c.offset])<<8)
	case 4:
		_ = b[3+c.offset]
		c.int32 = int32(b[3+c.offset]) | int32(b[2+c.offset])<<8 | int32(b[1+c.offset])<<16 | int32(b[0+c.offset])<<24
	default:
		panic("wrong size of FIXED_BIN")
	}
}

func (c *Fixed_bin) writeBuf() {
	b := *c.buff
	switch c.size {
	case 1:
		_ = b[0+c.offset]
		b[0+c.offset] = byte(c.int32)
	case 2:
		// binary.BigEndian.PutUint16
		_ = b[1+c.offset]
		b[0+c.offset] = byte(c.int32 >> 8)
		b[1+c.offset] = byte(c.int32)
	case 4:
		// binary.BigEndian.PutUint32
		_ = b[3+c.offset]
		b[0+c.offset] = byte(c.int32 >> 24)
		b[1+c.offset] = byte(c.int32 >> 16)
		b[2+c.offset] = byte(c.int32 >> 8)
		b[3+c.offset] = byte(c.int32)
	default:
		panic("wrong size of FIXED_BIN")
	}
}

func (c *Fixed_bin) BASED(based Objer) *Fixed_bin {
	c.base = based
	c.buff = nil
	c.buff = based.GetBuff()
	return c
}

func (c *Fixed_bin) INIT(init int32) *Fixed_bin {
	c.int32 = init
	c.writeBuf()
	return c
}

func (c *Fixed_bin) String() string {
	return strconv.Itoa(int(c.int32))
}

func (c *Fixed_bin) I32() int32 {
	c.buf2native()
	return (c.int32)
}

func (c *Fixed_bin) clone() Objer {
	ret := FIXED_BIN(c.size*8 - 1).INIT(c.int32)
	return ret
}

func (c *Fixed_bin) Set(val interface{}) error {
	switch val := val.(type) {
	case int32:
		c.int32 = val
	case int:
		c.int32 = int32(val)
	case int16:
		c.int32 = int32(val)
	case string:
		i, err := strconv.Atoi(val)
		if err != nil {
			return err
		}
		c.int32 = int32(i)
	case *Char:
		i, err := strconv.Atoi(val.String())
		if err != nil {
			return err
		}
		c.int32 = int32(i)
	default:
		return errors.New("BIN can accept int only")
	}
	c.writeBuf()
	return nil
}

func (c *Fixed_bin) P(val interface{}) Objer {
	ret := &Fixed_bin{Obj: c.Obj}
	ret.int32 = c.int32
	switch val := val.(type) {
	case int:
		ret.offset += uint32(val)
	case uint32:
		ret.offset += val
	}
	return ret
}
