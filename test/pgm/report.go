package test222

import (
	"gobol/src/asm"
	"gobol/src/asm/ds"
	"gobol/src/pl"
)

// PRINT("ON", "NOGEN")
func REPORT() { // REPORT   CSECT
	asm.BALR(R12, 0)    //          BALR  R12,0                              SET UP MY
	asm.R[R12] = report //          USING *,R12                              BASE REGISTER
	// *        ST    R13,SAVEA+4                        ENSURE SAVE AREA
	// *        LA    R13,SAVEA                          CHAIN BUILT CORRECTLY.
	asm.LA(R15, pl.Int(0)) // LA    R15,0
	// asm.LA(R15, 0)         // LA    R15,0
	// *---------------------------------------------------------------------*
	asm.OPEN(INDCB, "INPUT") // INOPEN   OPEN  (INDCB,INPUT)
	// *
	asm.OPEN(OUTDCB, "OUTPUT") // OUOPEN   OPEN  (OUTDCB,OUTPUT)
	// *
	// * WRITE OUTPUT RECORD
	asm.PUT(OUTDCB, PTITLE) // WTITLE   PUT   OUTDCB,PTITLE
	asm.WTO(PTITLE)
	// *---------------------------------------------------------------------*
READREC: // READREC  GET   INDCB,PAYREC                       READ IN EMPLOYEE REC
	EODAD := asm.GET(INDCB, PAYREC)
	if EODAD {
		goto INCLOS
	}
	// *---------------------------------------------------------------------*
	// * PRINT ALSO ON SCREEN
	asm.WTO(PAYREC) //          WTO   PAYREC
	// *
	pl.MVC(PEMPID, EMPID)      // CPYSTUFF MVC   PEMPID,EMPID
	pl.MVC(PEMPLOYE, EMPLOYEE) //          MVC   PEMPLOYE,EMPLOYEE
	pl.MVC(PSALARY, SALARY)    //          MVC   PSALARY,SALARY
	// *
	asm.PACK(ZSALARY, SALARY) // PACKIT   PACK  ZSALARY,SALARY                     PACK SALARY
	asm.AP(ZTOTSAL, ZSALARY)  //          AP    ZTOTSAL,ZSALARY                    ADD MONTHLY WAGE TO
	asm.PUT(OUTDCB, OUTAREA)  // WRITEPR  PUT   OUTDCB,OUTAREA                     WRITE TO PRINTER
	goto READREC              //          B     READREC                            AND REPEAT TILL FILE
	// *---------------------------------------------------------------------*
INCLOS:
	pl.MVC(ATOTAL, EDWD) // INCLOS   MVC   ATOTAL,EDWD
	// *        LA    R1,ATOTAL
	asm.ED(ATOTAL, ZTOTSAL)   //          ED    ATOTAL,ZTOTSAL 615194
	asm.PUT(OUTDCB, TOTALLNE) //          PUT   OUTDCB,TOTALLNE                    PRINT TOTAL LINE
	// * PRINT TOTAL ALSO ON SCREEN
	asm.WTO(TOTALLNE) //          WTO   TOTALLNE
	// *
	asm.CLOSE(INDCB)  // CLSALL   CLOSE (INDCB)                            WE GET HERE FROM EODAD
	asm.CLOSE(OUTDCB) //          CLOSE (OUTDCB)
	// *
	return //          BR    R14
	//4192      JIN BLANCHARD          4476
	//4199      RON JACOBS             2113
	//0                        $    615,194     TOTAL MONTHLY WAGES
}

var report = pl.CL0(4096)

// var report [4096]byte TODO: try it and stash

