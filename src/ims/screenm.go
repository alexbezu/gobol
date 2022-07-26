package ims

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"gobol/src/cmd/compile/internalll/syntax"
	"gobol/src/pl"
	"log"
	"strconv"
	"strings"

	"github.com/codenotary/immudb/pkg/api/schema"
)

// TN3270 Code table to transalte buffer addresses
var code_table = [...]uint8{0x40, 0xC1, 0xC2, 0xC3, 0xC4, 0xC5, 0xC6, 0xC7,
	0xC8, 0xC9, 0x4A, 0x4B, 0x4C, 0x4D, 0x4E, 0x4F,
	0x50, 0xD1, 0xD2, 0xD3, 0xD4, 0xD5, 0xD6, 0xD7,
	0xD8, 0xD9, 0x5A, 0x5B, 0x5C, 0x5D, 0x5E, 0x5F,
	0x60, 0x61, 0xE2, 0xE3, 0xE4, 0xE5, 0xE6, 0xE7,
	0xE8, 0xE9, 0x6A, 0x6B, 0x6C, 0x6D, 0x6E, 0x6F,
	0xF0, 0xF1, 0xF2, 0xF3, 0xF4, 0xF5, 0xF6, 0xF7,
	0xF8, 0xF9, 0x7A, 0x7B, 0x7C, 0x7D, 0x7E, 0x7F}

const (
	// Attributes
	attrsHI     uint8 = 0b11001000
	attrsALPHA  uint8 = 0b01000000
	attrsPROT   uint8 = 0b00100000
	attrsNUM    uint8 = 0b00010000
	attrsNODISP uint8 = 0b00001100
	attrsMOD    uint8 = 0b00000001
)

type field struct {
	posout     [2]byte
	label      string
	length     uint8
	attributes uint8
	value      *pl.Char
}

func NewTN3270screen() (ret TN3270screen) {
	ret.init()
	return ret
}

type TN3270screen struct {
	cur_dfld_type uint8
	DFLDs         [5]map[[2]byte]field
	DFLD          *map[[2]byte]field
	PFK           map[uint8]string
	PFKlabel      string
	MSGTYPE       string
	MFLDin        pl.Numed
	MFLDout       pl.Numed

	TRAN       string
	CURSOR     [2]byte
	Lterm      string
	seq_number uint16
}

func (s *TN3270screen) compile_tree(tree *syntax.File_asm) {
	for _, line := range tree.Lines {
		switch line.Instr {
		case "DEV":
			s.dev(line)
		case "DPAGE":
			s.dpage(line)
		case "DIV":
		case "DFLD":
			s.dfld(line, 0)
		case "DO": //s.do(line[2]);
		case "ENDDO": //s.do(line[2]);
		case "MFLD":
			s.mfld(line, 0)
		case "TITLE":
		case "FMT":
		case "FMTEND":
		case "EJECT":
		case "MSG":
			s.msg(line)
		case "MSGEND":
			s.msgend()
		case "SEG":
		case "END":
		case "PRINT":
		default:
			panic("compile_tree line.Instr")
		}
	}

	// 		 s.DFLD = s.DFLDs[s.cur_dfld_type];
	// 		 s.MFLDin.I.lterm.v = s.lterm;
}

