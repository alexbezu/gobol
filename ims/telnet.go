package ims

import (
	"crypto/rand"
	"encoding/base64"
	"log"
	"net"
	"os"
	"strings"

	"github.com/alexbezu/gobol/pl"
)

const (
	//  Telnet protocol commands
	ptcSE uint8 = 240 // End of subnegotiation parameters
	ptcSB uint8 = 250 // Sub-option to follow
	// ptcWILL uint8 = 251  // Will; request or confirm option begin
	// ptcWONT uint8 = 252  // Wont; deny option request
	ptcDO uint8 = 253 // Do = Request or confirm remote option
	// ptcDONT uint8 = 254  // Don't = Demand or confirm option halt
	ptcIAC uint8 = 0xff // Interpret as Command
	// ptcSEND uint8 = 0x1  // Sub-process negotiation SEND command
	// ptcIS   uint8 = 0x0  // Sub-process negotiation IS command

	// TN3270 Telnet Commands
	tnASSOCIATE  uint8 = 0
	tnCONNECT    uint8 = 1
	tnDEVICETYPE uint8 = 2
	tnFUNCTIONS  uint8 = 3
	tnIS         uint8 = 4
	// tnREASON     uint8 = 5
	// tnREJECT     uint8 = 6
	// tnREQUEST    uint8 = 7
	// tnRESPONSES  uint8 = 2
	tnSEND uint8 = 8
	tnEOR  uint8 = 0xef //End of Record

	// TN3270 Stream Orders
	streamOrderSF uint8 = 0x1D
	// streamOrderSFE uint8 = 0x29
	streamOrderSBA uint8 = 0x11
	// streamOrderSA  uint8 = 0x28
	// streamOrderMF  uint8 = 0x2C
	streamOrderIC uint8 = 0x13
	// streamOrderPT  uint8 = 0x5
	// streamOrderRA  uint8 = 0x3C
	// streamOrderEUA uint8 = 0x12
	// streamOrderGE  uint8 = 0x8

	// subsommands
	// subcmdBINARY  uint8 = 0
	// subcmdEOR     uint8 = 25
	// subcmdTTYPE   uint8 = 24
	subcmdTN3270E uint8 = 40

	// TN3270 Attention Identification (AIDS)
	aidENTER uint8 = 0x7d
	aidPF1   uint8 = 0xf1
	aidPF2   uint8 = 0xf2
	aidPF3   uint8 = 0xf3
	aidPF4   uint8 = 0xf4
	aidPF5   uint8 = 0xf5
	aidPF6   uint8 = 0xf6
	aidPF7   uint8 = 0xf7
	aidPF8   uint8 = 0xf8
	aidPF9   uint8 = 0xf9
	aidPF10  uint8 = 0x7a
	aidPF11  uint8 = 0x7b
	aidPF12  uint8 = 0x7c
	aidPF13  uint8 = 0xc1
	aidPF14  uint8 = 0xc2
	aidPF15  uint8 = 0xc3
	aidPF16  uint8 = 0xc4
	aidPF17  uint8 = 0xc5
	aidPF18  uint8 = 0xc6
	aidPF19  uint8 = 0xc7
	aidPF20  uint8 = 0xc8
	aidPF21  uint8 = 0xc9
	aidPF22  uint8 = 0x4a
	aidPF23  uint8 = 0x4b
	aidPF24  uint8 = 0x4c
	aidOICR  uint8 = 0xe6
	// aidMSR_MHS uint8 = 0xe7
	// aidSELECT  uint8 = 0x7e
	// aidNO      uint8 = 0x60
	// aidQREPLY  uint8 = 0x61
	// aidPA1     uint8 = 0x6c
	// aidPA2     uint8 = 0x6e
	// aidPA3     uint8 = 0x6b
	aidCLEAR uint8 = 0x6d
	// aidSYSREQ  uint8 = 0xf0

	// DATA_TYPE
	// dt3270_DATA    uint8 = 0x00 // The data portion of the message contains only the 3270 data stream.
	// dtSCS_DATA     uint8 = 0x01 // The data portion of the message contains SNA Character Stream data.
	// dtRESPONSE     uint8 = 0x02 // The data portion of the message constitutes device-status information and the RESPONSE-FLAG field indicates whether this is a positive or negative response (see below).
	dtBIND_IMAGE uint8 = 0x03 // The data portion of the message is the SNA bind image from the session established between the server and the host application.
	// dtUNBIND       uint8 = 0x04 // The data portion of the message is an Unbind reason code.
	// dtNVT_DATA     uint8 = 0x05 // The data portion of the message is to be interpreted as NVT data.
	// dtREQUEST      uint8 = 0x06 // There is no data portion present in the message.  Only the REQUEST-FLAG field has any meaning.
	// dtSSCP_LU_DATA uint8 = 0x07 // The data portion of the message is data from the SSCP-LU session.

	ERR_COND_CLEARED uint8 = 0x0
)

