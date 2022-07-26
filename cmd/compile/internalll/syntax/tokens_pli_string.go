// Code generated by "stringer -type tokens_pli -linecomment tokens_pli.go"; DO NOT EDIT.

package syntax

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[_EOF-1]
	_ = x[_Name-2]
	_ = x[_Literal-3]
	_ = x[_Newline-4]
	_ = x[_Operator-5]
	_ = x[_AssignOp-6]
	_ = x[_IncOp-7]
	_ = x[_Assign-8]
	_ = x[_Star-9]
	_ = x[_Lparen-10]
	_ = x[_Lbrack-11]
	_ = x[_Rparen-12]
	_ = x[_Rbrack-13]
	_ = x[_Comma-14]
	_ = x[_Semi-15]
	_ = x[_Colon-16]
	_ = x[_Dot-17]
	_ = x[_DotDotDot-18]
	_ = x[ZERODIVIDE-19]
	_ = x[ZDIV-20]
	_ = x[XU-21]
	_ = x[XN-22]
	_ = x[X-23]
	_ = x[WX-24]
	_ = x[WRITE-25]
	_ = x[WINMAIN-26]
	_ = x[WHILE-27]
	_ = x[WHENEVER-28]
	_ = x[WHEN-29]
	_ = x[WAIT-30]
	_ = x[WIDECHAR-31]
	_ = x[WCHAR-32]
	_ = x[VSAM-33]
	_ = x[VS-34]
	_ = x[VBS-35]
	_ = x[VB-36]
	_ = x[VARIABLE-37]
	_ = x[VARYINGZ-38]
	_ = x[VARZ-39]
	_ = x[VARYING-40]
	_ = x[VAR-41]
	_ = x[VALUERANGE-42]
	_ = x[VALUELISTFROM-43]
	_ = x[VALUELIST-44]
	_ = x[VALUE-45]
	_ = x[VALIDATE-46]
	_ = x[V-47]
	_ = x[UPTHRU-48]
	_ = x[UPDATE-49]
	_ = x[UNTIL-50]
	_ = x[UNSIGNED-51]
	_ = x[UNS-52]
	_ = x[UNLOCK-53]
	_ = x[UNION-54]
	_ = x[UNDERFLOW-55]
	_ = x[UFL-56]
	_ = x[UNDEFINEDFILE-57]
	_ = x[UNDF-58]
	_ = x[UNCONNECTED-59]
	_ = x[UNBUFFERED-60]
	_ = x[UNBUFF-61]
	_ = x[UNALIGNED-62]
	_ = x[UNAL-63]
	_ = x[U-64]
	_ = x[TYPE-65]
	_ = x[TSTACK-66]
	_ = x[TRKOFL-67]
	_ = x[TRANSMIT-68]
	_ = x[TRANSIENT-69]
	_ = x[TP-70]
	_ = x[TOTAL-71]
	_ = x[TO-72]
	_ = x[TITLE-73]
	_ = x[THREAD-74]
	_ = x[THEN-75]
	_ = x[TASK-76]
	_ = x[SYSTEM-77]
	_ = x[SUPPRESS-78]
	_ = x[SUPPORT-79]
	_ = x[SUBSCRIPTRANGE-80]
	_ = x[SUBRG-81]
	_ = x[SUB-82]
	_ = x[STRUCTURE-83]
	_ = x[STRINGVALUE-84]
	_ = x[STRINGSIZE-85]
	_ = x[STRZ-86]
	_ = x[STRINGRANGE-87]
	_ = x[STRG-88]
	_ = x[STRING-89]
	_ = x[STREAM-90]
	_ = x[STORAGE-91]
	_ = x[STOP-92]
	_ = x[STDCALL-93]
	_ = x[STATIC-94]
	_ = x[SNAP-95]
	_ = x[SKIP-96]
	_ = x[SIZE-97]
	_ = x[SIS-98]
	_ = x[SIGNED-99]
	_ = x[SIGNAL-100]
	_ = x[SET-101]
	_ = x[SEQUENTIAL-102]
	_ = x[SEQL-103]
	_ = x[SELECT-104]
	_ = x[SCALARVARYING-105]
	_ = x[REWRITE-106]
	_ = x[REVERT-107]
	_ = x[REUSE-108]
	_ = x[RETURNS-109]
	_ = x[RETURN-110]
	_ = x[RETCODE-111]
	_ = x[RESIGNAL-112]
	_ = x[RESERVES-113]
	_ = x[RESERVED-114]
	_ = x[REREAD-115]
	_ = x[REPLY-116]
	_ = x[REPLACE-117]
	_ = x[REPEAT-118]
	_ = x[REORDER-119]
	_ = x[RENAME-120]
	_ = x[RELEASE-121]
	_ = x[REGIONAL-122]
	_ = x[REFER-123]
	_ = x[REENTRANT-124]
	_ = x[REDUCIBLE-125]
	_ = x[RED-126]
	_ = x[RECURSIVE-127]
	_ = x[RECSIZE-128]
	_ = x[RECORD-129]
	_ = x[REAL-130]
	_ = x[READ-131]
	_ = x[RANGE-132]
	_ = x[R-133]
	_ = x[PUT-134]
	_ = x[PROCEDURE-135]
	_ = x[PROC-136]
	_ = x[PRIORITY-137]
	_ = x[PRINT-138]
	_ = x[PRECISION-139]
	_ = x[PREC-140]
	_ = x[POSITION-141]
	_ = x[POS-142]
	_ = x[POINTER-143]
	_ = x[PTR-144]
	_ = x[PICTURE-145]
	_ = x[PIC-146]
	_ = x[PENDING-147]
	_ = x[PASSWORD-148]
	_ = x[PARAMETER-149]
	_ = x[PARM-150]
	_ = x[PAGESIZE-151]
	_ = x[PAGE-152]
	_ = x[PACKED-153]
	_ = x[PACKAGE-154]
	_ = x[P-155]
	_ = x[OVERFLOW-156]
	_ = x[OFL-157]
	_ = x[OUTPUT-158]
	_ = x[OUTONLY-159]
	_ = x[OTHERWISE-160]
	_ = x[OTHER-161]
	_ = x[ORDINAL-162]
	_ = x[ORDER-163]
	_ = x[OPTLINK-164]
	_ = x[OPTIONS-165]
	_ = x[OPTIONAL-166]
	_ = x[OPEN-167]
	_ = x[ON-168]
	_ = x[OFFSET-169]
	_ = x[NULLINIT-170]
	_ = x[NOZERODIVIDE-171]
	_ = x[NOZDIV-172]
	_ = x[NOWRITE-173]
	_ = x[NOUNDERFLOW-174]
	_ = x[NOUFL-175]
	_ = x[NOTE-176]
	_ = x[NOSUBSCRIPTRANGE-177]
	_ = x[NOSUBRG-178]
	_ = x[NOSTRINGSIZE-179]
	_ = x[NOSTRZ-180]
	_ = x[NOSTRINGRANGE-181]
	_ = x[NOSTRG-182]
	_ = x[NOSIZE-183]
	_ = x[NORMAL-184]
	_ = x[NOPRINT-185]
	_ = x[NOOVERFLOW-186]
	_ = x[NOOFL-187]
	_ = x[NONVARYING-188]
	_ = x[NONVAR-189]
	_ = x[NONE-190]
	_ = x[NONCONNECTED-191]
	_ = x[NONCONN-192]
	_ = x[NONASSIGNABLE-193]
	_ = x[NONASGN-194]
	_ = x[NOLOCK-195]
	_ = x[NOINVALIDOP-196]
	_ = x[NOINLINE-197]
	_ = x[NOINIT-198]
	_ = x[NOFIXEDOVERFLOWOFL-199]
	_ = x[NOEXECOPS-200]
	_ = x[NODESCRIPTOR-201]
	_ = x[NOCONVERSION-202]
	_ = x[NOCONV-203]
	_ = x[NOCHECK-204]
	_ = x[NOCHARGRAPHIC-205]
	_ = x[NOCHARG-206]
	_ = x[NCP-207]
	_ = x[NAME-208]
	_ = x[MAIN-209]
	_ = x[M-210]
	_ = x[LOOP-211]
	_ = x[LOCATE-212]
	_ = x[LOCAL-213]
	_ = x[LITTLEENDIAN-214]
	_ = x[LIST-215]
	_ = x[LINKAGE-216]
	_ = x[LINESIZE-217]
	_ = x[LINE-218]
	_ = x[LIMITED-219]
	_ = x[LIKE-220]
	_ = x[LEAVE-221]
	_ = x[LABEL-222]
	_ = x[KEYTO-223]
	_ = x[KEYLOC-224]
	_ = x[KEYLENGTH-225]
	_ = x[KEYFROM-226]
	_ = x[KEYED-227]
	_ = x[KEY-228]
	_ = x[ITERATE-229]
	_ = x[IRREDUCIBLE-230]
	_ = x[IRRED-231]
	_ = x[INVALIDOP-232]
	_ = x[INTO-233]
	_ = x[INTERACTIVE-234]
	_ = x[INTER-235]
	_ = x[INTERNAL-236]
	_ = x[INT-237]
	_ = x[INPUT-238]
	_ = x[INOUT-239]
	_ = x[INONLY-240]
	_ = x[INLINE-241]
	_ = x[INITIAL-242]
	_ = x[INIT-243]
	_ = x[INI-244]
	_ = x[INDFOR-245]
	_ = x[INDEXED-246]
	_ = x[INDEXAREA-247]
	_ = x[INCLUDE-248]
	_ = x[IN-249]
	_ = x[IMPORTED-250]
	_ = x[IGNORE-251]
	_ = x[IF-252]
	_ = x[IEEE-253]
	_ = x[I-254]
	_ = x[HEXADEC-255]
	_ = x[HANDLE-256]
	_ = x[GX-257]
	_ = x[GRAPHIC-258]
	_ = x[GOTO-259]
	_ = x[GO-260]
	_ = x[GET-261]
	_ = x[GENKEY-262]
	_ = x[GENERIC-263]
	_ = x[G-264]
	_ = x[FS-265]
	_ = x[FROMALIEN-266]
	_ = x[FROM-267]
	_ = x[FREE-268]
	_ = x[FORTRAN-269]
	_ = x[FORMAT-270]
	_ = x[FOREVER-271]
	_ = x[FORCE-272]
	_ = x[FLUSH-273]
	_ = x[FLOAT-274]
	_ = x[FIXEDOVERFLOW-275]
	_ = x[FOFL-276]
	_ = x[FIXED-277]
	_ = x[FINISH-278]
	_ = x[FILE-279]
	_ = x[FETCHABLE-280]
	_ = x[FETCH-281]
	_ = x[FBS-282]
	_ = x[FB-283]
	_ = x[F-284]
	_ = x[EXTERNAL-285]
	_ = x[EXT-286]
	_ = x[EXPORTS-287]
	_ = x[EXIT-288]
	_ = x[EXEC-289]
	_ = x[EXCLUSIVE-290]
	_ = x[EXCL-291]
	_ = x[EVENT-292]
	_ = x[ERROR-293]
	_ = x[ENVIRONMENT-294]
	_ = x[ENV-295]
	_ = x[ENTRY-296]
	_ = x[ENDPAGE-297]
	_ = x[ENDFILE-298]
	_ = x[END-299]
	_ = x[ELSE-300]
	_ = x[EDIT-301]
	_ = x[DOWNTHRU-302]
	_ = x[DO-303]
	_ = x[DISPLAY-304]
	_ = x[DIRECT-305]
	_ = x[DIMACROSS-306]
	_ = x[DIMENSION-307]
	_ = x[DIM-308]
	_ = x[DETACH-309]
	_ = x[DESCRIPTORS-310]
	_ = x[DESCRIPTOR-311]
	_ = x[DELETE-312]
	_ = x[DELAY-313]
	_ = x[DEFINE-314]
	_ = x[DEFAULT-315]
	_ = x[DFT-316]
	_ = x[DEFINED-317]
	_ = x[DEF-318]
	_ = x[DECLARE-319]
	_ = x[DCL-320]
	_ = x[DECIMAL-321]
	_ = x[DEC-322]
	_ = x[DEACTIVATE-323]
	_ = x[DEACT-324]
	_ = x[DB-325]
	_ = x[DATE-326]
	_ = x[DATA-327]
	_ = x[D-328]
	_ = x[CTLASA-329]
	_ = x[CTL360-330]
	_ = x[COPY-331]
	_ = x[CONVERSION-332]
	_ = x[CONV-333]
	_ = x[CONTROLLED-334]
	_ = x[CTL-335]
	_ = x[CONSTANT-336]
	_ = x[CONST-337]
	_ = x[CONSECUTIVE-338]
	_ = x[CONNECTED-339]
	_ = x[CONN-340]
	_ = x[CONDITION-341]
	_ = x[COND-342]
	_ = x[COMPLEX-343]
	_ = x[CPLX-344]
	_ = x[COLUMN-345]
	_ = x[COL-346]
	_ = x[COBOL-347]
	_ = x[CLOSE-348]
	_ = x[CHECK-349]
	_ = x[CHARGRAPHIC-350]
	_ = x[CHARG-351]
	_ = x[CHARACTER-352]
	_ = x[CHAR-353]
	_ = x[CELL-354]
	_ = x[CDECL-355]
	_ = x[CALL-356]
	_ = x[C-357]
	_ = x[BYVALUE-358]
	_ = x[BYADDR-359]
	_ = x[BY-360]
	_ = x[BX-361]
	_ = x[BUILTIN-362]
	_ = x[BUFSP-363]
	_ = x[BUFNI-364]
	_ = x[BUFND-365]
	_ = x[BUFFOFF-366]
	_ = x[BUFFERS-367]
	_ = x[BUFF-368]
	_ = x[BUFFERED-369]
	_ = x[BUF-370]
	_ = x[BLKSIZE-371]
	_ = x[BKWD-372]
	_ = x[BIT-373]
	_ = x[BINARY-374]
	_ = x[BIN-375]
	_ = x[BIGENDIAN-376]
	_ = x[BEGIN-377]
	_ = x[BASED-378]
	_ = x[BACKWARDS-379]
	_ = x[B4-380]
	_ = x[B3-381]
	_ = x[B2-382]
	_ = x[B1-383]
	_ = x[B-384]
	_ = x[AUTOMATIC-385]
	_ = x[AUTO-386]
	_ = x[ATTENTION-387]
	_ = x[ATTN-388]
	_ = x[ATTACH-389]
	_ = x[ASSIGNABLE-390]
	_ = x[ASGN-391]
	_ = x[ASSEMBLER-392]
	_ = x[ASM-393]
	_ = x[ASCII-394]
	_ = x[AREA-395]
	_ = x[ANYCONDITION-396]
	_ = x[ANYCOND-397]
	_ = x[ALLOCATE-398]
	_ = x[ALLOC-399]
	_ = x[ALIGNED-400]
	_ = x[ALIAS-401]
	_ = x[ADDBUFF-402]
	_ = x[ACTIVATE-403]
	_ = x[ACT-404]
	_ = x[ABNORMAL-405]
	_ = x[A-406]
	_ = x[tokenCount-407]
}

