//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"os"

	"github.com/alexbezu/gobol/cmd/compile/internalll/syntax"
	"github.com/alexbezu/gobol/cmd/compile/internalll/translate"
	"github.com/alexbezu/gobol/ims"
)

func main() {
	// ast, _ := syntax.ParseFileAsm("../test/dcds.hlasm", nil)
	// fmt.Println(ast)
	/* var iopcb = &struct {
		Lterm_name *pl.Char
		Reserved   *pl.Char
		Status     *pl.Char
		Date_time  *pl.Char
		Msg_seq    *pl.Char
		Mod_name   *pl.Char
	}{
		Lterm_name: pl.CHAR(8).INIT("DEFAULT"),
		Reserved:   pl.CHAR(2).INIT("io"),
		Status:     pl.CHAR(2).INIT("  "),
		Date_time:  pl.CHAR(8).INIT("20220124"),
		Msg_seq:    pl.CHAR(4),
		Mod_name:   pl.CHAR(8),
	}
	pl.InitNumed(iopcb) */

	// 	dcl   1 month ( 12 ),   /* means January ... December */
	//         2 income   dec fixed ( 7 , 2 ),
	//         2 outgo    dec fixed ( 7 , 2 );
	//    month ( 1 ) . income = 1234.56;   /* store income-value of January  */
	//    month . income ( 2 ) = 2345.67;   /* store income-value of February */
	//    month ( 3 ) . outgo  = 3456.78;   /* store  outgo-value of March    */
	//    month . outgo ( 4 )  = 4567.89;   /* store  outgo-value of April    */

	// var month = [12]pl.NumT{
	// { "income": pl.FIXED_BIN(15),
	// 	 "outgo": pl.FIXED_BIN(15)},
	// }

	/* var month = pl.ARR("1:12", pl.NUMED(pl.NumT{
		"income": pl.ARR("7", pl.FIXED_BIN(15)),
		"outgo":  pl.ARR("7", pl.FIXED_BIN(15)),
	})) */

	// dcl   1 year ( 1999 : 2019 ),
	// 2 month ( 12 ),
	//   3 value ( 3 , 3 )   bin fixed (15);}
	/* var year = [2019 - 1999]pl.NumT{"month": [12]pl.NumT{
		"value": [3][3]pl.FIXED_BIN(15),
	}} */

	// var A = pl.ARR("9", "1:12", pl.NUMED(pl.NumT{
	// 	"LL":   pl.FIXED_BIN(15),
	// 	"ZZ":   pl.FIXED_BIN(15),
	// 	"Tran": pl.CHAR(8).INIT("TRANS001"),
	// 	"Addr": pl.CHAR(8).INIT("Ukraine2"),
	// }))

	if len(os.Args) < 2 {
		fmt.Println("Usage: ./gobol tn3270e or ./gobol asm ...")
		return
	}
	switch os.Args[1] {
	case "tn3270e":
		ims.TN3270Eserver()
	case "asm":
		switch os.Args[2] {
		case "run":
			file := os.Args[3]
			ast, _ := syntax.ParseFileAsm(file, nil)
			var tr translate.Translator_asm
			tr.Precompile_tree(ast)
			tr.Compile_tree(ast)
			fmt.Println(tr.Src)
		case "save":
		default:
		}
	default:
	}
}