var PFKs = [...]uint8{aidENTER, aidPF1, aidPF2, aidPF3, aidPF4, aidPF5, aidPF6, aidPF7, aidPF8, aidPF9, aidPF10, aidPF11, aidPF12, aidPF13, aidPF14,
	aidPF15, aidPF16, aidPF17, aidPF18, aidPF19, aidPF20, aidPF21, aidPF22, aidPF23, aidPF24, aidOICR}

func TN3270Eserver() {
	l, err := net.Listen("tcp", ":23567")
	// l, err := net.ListenTCP("tcp", &net.TCPAddr{IP: []byte{0, 0, 0, 0}, Port: 23567})
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	Connect2db(os.Getenv("DBHOST") + ":" + os.Getenv("DBPORT"))
	defer closedb()

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go func(c net.Conn) {
			screen := NewTN3270screen()
			tn3270eNegotiation(c, &screen)
			data := make([]byte, 2048)
			for {
				l, _ := c.Read(data)
				if screenHandler(data[:l], &screen) {
					c.Write(screen.formatScreen())
				} else {
					c.Close()
					return
				}
			}
		}(conn)
	}
}

func tn3270eNegotiation(socket net.Conn, screen *TN3270screen) {

	screen.readFormat("CALC") //HLWRLD DSN8IPD
	// make unique lterm
	buff := make([]byte, 6)
	rand.Read(buff)
	screen.Lterm = base64.RawURLEncoding.EncodeToString(buff)[:8]

	response := make([]byte, 2048)
	negotiation := map[string][]byte{
		"start":       {ptcIAC, ptcDO, subcmdTN3270E},
		"DEVICE_TYPE": {ptcIAC, ptcSB, subcmdTN3270E, tnSEND, tnDEVICETYPE, ptcIAC, ptcSE},
		"DEVICE_TYPE2": append(append(append(append([]byte{ptcIAC, ptcSB, subcmdTN3270E, tnDEVICETYPE, tnIS},
			[]byte("IBM-3278-2-E")...),
			tnCONNECT),
			[]byte(screen.Lterm)...),
			ptcIAC, ptcSE),
		"FUNCTIONS_REQUEST": {ptcIAC, ptcSB, subcmdTN3270E, tnFUNCTIONS, tnIS, tnASSOCIATE, tnDEVICETYPE, tnIS, ptcIAC, ptcSE},
		"BIND_IMAGE": append(append([]byte{dtBIND_IMAGE, ERR_COND_CLEARED, 0, 0, 0},
			[]byte("TELNET")...),
			ptcIAC, tnEOR),
	}

	socket.Write(negotiation["start"])
	socket.Read(response)
	// negotiation_start:
	// 	// assert(data == IAC + WILL + Subcommand['TN3270E'])
	socket.Write(negotiation["DEVICE_TYPE"])
	socket.Read(response)
	// negotiation_DEVICE_TYPE:
	// 	// assert(data[:5] == IAC + SB + Subcommand['TN3270E'] + TN_DEVICETYPE + TN_REQUEST)
	// 	// assert(data[-2:] == IAC + SE)
	socket.Write(negotiation["DEVICE_TYPE2"])
	socket.Read(response)
	// negotiation_DEVICE_TYPE2:
	// 	// assert(data == IAC + SB + Subcommand['TN3270E'] + TN_FUNCTIONS + TN_REQUEST + TN_ASSOCIATE + TN_DEVICETYPE + TN_IS + IAC + SE)
	socket.Write(negotiation["FUNCTIONS_REQUEST"])
	socket.Write(negotiation["BIND_IMAGE"])
	socket.Write(screen.formatScreen())
}

func screenHandler(arrdata []byte, screen *TN3270screen) bool {

	if len(arrdata) <= 5 {
		return false
	}

	arrdata = arrdata[5:] // git rid of the header

	AID := arrdata[0]
	if AID == aidCLEAR {
		screen.init()
		return true
	}

	//scatter input
	screen.CURSOR[0] = arrdata[1]
	screen.CURSOR[1] = arrdata[2]
	fields := strings.Split(string(arrdata[3:len(arrdata)-2]), string(streamOrderSBA))
	for _, f := range fields {
		if len(f) < 3 {
			continue
		}
		field := []byte(f)
		var POS [2]byte
		POS[0], POS[1] = field[0], field[1]
		dfld, ok := (*screen.DFLD)[POS]
		if ok {
			dfld.value.CopyBuff(field[2:])
		}
	}

	if AID == aidENTER {
		input := pl.CHAR(uint32(len(arrdata)))
		input.CopyBuff(arrdata)
		formatted := strings.ToUpper(input.String())
		i := strings.Index(formatted, "/FOR ")
		if i >= 0 {
			screnname := strings.TrimSpace(formatted[i+5 : len(formatted)-2])
			screen.readFormat(screnname)
			return true
		}
		i = strings.Index(formatted, "/QUIT")
		if i >= 0 {
			return false
		}
	}

	screen.MFLDsend(AID)

	screen.MFLDrecieve()
	return true
}
