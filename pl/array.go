package pl

import (
	"sort"
)

type Array struct {
	Obj
	factor     uint8
	start, end int
	arr        []Objer
	I          func(int) Objer
}

func ARR(dim ...interface{}) (ret *Array) {
	ret = &Array{}
	ret.offset = 0
	ret.factor = uint8(len(dim) - 1)
	switch ret.factor {
	case 0:
		panic("0 dimensions? really?")
	case 1:
		ret.start = 1
		ret.end = dim[0].(int)
		obj := dim[1].(Objer)
		for i := ret.start; i < ret.end; i++ {
			ret.arr = append(ret.arr, obj.clone())
		}
		ret.I = func(i int) Objer {
			return ret.arr[i]
		}
	default:
		panic("gobol does not support more than 1 dimensions")
	}
	ret.scatter_root(ret)
	ret.size = ret.offset
	ret.offset = 0
	heap_buff := make([]byte, ret.size)
	ret.buff = &heap_buff
	ret.scatter_buff(ret.buff)
	add2index(ret)
	return ret
}

func (a *Array) scatter_root(root Objer) {
	switch a.factor {
	case 1:
		var indexes []int
		for _, prop := range a.arr {
			indexes = append(indexes, obj2indexMap[prop])
		}
		sort.Ints(indexes)

		for _, k := range indexes {
			prop := orderedObjList[k]
			switch prop := prop.(type) {
			case *Char:
				prop.root = root
				prop.offset = a.offset
				a.offset += prop.size
			case *Fixed_bin:
				prop.root = root
				prop.offset = a.offset
				a.offset += prop.size
			case *Array:
				prop.scatter_root(root)
			case *Numed:
				prop.root = root
				prop.offset = a.offset
				prop.scatter_root(root)
				prop.size = prop.offset - a.offset
				a.offset, prop.offset = prop.offset, a.offset
			default:
				panic("TODO: scatter_root on ARR")
			}
		}
	default:
		panic("plexer does not support more than 3 dimensions")
	}
}

func (a *Array) scatter_buff(b *[]byte) {
	for _, prop := range a.arr {
		switch prop := prop.(type) {
		case *Char:
			prop.buff = b
		case *Fixed_bin:
			prop.buff = b
		case *Numed:
			prop.buff = b
			prop.scatter_based_buff(b)
		case *Array:
			prop.buff = b
			prop.scatter_buff(b)
		default:
			panic("scatter_based_buff: Unknown object: prop")
		}
	}
}

func (a *Array) ArraySet(arr []Objer) (ret *Array) {
	a.arr = arr
	a.factor = 1
	return a
}

func (a *Array) String() string {
	return "Array"
}

//         this.ret.unpack = function(source) {
//             switch(factor) {
//             case 1:
//                 for (let prop of arr) {
//                     if (prop instanceof plObj) {
//                         source.buf.copy(prop.buf, 0, prop.offset, prop.offset + prop.size);
//                         prop.buf2native();
//                     } else if (prop.size) { //this is how plexer detects initedObj
//                         source.unpack(prop);
//                     } else throw "TODO: unpack on ARR ";
//                 }
//                 break;
//             case 2: throw "TODO: scatter_root case 2"; break;
//             case 3: throw "TODO: scatter_root case 3"; break;
//             default: throw "plexer does not support more than 3 dimensions";
//             }
//         }

func (a *Array) clone() Objer {
	return a
}

func (a *Array) P(val interface{}) Objer {
	switch val := val.(type) {
	case int:
		return a.arr[val/4]
	default:
		panic("TODO: array P")
	}
}