// ['DEV', '',[('FEAT', ['IGNORE']),
//             ('PFK', ['LABEL2', 'PRTRANS1 F   ', 'PRTRANS1 2   ', '/FOR PRTRAN2.', 'PRTRANS1 21  ', 'PRTRANS1 24  ', None]),
//             ('DSCA',[144]),
//             ('SYSMSG', ['LABEL']),
//             ('TYPE', ['3270-A03'])]]
func (s *TN3270screen) dev(line syntax.Line) {
	// for k, values := range line.Params {
	for _, param := range line.Params {
		// k := values[0].ParamName
		switch param.ParamName {
		case "PFK":
			for i, value := range param.Values {
				if value.Extra != "" {
					index, _ := strconv.Atoi(value.Value) //PFK=(dfldname,10='Extra literal') where 'Extra literal' at index[10]
					s.PFK[uint8(index)] = value.Extra
				} else {
					index := PFKs[i]
					s.PFK[index] = value.Value //PFK=(dfldname,'literal') where 'literal' at index[1]
				}
			}
			//PFK[0] is a special case, where the dfldname does not have a label yet or label is hidden
			//label can be found alter in MFLD parameters
			//PFK=(dfldname, where dfldname is values[0]; see https://www.ibm.com/docs/en/ims/15.1.0?topic=statements-dev-statement
			s.PFKlabel = param.Values[0].Value
		case "TYPE":
			s.cur_dfld_type = 2
			// if (values == ['3270', '2'].join()) //values will 'join' automatically
			//	 this.cur_dfld_type = 2;
			// else if (values == ['3270-A03'].join())
			//	 this.cur_dfld_type = 3;
			s.DFLD = &s.DFLDs[s.cur_dfld_type]
		case "FEAT", "DSCA", "SYSMSG":
		default:
			panic("unkown DEV param")
		}
	}
}

// 'DPAGE CURSOR=((3,28,CURSOR)),FILL=PT'
// to
// this.CURSOR = POS(3,28)
func (s *TN3270screen) dpage(line syntax.Line) {
	for _, param := range line.Params {
		switch param.ParamName {
		case "CURSOR":
			r, _ := strconv.Atoi(param.Values[0].Value)
			c, _ := strconv.Atoi(param.Values[1].Value)
			s.CURSOR = s.POS(r, c)
		}
	}
}

// LBL   DFLD  'STR TEXT:',POS=(1,2),ATTR=(PROT),LTH=8
// to
// s.DFLD = {POS(1,2): [POSout(1,2), "LBL", 8, PROT, "STR TEXT"]}
func (s *TN3270screen) dfld(line syntax.Line, line_inc uint8) {
	pos := s.POS(1, 1)
	tmp := field{posout: s.POSout(1, 1), label: line.Label, length: 0, attributes: 0, value: &pl.Char{}}
	for _, param := range line.Params {
		switch param.ParamName {
		case "char":
			str := param.Values[0].Value
			tmp.value = pl.CHAR(uint32(len(str))).INIT(str)
			tmp.attributes |= 0b11110000
		case "POS":
			r, _ := strconv.Atoi(param.Values[0].Value)
			c, _ := strconv.Atoi(param.Values[1].Value)
			pos = s.POS(r+int(line_inc), c)
			tmp.posout = s.POSout(r+int(line_inc), c)
		case "LTH":
			length, _ := strconv.Atoi(param.Values[0].Value)
			tmp.length = uint8(length)
			tmp.value = pl.CHAR(uint32(length))
		case "ATTR":
			for _, attr := range param.Values {
				switch attr.Value {
				case "PROT":
					tmp.attributes |= attrsPROT
				case "NUM":
					tmp.attributes |= attrsNUM
				case "HI":
					tmp.attributes |= attrsHI
				case "NODISP":
					tmp.attributes |= attrsNODISP
					tmp.attributes &= 0b01111111
				case "MOD":
					tmp.attributes |= attrsMOD
				case "ALPHA":
					tmp.attributes |= attrsALPHA
				default:
					panic("dfld unknown ATTR")
				}

			}
		}
	}
	s.DFLDs[s.cur_dfld_type][pos] = tmp
}

