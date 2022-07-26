package pl

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
)

type NumT map[string]Objer

type Numed struct {
	Obj
	I NumT
}

// constuctor for Numed (js:objInit(obj))
func NUMED(obj NumT) (ret *Numed) {
	ret = &Numed{I: obj}
	ret.offset = 0
	ret.scatter_root(ret)
	ret.size = ret.offset
	ret.offset = 0
	tmp_buff := make([]byte, ret.size)
	ret.buff = &tmp_buff
	ret.scatter_based_buff(ret.buff)
	add2index(ret)
	return ret
}

//js:obj.__proto__.BASED = function(based)
func (initedObj *Numed) BASED(based Objer) *Numed {
	initedObj.base = based
	initedObj.buff = based.GetBuff()
	initedObj.scatter_based_buff(initedObj.buff)

	switch based := based.(type) {
	case *Char:
		based.subs = append(based.subs, initedObj)
		// based.bump_subs()
	case *Fixed_bin:
		based.subs = append(based.subs, initedObj)
	case *Numed:
		based.bump_subs()
	case *Array:
	default:
		// if (based instanceof PTR) {
		//     initedObj_bump_subs(based);
	}
	return initedObj
}

/**
 * 1. calc size of an Object
 * 2. calc offset for each property
 * 3. point each property with its root
 */
func (that *Numed) scatter_root(root Objer) {
	// When iterating over a map with a range loop, the iteration order is not specified so additional index is being used
	var indexes []int
	for _, prop := range that.I {
		indexes = append(indexes, obj2indexMap[prop])
	}
	sort.Ints(indexes)

	for _, k := range indexes {
		prop := orderedObjList[k]
		switch prop := prop.(type) {
		case *Char:
			prop.root = root
			prop.offset = that.offset
			that.offset += prop.size
		case *Fixed_bin: //TODO: contribution to google
			prop.root = root
			prop.offset = that.offset
			that.offset += prop.size
		case *Numed:
			prop.root = root
			prop.offset = that.offset
			prop.scatter_root(root)
			prop.size = prop.offset - that.offset
			that.offset, prop.offset = prop.offset, that.offset
		case *Array:
			prop.scatter_root(root)
		default:
			panic("scatter_root: Unknown object: prop")
		}
	}
}

func (that *Numed) scatter_based_buff(b *[]byte) {
	for _, prop := range that.I {
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
			panic("scatter_based_buff: Array")
		default:
			panic("scatter_based_buff: Unknown object: prop")
		}
	}
}

//initedObj_bump_subs
func (initedObj *Numed) bump_subs() {
	for _, sub := range initedObj.subs {
		initedObj.pack(sub) //TODO: pack and unpack are the same, bump_subs must be only in obj.go
	}
}

//initedObj_pack
//copy each item of the Numed struct to a long buffer at the corresponding place
func (initedObj *Numed) pack(sub Objer) {
	for _, prop := range initedObj.I { //     for (const key of Object.keys(initedObj)) {
		sub := sub.(*Obj)
		switch prop := prop.(type) {
		case *Char:
			// copy(sub.buff[prop.offset:], prop.buff[:prop.size]) // prop.buf.copy(sub.buf, prop.offset, 0, prop.size);
		case *Fixed_bin:
			// copy(sub.buff[prop.offset:prop.offset+prop.size], prop.buff)
		case *Numed:
			prop.pack(sub) //             initedObj_pack(prop, sub);
		case *Array:
			panic("func (initedObj *Numed) pack(sub Objer) -> case *Array:")
		default:
			panic("func (initedObj *Numed) pack(sub Objer) -> default:")
		}
	}
}

func (n *Numed) String() string {
	j, err := json.Marshal(n.I)
	if err == nil {
		return "Cannot convert numed"
	}
	return string(j)
}

func (n *Numed) buf2native() {
	panic("buf2native")
}

func (n *Numed) clone() Objer {
	ret := &Numed{Obj: n.Obj}
	return ret
}

func (n *Numed) Set(interface{}) error {
	panic("Numed Set")

}

func (n *Numed) P(val interface{}) Objer {
	switch val.(type) {

	}
	return n
}

func InitNumed(st interface{}) { //TODO: reflect
	// v := reflect.ValueOf(s)
	// values := make([]interface{}, v.NumField())

	// for i := 0; i < v.NumField(); i++ {
	// 	values[i] = v.Field(i) //.Interface()
	// }
	// fmt.Println(values)
	s := reflect.ValueOf(st).Elem()
	typeOfT := s.Type()
	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		fmt.Printf("%d: %s %s = %v\n", i, typeOfT.Field(i).Name, f.Type(), f.Interface())
		// i := f.Interface()
		// i.(*pl.Char).INIT("FF")
	}
}
