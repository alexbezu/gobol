package dc

import "github.com/alexbezu/gobol/pl"

var location_counter uint32
var LTORG = pl.CHAR(0xFFFF)

func F(init int) (ret *pl.Fixed_bin) {
	ret = pl.FIXED_BIN(31).BASED(LTORG)
	ret.SetOffset(location_counter)
	location_counter += 4
	ret.Set(init)
	return ret
}
