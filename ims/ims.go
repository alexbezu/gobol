//go:build immudb

package ims

import (
	"context"
	"log"
	"os"
	"strings"
	"time"

	"github.com/alexbezu/gobol/pl"
	"github.com/codenotary/immudb/pkg/api/schema"
	"google.golang.org/grpc"
)

// some kind of singleton for tcp
var conn *grpc.ClientConn
var err error
var c schema.ImmuServiceClient
var ctx context.Context
var cancel context.CancelFunc

func DLI(IMS_FUNC string, pcb pl.Objer, IO_AREA pl.Objer, SSAs ...string) {
	if conn == nil {
		Connect2db(os.Getenv("DBHOST") + ":" + os.Getenv("DBPORT"))
	}
	var iopcb = pl.NUMED(pl.NumT{
		"lterm_name": pl.CHAR(8),
		"reserved":   pl.CHAR(2),
		"status":     pl.CHAR(2),
		"date_time":  pl.CHAR(8),
		"msg_seq":    pl.FIXED_BIN(31),
		"mod_name":   pl.CHAR(8),
	}).BASED(pcb)

	var dbpcb = pl.NUMED(pl.NumT{
		"dbd_name":      pl.CHAR(8),
		"seg_level":     pl.CHAR(2), /* LEVEL OF LAST GOOD GET    */
		"status":        pl.CHAR(2),
		"proc_options":  pl.CHAR(4), /* PCB PROCESSING OPTION     */
		"dbd_reserved":  pl.FIXED_BIN(31),
		"seg_name":      pl.CHAR(8),       /* SEGMENT NAME LAST GOOD GET*/
		"length_fb_key": pl.FIXED_BIN(31), /* LENGTH(KEY_FEEDBACK)      */
		"num_sens_segs": pl.FIXED_BIN(31),
		"keyfba":        pl.CHAR(1),  /* Variable Length Field allocated by IMS */
		"key_fb_area":   pl.CHAR(50), /* FULLY CONCAT KEY LAST USED*/
	}).BASED(pcb)

	if dbpcb.I["seg_level"].String() != "" &&
		dbpcb.I["seg_level"].String() != "  " &&
		dbpcb.I["seg_level"].String() != "io" /*&& Number(dbpcb.I["seg_level"]*/ {

		//         let key = dbpcb.dbd_name;
		//         let segsize = 0;
		//         for (const ssa of SSAs) {
		//             const seglexed = /^([A-Z\s]{1,8})\(([A-Z]{1,8})[<=>\s]+(.+)\)/g.exec(ssa);
		//             const ssa_seg = seglexed[1];
		//             const ssa_field = seglexed[2];
		//             const ssa_value = seglexed[3];
		//             if (ssa_field == segs[ssa_seg]["key"]) {
		//                 key += `:${ssa_value}`;
		//                 const seg = segs[ssa_seg];
		//                 segsize = seg.base.size;
		//             } else {
		//                 throw "TODO:";
		//             }
		//             if (segsize == 0) {
		//                 iopcb.status.v = 'VG';
		//                 return;
		//             }
		//         }
		//         if (segsize == 0) segsize = 4096; //todo: segsize from another place

		//         const sab = new SharedArrayBuffer(Math.ceil(segsize/4)*4);
		//         const int32 = new Int32Array(sab);
		//         let buf = Buffer.from(sab);
		//         switch(String(IMS_FUNC)) {
		//             case "GHU ":
		//             case "GU  ":
		//                 dbpcb.key_fb_area.v = key;
		//                 worker.postMessage(["get", key, sab]);
		//                 Atomics.wait(int32, 0, 0);

		//                 if (IO_AREA instanceof pl.CHAR) {
		//                     buf.copy(IO_AREA.buf, 0, 0, IO_AREA.buf.length);
		//                     IO_AREA.bump_subs();
		//                 } else if (IO_AREA.constructor.name == "Object") {
		//                     let charbuf = new pl.CHAR(0);
		//                     charbuf.buf = buf;
		//                     charbuf.unpack(IO_AREA);
		//                 } else {
		//                     throw("redis.get: Unacceptable object: " + prop.constructor.name);
		//                 }

		//                 iopcb.status.v = '  ';

		//                 break;
		//             case "ISRT":
		//                 // ret = redis.set(key, IO_AREA.buf, redis.print);
		//                 break;
		//             case "REPL":
		//                 let replsize = 0;
		//                 if (IO_AREA instanceof pl.CHAR) {
		//                     replsize = IO_AREA.buf.length;
		//                     IO_AREA.buf.copy(buf, 0, 0, IO_AREA.buf.length);
		//                 } else if (IO_AREA.constructor.name === "Object" && IO_AREA.base !== undefined) {
		//                     replsize = IO_AREA.size;
		//                     IO_AREA.base.buf.copy(buf, 0, 0, replsize);
		//                 } else {
		//                     iopcb.status.v = 'BU';
		//                     throw("redis.rpush: Unacceptable object: " + prop.constructor.name);
		//                 }
		//                 worker.postMessage(["set", dbpcb.key_fb_area.native, sab, replsize]);
		//                 break;
		//             default:
		//                 throw "TODO:"
		//         }
	} else {
		switch IMS_FUNC {
		case "GU  ", "GN  ":
			var rpl *schema.MQpopReply
			rpl, err = c.MQpop(ctx, &schema.MQpopRequest{Qname: iopcb.I["mod_name"].String()})
			if err != nil {
				log.Fatal(err)
			}
			// fmt.Println(popresp)

			copy(*iopcb.GetBuff(), rpl.Value[0:8])
			copy(*IO_AREA.GetBuff(), rpl.Value[8:])

			iopcb.I["status"].Set("  ")
		case "ISRT":
			var put schema.MQputRequest

			lterm := pl.CHAR(8).INIT(strings.Repeat(" ", 8))
			if len(SSAs) > 0 { //set new format at first 8 bytes
				lterm.Set(SSAs[0])
			}
			put.Value = append(put.Value, *lterm.GetBuff()...)
			put.Value = append(put.Value, *IO_AREA.GetBuff()...)
			put.Qname = iopcb.I["lterm_name"].String()
			_, err := c.MQput(ctx, &put)
			if err != nil {
				log.Fatal(err)
			}
			// fmt.Println(putresp)
			iopcb.I["status"].Set("  ")
		case "PURG":
		default:
			panic("TODO: unknown IMS_FUNC")
		}
	}

}

func TDLI(parcnt int, IMS_FUNC string, pcb pl.Objer, IO_AREA pl.Objer, SSAs ...string) {
	DLI(IMS_FUNC, pcb, IO_AREA, SSAs...)
}

func Connect2db(endpoint string) {
	// host, port := os.Getenv("DBHOST"), os.Getenv("DBPORT")
	log.Print("Trying to connect to MQ: ")
	conn, err = grpc.Dial(endpoint, grpc.WithInsecure(), grpc.WithReturnConnectionError())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	} else {
		log.Print("Connected!")
	}

	c = schema.NewImmuServiceClient(conn)
	ctx, cancel = context.WithTimeout(context.Background(), time.Hour)
}

func closedb() {
	conn.Close()
	cancel()
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
