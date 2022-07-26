package pl

import (
	"errors"
	"math"
	"strconv"
)

type Fixed_dec struct {
	Obj
	int
}

func FIXED_DEC(size int) (ret *Fixed_dec) {
	decsize := int(math.Ceil((float64(size) + 1) / 2))
	tmp := make([]byte, decsize)
	ret = &Fixed_dec{Obj: Obj{size: uint32(decsize), buff: &tmp}, int: 0}
	add2index(ret)
	return ret
}

func (d *Fixed_dec) writeBuf() {
	var str = strconv.Itoa(d.int)
	if len(str)%2 == 0 {
		str = "0" + str
	}
	if d.int < 0 {
		str += "D"
	} else {
		str += "C"
	}
	// var buf2, _ = strconv.ParseInt(str, 16, 64)
	// var offset = d.size - len(buf2)
	// this.buf.set(buf2, offset < 0 ? 0 : offset);
}

func (d *Fixed_dec) INIT(init int) *Fixed_dec {
	d.int = init
	d.writeBuf()
	return d
}

func (d *Fixed_dec) String() string {
	return strconv.Itoa(d.int)
}

func (d *Fixed_dec) clone() Objer {
	ret := FIXED_DEC(int(d.size * 2)).INIT(d.int)
	return ret
}

func (d *Fixed_dec) Set(val interface{}) error {
	switch val := val.(type) {
	case int:
		d.int = val
	case int32:
		d.int = int(val)
	case int16:
		d.int = int(val)
	case string:
		i, err := strconv.Atoi(val)
		if err != nil {
			return err
		}
		d.int = i
	case *Char:
		i, err := strconv.Atoi(val.String())
		if err != nil {
			return err
		}
		d.int = i
	default:
		return errors.New("BIN can accept int only")
	}
	d.writeBuf()
	return nil
}

func (d1 *Fixed_dec) Add_Sub_Mul_Packed_code(d2 *Fixed_dec, op byte) {
	switch op {
	case '+':
		d1.int += d2.int
	}
	d1.writeBuf()
}

func (c *Fixed_dec) P(val interface{}) Objer {
	switch val.(type) {

	}
	return c
}