// LBL    MFLD  (LABEL,'PRTRAN3 E    '),LTH=13
// to
// is.MFLDin.root = {POS(1,2):  pl.CHAR(13).INIT('PRTRAN3 E    ')}
func (s *TN3270screen) mfld(line syntax.Line, line_inc uint8) {
	var ATTR = false
	var STR string
	var label string
	// var FILL string
	var LEN int
	for _, param := range line.Params {
		switch param.ParamName {
		case "ID":
			label = param.Values[0].Value
		case "LTH":
			LEN, err = strconv.Atoi(param.Values[0].Value)
			if err != nil {
				panic("LTH= error in mfld")
			}
		case "FILL":
			// FILL = values[0].Value
		case "JUST":
		case "char":
			STR = param.Values[0].Value
		case "list":
			// MFLD (DATE,DATE2)
			// MFLD (ENTERLBL,'DEFAULT STR IF press enter')
			label = param.Values[0].Value
			STR = param.Values[1].Value
			if label == s.PFKlabel && STR != "" {
				s.PFK[aidENTER] = STR
			}
		case "ATTR":
			ATTR = (param.Values[0].Value == "YES")
		default:
			label = param.ParamName
		}
	}

	if label == "" {
		buff := make([]byte, 8)
		rand.Read(buff)
		label = base64.RawURLEncoding.EncodeToString(buff)
	}

	if s.MSGTYPE == "INPUT" {
		if LEN == 0 {
			s.MFLDin.I[label] = pl.CHAR(uint32(len(STR)))
		} else {
			s.MFLDin.I[label] = pl.CHAR(uint32(LEN))
		}
		if STR != "" {
			s.MFLDin.I[label].Set(STR)
		}
		// if (FILL) s.MFLDin.I[label].v = FILL.repeat(LEN);
	} else if s.MSGTYPE == "OUTPUT" {
		if LEN > 0 {
			if ATTR { // s.MFLDout.I[POS] = {"ATTR": new pl.CHAR(2), "STR": new pl.CHAR(LEN - 2)};
				s.MFLDout.I[label+"_attr"] = pl.CHAR(2)
				s.MFLDout.I[label] = pl.CHAR(uint32(LEN) - 2)

			} else {
				s.MFLDout.I[label] = pl.CHAR(uint32(LEN)) // s.MFLDout.I[POS] = {"STR": new pl.CHAR(LEN)};
			}
			if STR != "" {
				s.MFLDout.I[label].Set(STR)
			}
		} else {
			panic("MSGTYPE == OUTPUT no length")
		}
	} else {
		panic("MSGTYPE nor INPUT nor OUTPUT either")
	}
}

//['MSG', 'PROIFM', [('TYPE', ['INPUT']), ('SOR', ['PROIFD']), ('NXT', ['PRTEST1'])]]
//to s.MSGTYPE ='INPUT'
func (s *TN3270screen) msg(line syntax.Line) {
	for _, param := range line.Params {
		switch param.ParamName {
		case "TYPE":
			s.MSGTYPE = param.Values[0].Value
		case "SOR":
		case "NXT":
		case "FILL":
		default:
			panic("MSG Unknown param: " + param.ParamName)
		}
	}
}

/**
 * At the end of MSG block creates inited Object with based buffer to simplify transfering data
 *
 * MSGEND statement
 * The MSGEND statement terminates a message input or output definition and is required as the last statement in the definition.
 * If this is the end of the job submitted, it must also be followed by an END compilation statement.
 */
func (s *TN3270screen) msgend() {
	if s.MSGTYPE == "INPUT" {
		s.MFLDin = *pl.NUMED(s.MFLDin.I)       // scatter root for dynamicly created items
		bufflength := len(*s.MFLDin.GetBuff()) //strconv.Itoa()
		s.MFLDin.I["LL"].Set(bufflength)
	} else if s.MSGTYPE == "OUTPUT" {
		s.MFLDout = *pl.NUMED(s.MFLDout.I)
	}
}

func (s *TN3270screen) POSout(r, c int) (ret [2]byte) {
	address := (r-1)*80 + c - 2
	ret[0] = code_table[(address>>6)&0x3F]
	ret[1] = code_table[address&0x3F]
	return ret
}

func (s *TN3270screen) POS(r, c int) (ret [2]byte) {
	address := (r-1)*80 + c - 1
	ret[0] = code_table[(address>>6)&0x3F]
	ret[1] = code_table[address&0x3F]
	return ret
}

