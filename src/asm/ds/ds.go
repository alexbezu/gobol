package ds

import "gobol/src/pl"

var dsect *pl.CL_0
var location_counter uint32

func DSECT() *pl.CL_0 {
	dsect = pl.CL0(1)
	location_counter = 0
	return dsect
}

func BufferSizeFromOffset(o pl.Objer) uint32 {
	return location_counter - o.GetOffset()
}

func X(init ...byte) (ret *pl.Char) {
	size := uint32(len(init))
	if size == 0 {
		size = 1
		init = make([]byte, size)
	}
	ret = pl.CHAR(size)
	if dsect != nil {
		ret.BASED(dsect)
	}
	ret.SetOffset(location_counter)
	location_counter += size
	ret.CopyBuff(init)
	return ret
}

func CL(size uint32) (ret *pl.Char) {
	ret = pl.CHAR(size)
	if dsect != nil {
		ret.BASED(dsect)
	}
	ret.SetOffset(location_counter)
	location_counter += size
	return ret
}

func H() (ret *pl.Fixed_bin) {
	ret = pl.FIXED_BIN(15)
	if dsect != nil {
		ret.BASED(dsect)
	}
	ret.SetOffset(location_counter)
	location_counter += 2
	return ret
}

func HL(size uint32) (ret *pl.Fixed_bin) {
	ret = pl.FIXED_BIN(15)
	if dsect == nil {
		panic("empty dsect")
	}
	ret.BASED(dsect)
	ret.SetOffset(location_counter)
	location_counter += size
	if dsect.GetSize() < location_counter {
		*dsect.GetBuff() = append(*dsect.GetBuff(), make([]byte, location_counter-dsect.GetSize())...)
	}
	return ret
}

func F() (ret *pl.Fixed_bin) {
	ret = pl.FIXED_BIN(31)
	if dsect == nil {
		panic("empty dsect")
	}
	ret.BASED(dsect)
	ret.SetOffset(location_counter)
	location_counter += 4
	if dsect.GetSize() < location_counter {
		*dsect.GetBuff() = append(*dsect.GetBuff(), make([]byte, location_counter-dsect.GetSize())...)
	}
	return ret
}

//set a new location counter
func ORG(o pl.Objer) bool {
	if dsect == nil {
		dsect = pl.CL0(o.GetSize())
	}
	dsect.BASED(o)
	location_counter = o.GetOffset()
	return false
}
