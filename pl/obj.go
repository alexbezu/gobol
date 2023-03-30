package pl

type Objer interface {
	String() string
	buf2native()
	GetBuff() *[]byte
	CopyBuff([]byte)
	Set(interface{}) error
	clone() Objer
	P(interface{}) Objer
	GetOffset() uint32
	GetSize() uint32
}

//sub string set
// type sss struct {
// 	pos    uint32
// 	length uint32
// }

type Obj struct {
	root   Objer   // for numed and Array
	base   Objer   // imitation BASED keyword
	subs   []Objer // pubsubtree pattern
	size   uint32
	buff   *[]byte // for EBCDIC and BigEndian
	pbuff  **[]byte
	offset uint32 // for numed and Array
	// sss    interface{} // for the SUBSTR function
}

var orderedObjList []Objer
var obj2indexMap = map[Objer]int{}

func add2index(item Objer) {
	orderedObjList = append(orderedObjList, item)
	obj2indexMap[item] = len(orderedObjList) - 1
}

func (c *Obj) String() string {
	panic("Obj string")
}

func (c *Obj) buf2native() {
	panic("Obj buf2native")
}

func (c *Obj) GetBuff() *[]byte {
	if c.pbuff != nil {
		return *c.pbuff
	}
	return c.buff
}

func (c *Obj) Set(interface{}) error {
	panic("Obj Set")
}

func (c *Obj) clone() Objer {
	panic("Obj clone")
}

func (c *Obj) CopyBuff(b []byte) {
	copy((*c.buff)[c.offset:c.offset+c.size], b)
}

func (o *Obj) P(val interface{}) Objer {
	panic("Obj P")
}

func (o *Obj) SetOffset(shift uint32) {
	o.offset = shift
}

func (o *Obj) GetOffset() uint32 {
	return o.offset
}

func (o *Obj) GetSize() uint32 {
	return o.size
}

// type Buff []byte

// func (b *Buff) String() string {
// 	return "sdfsdf"
// }
// func (b *Buff) buf2native() {

// }

// func (b *Buff) GetBuff() (ret *[]byte) {
// 	ret = b
// 	copy(*ret, *b)
// 	return ret
// }

// func (b *Buff) CopyBuff([]byte)
// func (b *Buff) Set(interface{}) error
// func (b *Buff) clone() Objer

type Int int

func (i Int) String() string {
	return "sdfsdf"
}
func (i Int) buf2native() {

}
func (i Int) GetBuff() (ret *[]byte) {
	return ret
}
func (i Int) CopyBuff([]byte) {
	panic("Obj Set")
}
func (i Int) Set(interface{}) error {
	panic("Obj Set")
}
func (i Int) clone() Objer {
	panic("Obj Set")
}

func (i Int) P(val interface{}) Objer {
	panic("Obj P")
}

func (i Int) GetOffset() uint32 {
	return 0
}

func (i Int) GetSize() uint32 {
	panic("Obj Set")
}