func (s *TN3270screen) MFLDsend(AID uint8) {
	// set special label (pass pressed PFK to MFLD if label exists here)
	if s.PFKlabel != "" {
		_, ok := s.MFLDin.I[s.PFKlabel]
		if ok {
			s.MFLDin.I[s.PFKlabel].Set(s.PFK[AID])
		}
	}
	// copy DFLD to MFLD
	for _, field := range *s.DFLD {
		_, ok := s.MFLDin.I[field.label]
		if ok {
			s.MFLDin.I[field.label].Set(field.value.String())
		}
	}
	s.MFLDin.I["lterm"].Set(s.Lterm)
	var put schema.MQputRequest
	put.Value = append(put.Value, *s.MFLDin.GetBuff()...)
	put.Qname = s.TRAN
	for len(put.Qname) < 8 {
		put.Qname += " "
	}
	_, err := c.MQput(ctx, &put)
	if err != nil {
		log.Fatal(err)
	}
}

func (s *TN3270screen) MFLDrecieve() {
	var rpl *schema.MQpopReply
	var get schema.MQpopRequest
	get.Qname = s.Lterm
	rpl, err = c.MQpop(ctx, &get)
	if err != nil {
		log.Fatal(err)
	}
	copy(*s.MFLDout.GetBuff(), rpl.Value)

	// this is message switch
	// if s.MFLDout.I["newformat"].String() != "" {
	// 			 this.readFormat(String(this.MFLDout.root.newformat));
	// 			 buf.copy(this.MFLDout.inArea.buf, 0, 0, this.MFLDout.inArea.size);
	// 			 this.MFLDout.inArea.bump_subs();
	// }

	//copy MFLD to DFLD
	for _, dfield := range *s.DFLD {
		mfld, ok := s.MFLDout.I[dfield.label]
		if ok {
			dfield.value.Set(mfld)
		}
	}
}

func (s *TN3270screen) init() { //clear in js
	s.cur_dfld_type = 2
	s.DFLDs = [5]map[[2]byte]field{}
	s.DFLDs[2] = map[[2]byte]field{}
	s.DFLD = &s.DFLDs[2]

	s.PFK = map[uint8]string{}
	s.PFKlabel = ""
	s.MSGTYPE = ""
	s.TRAN = strings.Repeat(" ", 8)
	s.CURSOR = s.POS(1, 2)
	s.seq_number = 0 // from TN3270Header field SEQ-NUMBER
	s.MFLDin = *pl.NUMED(pl.NumT{
		"lterm": pl.CHAR(8),
		"LL":    pl.FIXED_BIN(15),
		"ZZ":    pl.FIXED_BIN(15),
	})
	s.MFLDout = *pl.NUMED(pl.NumT{
		"newformat": pl.CHAR(8),
		"LL":        pl.FIXED_BIN(15),
		"ZZ":        pl.FIXED_BIN(15),
	})
}

// 	 do(Parameters) {
// 		 let [count, line_increment, position_increment] = Parameters[0];
// 		 if (!line_increment) line_increment = 1;
// 		 for (let i = 0; i < Number(count); i += line_increment) {
// 			 for (const line of Parameters[1]) {
// 				 switch(line[0]) {
// 				 case "DFLD":
// 					 this.dfld(line, i);
// 				 break;
// 				 case "MFLD":
// 					 this.mfld(line[1], line[2], i);
// 				 break;
// 				 }
// 			 }
// 		 }
// 	 }

func (s *TN3270screen) formatScreen() []byte {
	s.seq_number += 1
	DFLD := []byte{0, 0, 0, byte(s.seq_number >> 8), byte(s.seq_number), 0xF5, 0xC3}
	DFLD = append(DFLD, streamOrderSBA, s.CURSOR[0], s.CURSOR[1], streamOrderIC)
	for _, value := range *s.DFLD {
		DFLD = append(DFLD, streamOrderSBA, value.posout[0], value.posout[1], streamOrderSF, value.attributes)
		if value.value != nil {
			b := *value.value.GetBuff()
			DFLD = append(DFLD, b...)
		}
	}
	DFLD = append(DFLD, ptcIAC, tnEOR)
	return DFLD
}

func (s *TN3270screen) readFormat(screenfile string) {
	ast, err := syntax.ParseFileAsm("/home/oleksii/plexer/transactions/"+screenfile+".hlasm", nil)
	fmt.Println(ast)
	if err != nil {
		panic("bad format file")
	}
	s.init()
	s.compile_tree(ast)
	s.TRAN = screenfile
}