const _tokens_pli_name = "EOFnameliteralnewlineopop=opop=*([)],;:....ZERODIVIDEZDIVXUXNXWXWRITEWINMAINWHILEWHENEVERWHENWAITWIDECHARWCHARVSAMVSVBSVBVARIABLEVARYINGZVARZVARYINGVARVALUERANGEVALUELISTFROMVALUELISTVALUEVALIDATEVUPTHRUUPDATEUNTILUNSIGNEDUNSUNLOCKUNIONUNDERFLOWUFLUNDEFINEDFILEUNDFUNCONNECTEDUNBUFFEREDUNBUFFUNALIGNEDUNALUTYPETSTACKTRKOFLTRANSMITTRANSIENTTPTOTALTOTITLETHREADTHENTASKSYSTEMSUPPRESSSUPPORTSUBSCRIPTRANGESUBRGSUBSTRUCTURESTRINGVALUESTRINGSIZESTRZSTRINGRANGESTRGSTRINGSTREAMSTORAGESTOPSTDCALLSTATICSNAPSKIPSIZESISSIGNEDSIGNALSETSEQUENTIALSEQLSELECTSCALARVARYINGREWRITEREVERTREUSERETURNSRETURNRETCODERESIGNALRESERVESRESERVEDREREADREPLYREPLACEREPEATREORDERRENAMERELEASEREGIONALREFERREENTRANTREDUCIBLEREDRECURSIVERECSIZERECORDREALREADRANGERPUTPROCEDUREPROCPRIORITYPRINTPRECISIONPRECPOSITIONPOSPOINTERPTRPICTUREPICPENDINGPASSWORDPARAMETERPARMPAGESIZEPAGEPACKEDPACKAGEPOVERFLOWOFLOUTPUTOUTONLYOTHERWISEOTHERORDINALORDEROPTLINKOPTIONSOPTIONALOPENONOFFSETNULLINITNOZERODIVIDENOZDIVNOWRITENOUNDERFLOWNOUFLNOTENOSUBSCRIPTRANGENOSUBRGNOSTRINGSIZENOSTRZNOSTRINGRANGENOSTRGNOSIZENORMALNOPRINTNOOVERFLOWNOOFLNONVARYINGNONVARNONENONCONNECTEDNONCONNNONASSIGNABLENONASGNNOLOCKNOINVALIDOPNOINLINENOINITNOFIXEDOVERFLOWNOEXECOPSNODESCRIPTORNOCONVERSIONNOCONVNOCHECKNOCHARGRAPHICNOCHARGNCPNAMEMAINMLOOPLOCATELOCALLITTLEENDIANLISTLINKAGELINESIZELINELIMITEDLIKELEAVELABELKEYTOKEYLOCKEYLENGTHKEYFROMKEYEDKEYITERATEIRREDUCIBLEIRREDINVALIDOPINTOINTERACTIVEINTERINTERNALINTINPUTINOUTINONLYINLINEINITIALINITINIINDFORINDEXEDINDEXAREAINCLUDEINIMPORTEDIGNOREIFIEEEIHEXADECHANDLEGXGRAPHICGOTOGOGETGENKEYGENERICGFSFROMALIENFROMFREEFORTRANFORMATFOREVERFORCEFLUSHFLOATFIXEDOVERFLOWFOFLFIXEDFINISHFILEFETCHABLEFETCHFBSFBFEXTERNALEXTEXPORTSEXITEXECEXCLUSIVEEXCLEVENTERRORENVIRONMENTENVENTRYENDPAGEENDFILEENDELSEEDITDOWNTHRUDODISPLAYDIRECTDIMACROSSDIMENSIONDIMDETACHDESCRIPTORSDESCRIPTORDELETEDELAYDEFINEDEFAULTDFTDEFINEDDEFDECLAREDCLDECIMALDECDEACTIVATEDEACTDBDATEDATADCTLASACTL360COPYCONVERSIONCONVCONTROLLEDCTLCONSTANTCONSTCONSECUTIVECONNECTEDCONNCONDITIONCONDCOMPLEXCPLXCOLUMNCOLCOBOLCLOSECHECKCHARGRAPHICCHARGCHARACTERCHARCELLCDECLCALLCBYVALUEBYADDRBYBXBUILTINBUFSPBUFNIBUFNDBUFFOFFBUFFERSBUFFBUFFEREDBUFBLKSIZEBKWDBITBINARYBINBIGENDIANBEGINBASEDBACKWARDSB4B3B2B1BAUTOMATICAUTOATTENTIONATTNATTACHASSIGNABLEASGNASSEMBLERASMASCIIAREAANYCONDITIONANYCONDALLOCATEALLOCALIGNEDALIASADDBUFFACTIVATEACTABNORMALA"