//          LTORG
// *---------------------------------------------------------------------*
var INDCB = &asm.DCB{MACRF: "GM", DDNAME: "INDD", DSORG: "PS"}   // INDCB    DCB   MACRF=GM,DDNAME=INDD,DSORG=PS,EODAD=INCLOS
var OUTDCB = &asm.DCB{MACRF: "PM", DDNAME: "OUTDD", DSORG: "PS"} // OUTDCB   DCB   MACRF=PM,DDNAME=OUTDD,DSORG=PS
// *                          PAYROLL REPORT STRUCTURE
var PAYREC = pl.CL0(80)    // PAYREC   DS    0CL80                              HANDLE FOR THE STRU
var EMPID = pl.CHAR(4)     // EMPLOYEE ID
var rndsdfsf = pl.CHAR(6)  //          DS    CL6                                FILLER TO POSITION10
var EMPLOYEE = pl.CHAR(21) // EMPLOYEE DS    CL21                               NAME OF EMPLOYEE
var rndsgh56 = pl.CHAR(2)  //          DS    CL2                                FILLER TO POSITION34
var SALARY = pl.CHAR(4)    // SALARY   DS    CL4                                MONTHLY SALARY
var TOEND = pl.CHAR(43)    // TOEND    DS    CL43                               80 BYTES SO FAR
// *--------S-T-A-R-T----O-F----O-U-T-P-U-T----S-T-R-U-C-T-U-R-E---------*
var PTITLE = pl.CL0(121)                                        // PTITLE   DS    0CL121
var rnd2034nd = pl.CHAR(27).INIT(" P A Y R O L L  R E P O R T") //          DC    CL27' P A Y R O L L  R E P O R T'
var rndlksd84 = pl.CHAR(20).INIT("  -  B I M  C O R P.")        //          DC    CL20'  -  B I M  C O R P.'
var rndllkds7 = pl.CHAR(74).INIT(" ")                           //          DC    CL74' '
var OUTAREA = pl.CL0(133)                                       // OUTAREA  DS    0CL133
var EMPTY = pl.CHAR(1).INIT(" ")                                // EMPTY    DC    CL1' '
var PEMPID = pl.CHAR(4).INIT(" ")                               // PEMPID   DC    CL4' '
var rndsd49d = pl.CHAR(6).INIT(" ")                             //          DC    CL6' '
var PEMPLOYE = pl.CHAR(20).INIT(" ")                            // PEMPLOYE DC    CL20' '
var rndsdf7v = pl.CHAR(2).INIT(" ")                             //          DC    CL2' '
var PDOLLAR = pl.CHAR(1).INIT(" ")                              // PDOLLAR  DC    CL1' '
var PSALARY = pl.CHAR(5).INIT(" ")                              // PSALARY  DC    CL5' '
var OFILLER = pl.CHAR(94).INIT(" ")                             // OFILLER  DC    CL94' '
// *
var ZSALARY = pl.FIXED_DEC(5)                                                     // ZSALARY  DC    PL3'0'                             INITIALIZE SLARY PAC
var ZTOTSAL = pl.FIXED_DEC(9)                                                     // ZTOTSAL  DC    PL05'0'                            INITIALIZE TOTAL WA
var EDWD = ds.X(0x40, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x6B, 0x20, 0x20, 0x20) // EDWD     DC    X'402020202020206B202020'
var TOTALLNE = pl.CL0(133)                                                        // TOTALLNE DS    0CL133
var SKIP = pl.CHAR(1).INIT("0")                                                   // SKIP     DC    CL1'0'
var TFILL1 = pl.CHAR(9).INIT(" ")                                                 // TFILL1   DC    CL09' '
var TFILL2 = pl.CHAR(15).INIT(" ")                                                // TFILL2   DC    CL15' '
var TDOLLAR = pl.CHAR(1).INIT("$")                                                // TDOLLAR  DC    CL1'$'
var ATOTAL = pl.CHAR(11).INIT(" ")                                                // ATOTAL   DC    CL11' '
var TFILL3 = pl.CHAR(5).INIT(" ")                                                 // TFILL3   DC    CL5' '
var TOTMSG = pl.CHAR(61).INIT("TOTAL MONTHLY WAGES")                              // TOTMSG   DC    CL61'TOTAL MONTHLY WAGES'
var TFILL4 = pl.CHAR(30).INIT(" ")                                                // TFILL4   DC    CL30' '
// *
const R0 = 0
const R1 = 1
const R2 = 2
const R3 = 3
const R4 = 4
const R5 = 5
const R6 = 6
const R7 = 7
const R8 = 8
const R9 = 9
const R10 = 10
const R11 = 11
const R12 = 12
const R13 = 13
const R14 = 14
const R15 = 15

var _ = asm.END(report)