var _tokens_pli_index = [...]uint16{0, 3, 7, 14, 21, 23, 26, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 43, 53, 57, 59, 61, 62, 64, 69, 76, 81, 89, 93, 97, 105, 110, 114, 116, 119, 121, 129, 137, 141, 148, 151, 161, 174, 183, 188, 196, 197, 203, 209, 214, 222, 225, 231, 236, 245, 248, 261, 265, 276, 286, 292, 301, 305, 306, 310, 316, 322, 330, 339, 341, 346, 348, 353, 359, 363, 367, 373, 381, 388, 402, 407, 410, 419, 430, 440, 444, 455, 459, 465, 471, 478, 482, 489, 495, 499, 503, 507, 510, 516, 522, 525, 535, 539, 545, 558, 565, 571, 576, 583, 589, 596, 604, 612, 620, 626, 631, 638, 644, 651, 657, 664, 672, 677, 686, 695, 698, 707, 714, 720, 724, 728, 733, 734, 737, 746, 750, 758, 763, 772, 776, 784, 787, 794, 797, 804, 807, 814, 822, 831, 835, 843, 847, 853, 860, 861, 869, 872, 878, 885, 894, 899, 906, 911, 918, 925, 933, 937, 939, 945, 953, 965, 971, 978, 989, 994, 998, 1014, 1021, 1033, 1039, 1052, 1058, 1064, 1070, 1077, 1087, 1092, 1102, 1108, 1112, 1124, 1131, 1144, 1151, 1157, 1168, 1176, 1182, 1197, 1206, 1218, 1230, 1236, 1243, 1256, 1263, 1266, 1270, 1274, 1275, 1279, 1285, 1290, 1302, 1306, 1313, 1321, 1325, 1332, 1336, 1341, 1346, 1351, 1357, 1366, 1373, 1378, 1381, 1388, 1399, 1404, 1413, 1417, 1428, 1433, 1441, 1444, 1449, 1454, 1460, 1466, 1473, 1477, 1480, 1486, 1493, 1502, 1509, 1511, 1519, 1525, 1527, 1531, 1532, 1539, 1545, 1547, 1554, 1558, 1560, 1563, 1569, 1576, 1577, 1579, 1588, 1592, 1596, 1603, 1609, 1616, 1621, 1626, 1631, 1644, 1648, 1653, 1659, 1663, 1672, 1677, 1680, 1682, 1683, 1691, 1694, 1701, 1705, 1709, 1718, 1722, 1727, 1732, 1743, 1746, 1751, 1758, 1765, 1768, 1772, 1776, 1784, 1786, 1793, 1799, 1808, 1817, 1820, 1826, 1837, 1847, 1853, 1858, 1864, 1871, 1874, 1881, 1884, 1891, 1894, 1901, 1904, 1914, 1919, 1921, 1925, 1929, 1930, 1936, 1942, 1946, 1956, 1960, 1970, 1973, 1981, 1986, 1997, 2006, 2010, 2019, 2023, 2030, 2034, 2040, 2043, 2048, 2053, 2058, 2069, 2074, 2083, 2087, 2091, 2096, 2100, 2101, 2108, 2114, 2116, 2118, 2125, 2130, 2135, 2140, 2147, 2154, 2158, 2166, 2169, 2176, 2180, 2183, 2189, 2192, 2201, 2206, 2211, 2220, 2222, 2224, 2226, 2228, 2229, 2238, 2242, 2251, 2255, 2261, 2271, 2275, 2284, 2287, 2292, 2296, 2308, 2315, 2323, 2328, 2335, 2340, 2347, 2355, 2358, 2366, 2367, 2367}

func (i tokens_pli) String() string {
	i -= 1
	if i >= tokens_pli(len(_tokens_pli_index)-1) {
		return "tokens_pli(" + strconv.FormatInt(int64(i+1), 10) + ")"
	}
	return _tokens_pli_name[_tokens_pli_index[i]:_tokens_pli_index[i+1]]
}
